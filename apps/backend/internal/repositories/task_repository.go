package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	pgx *pgxpool.Pool
}

func NewTaskRepository(pgx *pgxpool.Pool) TaskRepository {
	return TaskRepository{pgx}
}

func (tr TaskRepository) GetPaginated(ctx context.Context, q models.TaskSearchModel) (models.TaskPaginatedModel, error) {
	sortByMap := map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"priority":  "priority",
		"dueDate":   "due_date",
	}
	sortOrderMap := map[string]string{"asc": "ASC", "desc": "DESC"}
	sortColumn, _ := sortByMap[q.SortBy]
	sortOrder, _ := sortOrderMap[q.SortOrder]

	offset := (q.PageNumber - 1) * q.PageSize
	searchPattern := "%" + q.Query + "%"

	orderClause := ""
	if q.SortBy == "dueDate" {
		orderClause = `CASE WHEN f.due_date IS NOT NULL THEN 0 ELSE 1 END, f.due_date ASC NULLS LAST, f.priority DESC`
	} else {
		orderClause = "f." + sortColumn + " " + sortOrder
	}

	query := `WITH filtered AS (
        SELECT t.id, t.project_id, COALESCE(t.title, '') AS title, COALESCE(t.details, '') AS details, t.status_id, t.priority, t.due_date, t.created_at, t.updated_at
        FROM tasks t
        INNER JOIN projects p ON t.project_id = p.id AND p.deleted_at IS NULL
        WHERE t.deleted_at IS NULL
            AND ($1::uuid[] IS NULL OR CARDINALITY($1::uuid[]) = 0 OR t.id = ANY($1))
            AND ($2::uuid[] IS NULL OR CARDINALITY($2::uuid[]) = 0 OR t.project_id = ANY($2))
            AND ($3::uuid[] IS NULL OR CARDINALITY($3::uuid[]) = 0 OR t.status_id = ANY($3))
            AND ($4 = '' OR t.title ILIKE $4 OR t.details ILIKE $4)
    ), counted AS (
        SELECT COUNT(*) as total FROM filtered
    )
    SELECT f.id, f.project_id, f.title, f.details, f.status_id, f.priority, f.due_date, f.created_at, f.updated_at, c.total
    FROM filtered f
    CROSS JOIN counted c
    ORDER BY ` + orderClause + `
    LIMIT $5 OFFSET $6`

	rows, err := tr.pgx.Query(ctx, query, q.ID, q.ProjectID, q.StatusID, searchPattern, q.PageSize, offset)
	if err != nil {
		return models.TaskPaginatedModel{}, huma.Error400BadRequest("Unable to query tasks", err)
	}
	defer rows.Close()

	var items []models.TaskModel
	var totalCount int
	for rows.Next() {
		var t models.TaskModel
		var statusID sql.NullString
		var dueDate sql.NullTime
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Details, &statusID, &t.Priority, &dueDate, &t.CreatedAt, &t.UpdatedAt, &totalCount); err != nil {
			return models.TaskPaginatedModel{}, huma.Error400BadRequest("Unable to scan task", err)
		}
		if statusID.Valid {
			t.StatusID = statusID.String
		} else {
			t.StatusID = ""
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return models.TaskPaginatedModel{}, huma.Error400BadRequest("Error reading task rows", err)
	}
	if items == nil {
		items = []models.TaskModel{}
	}

	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + q.PageSize - 1) / q.PageSize
	}

	return models.TaskPaginatedModel{Items: items, PageNumber: q.PageNumber, PageSize: q.PageSize, TotalPages: totalPages, TotalCount: totalCount}, nil
}

