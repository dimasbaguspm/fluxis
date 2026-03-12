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
	ErrOrgNotFound         = httpx.NotFound("organisation not found")
	ErrSlugIsTaken        = httpx.Conflict("slug has been taken")
	ErrOrgMemberNotFound  = httpx.NotFound("organisation member not found")
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

	data := make([]domain.OrganisationModel, 0, len(orgs))

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

func (s *Service) SearchOrganisations(ctx context.Context, q domain.Organisations) (domain.OrganisationPagedModel, error) {
	q.ApplyDefaults()

	rows, err := s.Repo.SearchOrganisations(ctx, repository.SearchOrganisationsParams{
		Column1: q.ID,
		Column2: q.Name,
		Column3: q.SortBy,
		Column4: q.SortOrder,
		Limit:   int32(q.PageSize),
		Column6: int32(q.PageNumber),
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			emptyPage := domain.OrganisationPagedModel{}
			return emptyPage.Empty(q.PageNumber, q.PageSize), nil
		}
		return domain.OrganisationPagedModel{}, fmt.Errorf("search organisations: %w", err)
	}

	var totalCount int64
	items := make([]domain.OrganisationModel, 0, len(rows))

	for _, row := range rows {
		if totalCount == 0 {
			totalCount = row.TotalCount
		}
		items = append(items, domain.OrganisationModel{
			ID:        row.ID,
			Name:      row.Name,
			Slug:      row.Slug,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}

	totalPages := int(totalCount) / q.PageSize
	if int(totalCount)%q.PageSize != 0 {
		totalPages++
	}

	return domain.OrganisationPagedModel{
		Items:      items,
		TotalCount: int(totalCount),
		TotalPages: totalPages,
		Page:       q.PageNumber,
		PageSize:   q.PageSize,
	}, nil
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

	if _, err := s.Repo.CreateOrgMember(ctx, repository.CreateOrgMemberParams{
		OrgID:  org.ID,
		UserID: userID,
		Role:   repository.OrgRoleAdmin,
	}); err != nil {
		return domain.OrganisationModel{}, fmt.Errorf("create org member: %w", err)
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
