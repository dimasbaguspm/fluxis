package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dimasbaguspm/fluxis/internal/org/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) GetListOrganisationMembers(ctx context.Context, orgId pgtype.UUID) ([]domain.OrganisationMemberModel, error) {
	members, err := s.Repo.ListOrgMembers(ctx, repository.ListOrgMembersParams{
		OrgID: orgId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.OrganisationMemberModel{}, nil
		}
		return []domain.OrganisationMemberModel{}, fmt.Errorf("get org members: %w", err)
	}

	data := make([]domain.OrganisationMemberModel, len(members))

	for _, member := range members {
		data = append(data, domain.OrganisationMemberModel{
			UserID:   member.UserID,
			Name:     member.DisplayName,
			Email:    member.Email,
			Role:     string(member.Role),
			JoinedAt: member.JoinedAt.Time,
		})
	}

	return data, nil
}

func (s *Service) AddOrganisationMember(ctx context.Context, orgId pgtype.UUID, p domain.OrganisationMemberCreateModel) error {
	_, err := s.Repo.CreateOrgMember(ctx, repository.CreateOrgMemberParams{
		OrgID:  orgId,
		UserID: p.UserID,
		Role:   repository.OrgRole(p.Role),
	})

	if err != nil {
		return fmt.Errorf("add org member err: %w", err)
	}

	return nil
}

func (s *Service) UpdateOrganisationMemberRole(ctx context.Context, orgId pgtype.UUID, p domain.OrganisationMemberUpdateModel) error {
	_, err := s.Repo.UpdateOrgMemberRole(ctx, repository.UpdateOrgMemberRoleParams{
		OrgID:  orgId,
		UserID: p.UserID,
		Role:   repository.OrgRole(p.Role),
	})

	if err != nil {
		return fmt.Errorf("update org member role err: %w", err)
	}

	return nil
}

func (s *Service) RemoveOrganisationMember(ctx context.Context, orgId pgtype.UUID, p domain.OrganisationMemberRemoveModel) error {
	err := s.Repo.DeleteOrgMember(ctx, repository.DeleteOrgMemberParams{
		OrgID:  orgId,
		UserID: p.UserID,
	})

	if err != nil {
		return fmt.Errorf("delete org member role err: %w", err)
	}

	return nil
}
