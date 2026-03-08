package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dimasbaguspm/fluxis/internal/org/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/transformer"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrOrgNotFound = httpx.NotFound("organisation not found")
	ErrSlugIsTaken = httpx.Conflict("slug has been taken")
)

func (s *Service) ListOrgs(ctx context.Context, q domain.OrganisationSearchModel) ([]domain.OrganisationModel, error) {
	orgs, err := s.Repo.ListOrg(ctx, repository.ListOrgParams{
		Column1: q.OrgId,
		Column2: q.UserId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.OrganisationModel{}, nil
		}
		return []domain.OrganisationModel{}, fmt.Errorf("list orgs: %w", err)
	}

	data := make([]domain.OrganisationModel, len(orgs))

	for _, org := range orgs {
		data = append(data, domain.OrganisationModel{
			ID:        org.ID,
			Name:      org.Name,
			Slug:      org.Slug,
			CreatedAt: org.CreatedAt.Time,
			UpdatedAt: org.UpdatedAt.Time,
		})
	}

	return data, nil
}

func (s *Service) GetOrgById(ctx context.Context, id pgtype.UUID) (domain.OrganisationModel, error) {
	org, err := s.Repo.GetOrgById(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.OrganisationModel{}, ErrOrgNotFound
		}
		return domain.OrganisationModel{}, fmt.Errorf("get org by id: %w", err)
	}

	totalMembers, _ := s.Repo.CountOrgMembers(ctx, repository.CountOrgMembersParams{
		OrgID: org.ID,
	})

	return domain.OrganisationModel{
		ID:           org.ID,
		Name:         org.Name,
		Slug:         org.Slug,
		TotalMembers: totalMembers,
		CreatedAt:    org.CreatedAt.Time,
		UpdatedAt:    org.UpdatedAt.Time,
	}, nil
}

func (s *Service) GetOrgBySlug(ctx context.Context, slug string) (domain.OrganisationModel, error) {
	org, err := s.Repo.GetOrgBySlug(ctx, slug)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.OrganisationModel{}, ErrOrgNotFound
		}
		return domain.OrganisationModel{}, fmt.Errorf("get org by id: %w", err)
	}

	totalMembers, _ := s.Repo.CountOrgMembers(ctx, repository.CountOrgMembersParams{
		OrgID: org.ID,
	})

	return domain.OrganisationModel{
		ID:           org.ID,
		Name:         org.Name,
		Slug:         org.Slug,
		TotalMembers: totalMembers,
		CreatedAt:    org.CreatedAt.Time,
		UpdatedAt:    org.UpdatedAt.Time,
	}, nil
}

func (s *Service) CreateOrg(ctx context.Context, p domain.OrganisationCreateModel) (domain.OrganisationModel, error) {
	userID := httpx.MustUserID(ctx)
	org, err := s.Repo.CreateOrg(ctx, repository.CreateOrgParams{
		Name: p.Name,
		Slug: transformer.CreateSlug(p.Name),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.OrganisationModel{}, ErrSlugIsTaken
		}
		return domain.OrganisationModel{}, fmt.Errorf("create org: %w", err)
	}

	s.Repo.CreateOrgMember(ctx, repository.CreateOrgMemberParams{
		OrgID:  org.ID,
		UserID: userID,
		Role:   repository.OrgRoleAdmin,
	})

	totalMembers, _ := s.Repo.CountOrgMembers(ctx, repository.CountOrgMembersParams{
		OrgID: org.ID,
	})

	return domain.OrganisationModel{
		ID:           org.ID,
		Name:         org.Name,
		Slug:         org.Slug,
		TotalMembers: totalMembers,
		CreatedAt:    org.CreatedAt.Time,
		UpdatedAt:    org.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateOrg(ctx context.Context, id pgtype.UUID, p domain.OrganisationUpdateModel) (domain.OrganisationModel, error) {
	org, err := s.Repo.UpdateOrg(ctx, repository.UpdateOrgParams{
		ID:      id,
		Column1: p.Name,
		Column2: transformer.CreateSlug(p.Name),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.OrganisationModel{}, ErrOrgNotFound
		}
		return domain.OrganisationModel{}, fmt.Errorf("update org: %w", err)
	}

	totalMembers, _ := s.Repo.CountOrgMembers(ctx, repository.CountOrgMembersParams{
		OrgID: org.ID,
	})

	return domain.OrganisationModel{
		ID:           org.ID,
		Name:         org.Name,
		Slug:         org.Slug,
		TotalMembers: totalMembers,
		CreatedAt:    org.CreatedAt.Time,
		UpdatedAt:    org.UpdatedAt.Time,
	}, nil
}

func (s *Service) DeleteOrg(ctx context.Context, id pgtype.UUID) error {
	err := s.Repo.DeleteOrg(ctx, id)
	if err != nil {
		return fmt.Errorf("delete org: %w", err)
	}
	return nil
}
