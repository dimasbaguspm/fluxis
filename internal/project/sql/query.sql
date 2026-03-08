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
