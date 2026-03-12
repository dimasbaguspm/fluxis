-- name: CreateProject :one
INSERT INTO projects (org_id, key, name, description, visibility)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, org_id, key, name, description, visibility, created_at, updated_at, deleted_at;

-- name: GetProject :one
SELECT id, org_id, key, name, description, visibility, created_at, updated_at, deleted_at
FROM projects
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetProjectByKey :one
SELECT id, org_id, key, name, description, visibility, created_at, updated_at, deleted_at
FROM projects
WHERE org_id = $1 AND key = $2 AND deleted_at IS NULL;

-- name: ListProjectsByOrg :many
SELECT id, org_id, key, name, description, visibility, created_at, updated_at, deleted_at
FROM projects
WHERE org_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListProjectsByOrgPaged :many
WITH filtered_projects AS (
  SELECT
    id, org_id, key, name, description, visibility, created_at, updated_at, deleted_at,
    COUNT(*) OVER () as total_count
  FROM
    projects
  WHERE
    deleted_at IS NULL
    AND (array_length($1::uuid[], 1) IS NULL OR org_id = ANY($1::uuid[]))
    AND (array_length($2::uuid[], 1) IS NULL OR id = ANY($2::uuid[]))
    AND ($3::text = '' OR name ILIKE '%' || $3 || '%')
)
SELECT
  id, org_id, key, name, description, visibility, created_at, updated_at, deleted_at, total_count
FROM
  filtered_projects
ORDER BY
  created_at DESC
LIMIT $4
OFFSET $5;

-- name: UpdateProject :one
UPDATE projects
SET name = $2, description = $3, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, org_id, key, name, description, visibility, created_at, updated_at, deleted_at;

-- name: UpdateProjectVisibility :one
UPDATE projects
SET visibility = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, org_id, key, name, description, visibility, created_at, updated_at, deleted_at;

-- name: DeleteProject :one
UPDATE projects
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, org_id, key, name, description, visibility, created_at, updated_at, deleted_at;

-- name: HardDeleteProject :exec
DELETE FROM projects
WHERE id = $1;
