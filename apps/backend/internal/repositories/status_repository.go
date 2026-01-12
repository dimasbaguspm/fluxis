package repositories

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatusRepository struct {
	pgx *pgxpool.Pool
}

func NewStatusRepository(pgx *pgxpool.Pool) StatusRepository {
	return StatusRepository{pgx}
}

func (sr StatusRepository) GetByProject(ctx context.Context, projectId string) ([]models.StatusModel, error) {
	sql := `SELECT id, project_id, name, slug, position, is_default, created_at, updated_at
		FROM statuses
		WHERE project_id = $1 AND deleted_at IS NULL
		ORDER BY position ASC`

	rows, err := sr.pgx.Query(ctx, sql, projectId)
	if err != nil {
		return nil, huma.Error400BadRequest("Unable to query statuses", err)
	}
	defer rows.Close()

	var items []models.StatusModel
	for rows.Next() {
		var s models.StatusModel
		err := rows.Scan(&s.ID, &s.ProjectID, &s.Name, &s.Slug, &s.Position, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, huma.Error400BadRequest("Unable to scan status", err)
		}
		items = append(items, s)
	}

	if err := rows.Err(); err != nil {
		return nil, huma.Error400BadRequest("Error reading status rows", err)
	}

	if items == nil {
		items = []models.StatusModel{}
	}

	return items, nil
}

func (sr StatusRepository) Create(ctx context.Context, projectId string, payload models.StatusCreateModel) (models.StatusModel, error) {
	var data models.StatusModel

	// generate slug in Go using helper
	slug := common.Slugify(payload.Name)

	sql := `INSERT INTO statuses (project_id, name, slug, position, is_default)
		VALUES ($1, $2, $3,
			(SELECT COALESCE(MAX(position), -1) + 1 FROM statuses WHERE project_id = $1 AND deleted_at IS NULL),
			false)
		RETURNING id, project_id, name, slug, position, is_default, created_at, updated_at`

	err := sr.pgx.QueryRow(ctx, sql, projectId, payload.Name, slug).Scan(
		&data.ID, &data.ProjectID, &data.Name, &data.Slug, &data.Position, &data.IsDefault, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		return models.StatusModel{}, huma.Error400BadRequest("Unable to create status", err)
	}

	return data, nil
}

func (sr StatusRepository) Update(ctx context.Context, id string, payload models.StatusUpdateModel) (models.StatusModel, error) {
	var data models.StatusModel

	// compute slug in Go and pass as parameter
	slug := common.Slugify(payload.Name)

	sql := `UPDATE statuses
		SET name = COALESCE(NULLIF($1, ''), name),
			slug = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING id, project_id, name, slug, position, is_default, created_at, updated_at`

	err := sr.pgx.QueryRow(ctx, sql, payload.Name, slug, id).Scan(
		&data.ID, &data.ProjectID, &data.Name, &data.Slug, &data.Position, &data.IsDefault, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.StatusModel{}, huma.Error404NotFound("No status found")
		}
		return models.StatusModel{}, huma.Error400BadRequest("Unable to update status", err)
	}

	return data, nil
}

func (sr StatusRepository) Delete(ctx context.Context, id string) error {
	sql := `UPDATE statuses 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND deleted_at IS NULL`

	cmdTag, err := sr.pgx.Exec(ctx, sql, id)
	if err != nil {
		return huma.Error400BadRequest("Unable to delete status", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return huma.Error404NotFound("No status found")
	}

	return nil
}

func (sr StatusRepository) Reorder(ctx context.Context, projectId string, ids []string) ([]models.StatusModel, error) {
	tx, err := sr.pgx.Begin(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Unable to start transaction", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	sql := `WITH np AS (
		SELECT u.id::uuid AS id, (u.ord - 1) AS pos
		FROM unnest($1::uuid[]) WITH ORDINALITY AS u(id, ord)
	),
	upd AS (
	  UPDATE statuses s
	  SET position = np.pos
	  FROM np
	  WHERE s.id = np.id AND s.project_id = $2 AND s.deleted_at IS NULL
	  RETURNING s.id, s.project_id, s.name, s.slug, s.position, s.is_default, s.created_at, s.updated_at
	)
	SELECT id, project_id, name, slug, position, is_default, created_at, updated_at
	FROM upd
	ORDER BY position ASC`

	rows, err := tx.Query(ctx, sql, ids, projectId)
	if err != nil {
		return nil, huma.Error400BadRequest("Unable to reorder statuses", err)
	}
	defer rows.Close()

	var items []models.StatusModel
	for rows.Next() {
		var s models.StatusModel
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Name, &s.Slug, &s.Position, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, huma.Error400BadRequest("Unable to scan reordered status", err)
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, huma.Error400BadRequest("Error reading reordered rows", err)
	}

	if int64(len(items)) != int64(len(ids)) {
		return nil, huma.Error400BadRequest("Unable to update status positions: invalid id or not in project")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, huma.Error400BadRequest("Unable to commit reorder transaction", err)
	}

	return items, nil
}

// ValidateReorderCounts returns (total, matched) where total is number of non-deleted statuses
// for the project, and matched is how many of the provided ids belong to that project.
func (sr StatusRepository) ValidateReorderCounts(ctx context.Context, projectId string, ids []string) (int, int, error) {
	var total, matched int
	sql := `SELECT
		(SELECT COUNT(1) FROM statuses WHERE project_id = $1 AND deleted_at IS NULL) AS total,
		(SELECT COUNT(1) FROM statuses WHERE project_id = $1 AND id = ANY($2::uuid[]) AND deleted_at IS NULL) AS matched`

	err := sr.pgx.QueryRow(ctx, sql, projectId, ids).Scan(&total, &matched)
	if err != nil {
		return 0, 0, huma.Error400BadRequest("Unable to validate reorder payload", err)
	}
	return total, matched, nil
}
