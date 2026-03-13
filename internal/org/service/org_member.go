package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dimasbaguspm/fluxis/internal/org/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
	"github.com/dimasbaguspm/fluxis/pkg/syncx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) ListMembers(ctx context.Context, orgId pgtype.UUID, q domain.OrganisationMembersSearchModel) (domain.OrganisationMembersPagedModel, error) {
	q.ApplyDefaults()

	members, err := s.Repo.ListOrgMembers(ctx, repository.ListOrgMembersParams{
		OrgID:   orgId,
		Column2: q.UserID,
		Column3: q.Email,
		Column4: q.DisplayName,
		Limit:   int32(q.PageSize),
		Offset:  int32((q.PageNumber - 1) * q.PageSize),
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			emptyResult := domain.OrganisationMembersPagedModel{}
			return emptyResult.Empty(q.PageNumber, q.PageSize), nil
		}
		emptyResult := domain.OrganisationMembersPagedModel{}
		return emptyResult.Empty(q.PageNumber, q.PageSize), fmt.Errorf("get org members: %w", err)
	}

	totalCount := int64(0)
	data := make([]domain.OrganisationMemberModel, 0, len(members))

	for _, member := range members {
		totalCount = member.TotalCount
		data = append(data, domain.OrganisationMemberModel{
			UserID:   member.UserID,
			Name:     member.DisplayName,
			Email:    member.Email,
			Role:     string(member.Role),
			JoinedAt: member.JoinedAt.Time,
		})
	}

	totalPages := 0
	if totalCount > 0 {
		totalPages = int((totalCount + int64(q.PageSize) - 1) / int64(q.PageSize))
	}

	return domain.OrganisationMembersPagedModel{
		Items:      data,
		TotalCount: int(totalCount),
		TotalPages: totalPages,
		PageNumber: q.PageNumber,
		PageSize:   q.PageSize,
	}, nil
}

func (s *Service) AddMember(ctx context.Context, orgId pgtype.UUID, p domain.OrganisationMemberCreateModel) error {
	var userId pgtype.UUID
	if err := userId.Scan(p.UserId); err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	var (
		org  domain.OrganisationModel
		user domain.UserModel
	)

	err := syncx.Run(ctx, func(ctx context.Context) (err error) {
		org, err = s.GetOrgById(ctx, orgId)
		return err
	}, func(ctx context.Context) (err error) {
		user, err = s.User.GetSingleUserById(ctx, userId)
		return err
	})

	if err != nil {
		return fmt.Errorf("get org or user: %w", err)
	}

	_, err = s.Repo.CreateOrgMember(ctx, repository.CreateOrgMemberParams{
		OrgID:  org.ID,
		UserID: user.ID,
		Role:   repository.OrgRole(p.Role),
	})

	if err != nil {
		return fmt.Errorf("add org member err: %w", err)
	}

	if err := s.Bus.Publish(ctx, pubsub.OrgMemberAdded, map[string]string{
		"orgId":  uuid.UUID(org.ID.Bytes).String(),
		"userId": uuid.UUID(user.ID.Bytes).String(),
	}); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.OrgMemberAdded), "error", err)
	}

	return nil
}

func (s *Service) UpdateMemberRole(ctx context.Context, orgId, userId pgtype.UUID, p domain.OrganisationMemberUpdateModel) error {
	_, err := s.Repo.UpdateOrgMemberRole(ctx, repository.UpdateOrgMemberRoleParams{
		OrgID:  orgId,
		UserID: userId,
		Role:   repository.OrgRole(p.Role),
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrgMemberNotFound
		}
		return fmt.Errorf("update org member role err: %w", err)
	}

	if err := s.Bus.Publish(ctx, pubsub.OrgMemberUpdated, map[string]string{
		"orgId":  uuid.UUID(orgId.Bytes).String(),
		"userId": uuid.UUID(userId.Bytes).String(),
	}); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.OrgMemberUpdated), "error", err)
	}

	return nil
}

func (s *Service) RemoveMember(ctx context.Context, orgId, userId pgtype.UUID) error {
	err := s.Repo.DeleteOrgMember(ctx, repository.DeleteOrgMemberParams{
		OrgID:  orgId,
		UserID: userId,
	})

	if err != nil {
		return fmt.Errorf("delete org member role err: %w", err)
	}

	if err := s.Bus.Publish(ctx, pubsub.OrgMemberRemoved, map[string]string{
		"orgId":  uuid.UUID(orgId.Bytes).String(),
		"userId": uuid.UUID(userId.Bytes).String(),
	}); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.OrgMemberRemoved), "error", err)
	}

	return nil
}
