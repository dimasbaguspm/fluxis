-- name: CreateSprint :one
INSERT INTO sprints (project_id, name, goal, status, planned_started_at, planned_completed_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, project_id, name, goal, status, planned_started_at, planned_completed_at, started_at, completed_at, created_at, updated_at, deleted_at;

-- name: GetSprint :one
SELECT id, project_id, name, goal, status, planned_started_at, planned_completed_at, started_at, completed_at, created_at, updated_at, deleted_at
FROM sprints
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListSprintsByProject :many
SELECT id, project_id, name, goal, status, planned_started_at, planned_completed_at, started_at, completed_at, created_at, updated_at, deleted_at
FROM sprints
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListSprintsPaged :many
WITH filtered_sprints AS (
  SELECT
    id, project_id, name, goal, status, planned_started_at, planned_completed_at, started_at, completed_at, created_at, updated_at, deleted_at,
    COUNT(*) OVER () as total_count
  FROM
    sprints
  WHERE
    deleted_at IS NULL
    AND (array_length($1::uuid[], 1) IS NULL OR id = ANY($1::uuid[]))
    AND (array_length($2::uuid[], 1) IS NULL OR project_id = ANY($2::uuid[]))
    AND ($3::text = '' OR name ILIKE '%' || $3 || '%')
)
SELECT
  id, project_id, name, goal, status, planned_started_at, planned_completed_at, started_at, completed_at, created_at, updated_at, deleted_at, total_count
FROM
  filtered_sprints
ORDER BY
  created_at DESC
LIMIT $4
OFFSET $5;

-- name: UpdateSprint :one
UPDATE sprints
SET name = $2, goal = $3, status = $4, planned_started_at = $5, planned_completed_at = $6, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, project_id, name, goal, status, planned_started_at, planned_completed_at, started_at, completed_at, created_at, updated_at, deleted_at;

-- name: StartSprint :one
UPDATE sprints
SET status = 'active', started_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, project_id, name, goal, status, planned_started_at, planned_completed_at, started_at, completed_at, created_at, updated_at, deleted_at;

-- name: CompleteSprint :one
UPDATE sprints
SET status = 'completed', completed_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, project_id, name, goal, status, planned_started_at, planned_completed_at, started_at, completed_at, created_at, updated_at, deleted_at;

-- name: DeleteSprint :one
UPDATE sprints
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, project_id, name, goal, status, planned_started_at, planned_completed_at, started_at, completed_at, created_at, updated_at, deleted_at;

-- name: HardDeleteSprint :exec
DELETE FROM sprints
WHERE id = $1;
