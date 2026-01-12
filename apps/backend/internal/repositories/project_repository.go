package repositories

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pgx *pgxpool.Pool
}

func NewProjectRepository(pgx *pgxpool.Pool) ProjectRepository {
	return ProjectRepository{pgx}
}

func (pr ProjectRepository) GetPaginated(ctx context.Context, query models.ProjectSearchModel) (models.ProjectPaginatedModel, error) {
	sortMap := map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"status":    "status",
	}
	sortColumn, ok := sortMap[query.SortBy]
	if !ok {
		sortColumn = "created_at" // default fallback
	}

	sortDirection := "DESC"
	if query.SortOrder == "asc" {
		sortDirection = "ASC"
	}

	offset := (query.PageNumber - 1) * query.PageSize
	searchPattern := "%" + query.Query + "%"

	sql := `
		WITH filtered AS (
			SELECT id, name, description, status, created_at, updated_at
			FROM projects
			WHERE deleted_at IS NULL
				AND ($1::uuid[] IS NULL OR CARDINALITY($1::uuid[]) = 0 OR id = ANY($1))
				AND ($2::text[] IS NULL OR CARDINALITY($2::text[]) = 0 OR status::text = ANY($2))
				AND ($3 = '' OR name ILIKE $3 OR description ILIKE $3)
		),
		counted AS (
			SELECT COUNT(*) as total FROM filtered
		)
		SELECT 
			f.id,
			f.name,
			COALESCE(f.description, '') as description,
			f.status::text as status,
			f.created_at,
			f.updated_at,
			c.total
		FROM filtered f
		CROSS JOIN counted c
		ORDER BY f.` + sortColumn + ` ` + sortDirection + `
		LIMIT $4 OFFSET $5
	`

	rows, err := pr.pgx.Query(ctx, sql, query.ID, query.Status, searchPattern, query.PageSize, offset)
	if err != nil {
		return models.ProjectPaginatedModel{}, huma.Error400BadRequest("Unable to query projects", err)
	}
	defer rows.Close()

	var items []models.ProjectModel
	var totalCount int

	for rows.Next() {
		var item models.ProjectModel
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Status, &item.CreatedAt, &item.UpdatedAt, &totalCount)
		if err != nil {
			return models.ProjectPaginatedModel{}, huma.Error400BadRequest("Unable to scan project data", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return models.ProjectPaginatedModel{}, huma.Error400BadRequest("Error reading project rows", err)
	}

	if items == nil {
		items = []models.ProjectModel{}
	}

	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + query.PageSize - 1) / query.PageSize
	}

	return models.ProjectPaginatedModel{
		Items:      items,
		PageNumber: query.PageNumber,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}, nil
}

func (pr ProjectRepository) GetDetail(ctx context.Context, id string) (models.ProjectModel, error) {

	var data models.ProjectModel

	sql := `SELECT id, name, description, status, created_at, updated_at
				FROM projects
				WHERE id = $1::uuid AND deleted_at IS NULL`

	err := pr.pgx.QueryRow(ctx, sql, id).Scan(&data.ID, &data.Name, &data.Description, &data.Status, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ProjectModel{}, huma.Error404NotFound("No project found")
		}
		return models.ProjectModel{}, huma.Error400BadRequest("Unable to query the project details", err)
	}

	return data, nil
}

func (pr ProjectRepository) Create(ctx context.Context, payload models.ProjectCreateModel) (models.ProjectModel, error) {
	var data models.ProjectModel

	sql := `INSERT into projects (name, description, status) VALUES ($1, $2, $3)
				RETURNING id, name, description, status, created_at, updated_at`

	err := pr.pgx.QueryRow(ctx, sql, payload.Name, payload.Description, payload.Status).Scan(&data.ID, &data.Name, &data.Description, &data.Status, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		return models.ProjectModel{}, huma.Error400BadRequest("Unable to create project", err)
	}

	return data, nil
}

func (pr ProjectRepository) Update(ctx context.Context, id string, payload models.ProjectUpdateModel) (models.ProjectModel, error) {
	var data models.ProjectModel

	sql := `UPDATE projects
			SET name = COALESCE($1, name),
				description = COALESCE($2, description),
				status = COALESCE($3, status),
				updated_at = CURRENT_TIMESTAMP
				WHERE id = $4::uuid AND deleted_at IS NULL
			RETURNING id, name, description, status, created_at, updated_at`
	err := pr.pgx.QueryRow(ctx, sql, payload.Name, payload.Description, payload.Status, id).Scan(&data.ID, &data.Name, &data.Description, &data.Status, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ProjectModel{}, huma.Error404NotFound("No project found")
		}
		return models.ProjectModel{}, huma.Error400BadRequest("Unable to update the record", err)
	}

	return data, nil
}

func (pr ProjectRepository) Delete(ctx context.Context, id string) error {
	sql := `UPDATE projects
					SET deleted_at = CURRENT_TIMESTAMP      
					WHERE id = $1::uuid AND deleted_at IS NULL`

	cmdTag, err := pr.pgx.Exec(ctx, sql, id)
	if err != nil {
		return huma.Error400BadRequest("Unable to delete the record", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return huma.Error404NotFound("No project found")
	}

	return nil
}