func (tr TaskRepository) GetDetail(ctx context.Context, id string) (models.TaskModel, error) {
	var t models.TaskModel
	var statusID sql.NullString
	var dueDate sql.NullTime

	query := `SELECT t.id, t.project_id, t.title, t.details, t.status_id, t.priority, t.due_date, t.created_at, t.updated_at
        FROM tasks t
        INNER JOIN projects p ON t.project_id = p.id AND p.deleted_at IS NULL
        WHERE t.id = $1::uuid AND t.deleted_at IS NULL`

	err := tr.pgx.QueryRow(ctx, query, id).Scan(&t.ID, &t.ProjectID, &t.Title, &t.Details, &statusID, &t.Priority, &dueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.TaskModel{}, huma.Error404NotFound("No task found")
		}
		return models.TaskModel{}, huma.Error400BadRequest("Unable to query the task", err)
	}
	if statusID.Valid {
		t.StatusID = statusID.String
	} else {
		t.StatusID = ""
	}
	if dueDate.Valid {
		t.DueDate = &dueDate.Time
	}
	return t, nil
}

func (tr TaskRepository) Create(ctx context.Context, payload models.TaskCreateModel) (models.TaskModel, error) {
	var t models.TaskModel
	query := `INSERT INTO tasks (project_id, title, details, status_id, priority, due_date)
        VALUES ($1::uuid, $2, $3, $4::uuid, $5, $6)
        RETURNING id, project_id, title, details, status_id, priority, due_date, created_at, updated_at`

	var statusParam interface{}
	if payload.StatusID == "" {
		statusParam = nil
	} else {
		statusParam = payload.StatusID
	}

	var statusScan sql.NullString
	var dueDateScan sql.NullTime
	err := tr.pgx.QueryRow(ctx, query, payload.ProjectID, payload.Title, payload.Details, statusParam, payload.Priority, payload.DueDate).Scan(&t.ID, &t.ProjectID, &t.Title, &t.Details, &statusScan, &t.Priority, &dueDateScan, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return models.TaskModel{}, huma.Error400BadRequest("Unable to create task", err)
	}
	if statusScan.Valid {
		t.StatusID = statusScan.String
	} else {
		t.StatusID = ""
	}
	if dueDateScan.Valid {
		t.DueDate = &dueDateScan.Time
	}
	return t, nil
}

func (tr TaskRepository) Update(ctx context.Context, id string, payload models.TaskUpdateModel) (models.TaskModel, error) {
	var t models.TaskModel
	query := `UPDATE tasks t SET title = COALESCE(NULLIF($1, ''), t.title), details = COALESCE(NULLIF($2, ''), t.details), status_id = $3, priority = COALESCE($4, t.priority), due_date = $5, updated_at = CURRENT_TIMESTAMP FROM projects p WHERE t.id = $6::uuid AND t.deleted_at IS NULL AND t.project_id = p.id AND p.deleted_at IS NULL RETURNING t.id, t.project_id, t.title, t.details, t.status_id, t.priority, t.due_date, t.created_at, t.updated_at`

	var statusParam interface{}
	if payload.StatusID == "" {
		statusParam = nil
	} else {
		statusParam = payload.StatusID
	}

	var dueDateParam interface{}
	if payload.DueDate == nil {
		dueDateParam = nil
	} else {
		dueDateParam = *payload.DueDate
	}

	var statusScan sql.NullString
	var dueDateScan sql.NullTime
	err := tr.pgx.QueryRow(ctx, query, payload.Title, payload.Details, statusParam, payload.Priority, dueDateParam, id).Scan(&t.ID, &t.ProjectID, &t.Title, &t.Details, &statusScan, &t.Priority, &dueDateScan, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.TaskModel{}, huma.Error404NotFound("No task found")
		}
		return models.TaskModel{}, huma.Error400BadRequest("Unable to update task", err)
	}
	if statusScan.Valid {
		t.StatusID = statusScan.String
	} else {
		t.StatusID = ""
	}
	if dueDateScan.Valid {
		t.DueDate = &dueDateScan.Time
	}
	return t, nil
}

func (tr TaskRepository) Delete(ctx context.Context, id string) error {
	sql := `UPDATE tasks t SET deleted_at = CURRENT_TIMESTAMP FROM projects p WHERE t.id = $1::uuid AND t.deleted_at IS NULL AND t.project_id = p.id AND p.deleted_at IS NULL`
	cmdTag, err := tr.pgx.Exec(ctx, sql, id)
	if err != nil {
		return huma.Error400BadRequest("Unable to delete task", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return huma.Error404NotFound("No task found")
	}
	return nil
}
