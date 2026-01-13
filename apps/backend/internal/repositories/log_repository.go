package repositories

import (
	"context"
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LogRepository struct {
	pgx *pgxpool.Pool
}

func NewLogRepository(pgx *pgxpool.Pool) LogRepository {
	return LogRepository{pgx: pgx}
}

func (lr *LogRepository) Insert(ctx context.Context, entry models.LogCreateModel) error {
	sqlStr := `INSERT INTO logs (project_id, task_id, status_id, entry) VALUES ($1::uuid, $2::uuid, $3::uuid, $4)`
	var taskID interface{}
	if entry.TaskID == nil || *entry.TaskID == "" {
		taskID = nil
	} else {
		taskID = *entry.TaskID
	}
	var statusID interface{}
	if entry.StatusID == nil || *entry.StatusID == "" {
		statusID = nil
	} else {
		statusID = *entry.StatusID
	}

	_, err := lr.pgx.Exec(ctx, sqlStr, entry.ProjectID, taskID, statusID, entry.Entry)
	if err != nil {
		return huma.Error400BadRequest("Unable to write log", err)
	}
	return nil
}

func (lr *LogRepository) GetPaginated(ctx context.Context, projectID string, q models.LogSearchModel) (models.LogPaginatedModel, error) {
	offset := (q.PageNumber - 1) * q.PageSize
	searchPattern := "%" + q.Query + "%"

	query := `WITH filtered AS (
		SELECT id, project_id, task_id, status_id, entry, created_at
		FROM logs
		WHERE project_id = $1::uuid
			AND ($2::uuid[] IS NULL OR CARDINALITY($2::uuid[]) = 0 OR task_id = ANY($2))
			AND ($3::uuid[] IS NULL OR CARDINALITY($3::uuid[]) = 0 OR status_id = ANY($3))
			AND ($4 = '' OR entry ILIKE $4)
	), counted AS (
		SELECT COUNT(*) as total FROM filtered
	)
	SELECT f.id, f.project_id, f.task_id, f.status_id, f.entry, f.created_at, c.total
	FROM filtered f
	CROSS JOIN counted c
	ORDER BY f.created_at DESC
	LIMIT $5 OFFSET $6`

	rows, err := lr.pgx.Query(ctx, query, projectID, q.TaskID, q.StatusID, searchPattern, q.PageSize, offset)
	if err != nil {
		return models.LogPaginatedModel{}, huma.Error400BadRequest("Unable to query logs", err)
	}
	defer rows.Close()

	var items []models.LogModel
	var totalCount int
	for rows.Next() {
		var l models.LogModel
		var taskID sql.NullString
		var statusID sql.NullString
		if err := rows.Scan(&l.ID, &l.ProjectID, &taskID, &statusID, &l.Entry, &l.CreatedAt, &totalCount); err != nil {
			return models.LogPaginatedModel{}, huma.Error400BadRequest("Unable to scan log", err)
		}
		if taskID.Valid {
			t := taskID.String
			l.TaskID = &t
		} else {
			l.TaskID = nil
		}
		if statusID.Valid {
			s := statusID.String
			l.StatusID = &s
		} else {
			l.StatusID = nil
		}
		items = append(items, l)
	}
	if err := rows.Err(); err != nil {
		return models.LogPaginatedModel{}, huma.Error400BadRequest("Error reading log rows", err)
	}
	if items == nil {
		items = []models.LogModel{}
	}

	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + q.PageSize - 1) / q.PageSize
	}

	return models.LogPaginatedModel{Items: items, PageNumber: q.PageNumber, PageSize: q.PageSize, TotalPages: totalPages, TotalCount: totalCount}, nil
}
