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
VALUES ($1, $2, (SELECT COALESCE(MAX(position), -1) + 1 FROM board_columns WHERE board_id = $1 AND deleted_at IS NULL))
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

-- name: CountBoardColumns :one
SELECT COUNT(*) as count FROM board_columns
WHERE board_id = $1 AND deleted_at IS NULL;

-- name: CheckBoardColumnsExist :one
SELECT COUNT(*) as count FROM board_columns
WHERE board_id = $1 AND id = ANY($2::uuid[]) AND deleted_at IS NULL;

-- name: ReorderBoardColumnsInBatch :many
UPDATE board_columns
SET position = positions.pos, updated_at = NOW()
FROM (
  SELECT id, ROW_NUMBER() OVER () - 1 as pos
  FROM UNNEST($2::uuid[]) AS t(id)
) AS positions
WHERE board_columns.id = positions.id
  AND board_columns.board_id = $1
  AND board_columns.deleted_at IS NULL
RETURNING board_columns.id, board_columns.board_id, board_columns.name, board_columns.position, board_columns.created_at, board_columns.updated_at, board_columns.deleted_at;
