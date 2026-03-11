-- name: CreateBoard :one
INSERT INTO boards (sprint_id, name, position)
VALUES ($1, $2, (SELECT COALESCE(MAX(position), -1) + 1 FROM boards WHERE sprint_id = $1 AND deleted_at IS NULL))
RETURNING *;

-- name: GetBoard :one
SELECT * FROM boards WHERE id = $1 AND deleted_at IS NULL;

-- name: ListBoardsBySprint :many
SELECT * FROM boards WHERE sprint_id = $1 AND deleted_at IS NULL ORDER BY position ASC;

-- name: UpdateBoard :one
UPDATE boards
SET name = $2, sprint_id = $3, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ReorderBoard :one
UPDATE boards SET position = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: DeleteBoard :one
UPDATE boards SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: CreateBoardColumn :one
INSERT INTO board_columns (board_id, name, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBoardColumn :one
SELECT * FROM board_columns WHERE id = $1 AND deleted_at IS NULL;

-- name: ListBoardColumns :many
SELECT * FROM board_columns WHERE board_id = $1 AND deleted_at IS NULL ORDER BY position ASC;

-- name: UpdateBoardColumn :one
UPDATE board_columns SET name = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: ReorderBoardColumn :one
UPDATE board_columns SET position = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: DeleteBoardColumn :one
UPDATE board_columns SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *;
