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

-- name: DeleteBoard :one
UPDATE boards SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *;

-- name: ReorderBoardsInBatch :many
-- Atomically validates and reorders boards within a sprint with row-level locking
-- All boards must belong to the same sprint; results ordered by position to maintain input array order
WITH validation AS (
  -- Validate: all provided IDs exist and belong to this sprint
  SELECT id, ROW_NUMBER() OVER () - 1 as pos
  FROM UNNEST($2::uuid[]) AS t(id)
  WHERE EXISTS (
    SELECT 1 FROM boards b
    WHERE b.id = t.id AND b.sprint_id = $1 AND b.deleted_at IS NULL
  )
),updated AS (
  UPDATE boards
  SET position = validation.pos, updated_at = NOW()
  FROM validation
  WHERE boards.id = validation.id
    AND boards.sprint_id = $1
    AND boards.deleted_at IS NULL
    -- Validate: provided count matches total boards in sprint
    AND (
      SELECT COUNT(*) FROM boards
      WHERE sprint_id = $1 AND deleted_at IS NULL
    ) = array_length($2::uuid[], 1)
    -- Validate: all array elements are valid boards for this sprint
    AND (
      SELECT COUNT(*) FROM validation
    ) = array_length($2::uuid[], 1)
  RETURNING boards.id, boards.sprint_id, boards.name, boards.position, boards.created_at, boards.updated_at, boards.deleted_at
)
SELECT * FROM updated ORDER BY position;

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

-- name: ReorderBoardColumnsInBatch :many
-- Atomically validates and reorders columns with row-level locking
-- Results ordered by position to maintain input array order
WITH validation AS (
  -- Validate: all provided IDs exist and belong to this board
  SELECT id, ROW_NUMBER() OVER () - 1 as pos
  FROM UNNEST($2::uuid[]) AS t(id)
  WHERE EXISTS (
    SELECT 1 FROM board_columns bc
    WHERE bc.id = t.id AND bc.board_id = $1 AND bc.deleted_at IS NULL
  )
),updated AS (
  UPDATE board_columns
  SET position = validation.pos, updated_at = NOW()
  FROM validation
  WHERE board_columns.id = validation.id
    AND board_columns.board_id = $1
    AND board_columns.deleted_at IS NULL
    -- Validate: provided count matches total columns in board
    AND (
      SELECT COUNT(*) FROM board_columns
      WHERE board_id = $1 AND deleted_at IS NULL
    ) = array_length($2::uuid[], 1)
    -- Validate: all array elements are valid columns for this board
    AND (
      SELECT COUNT(*) FROM validation
    ) = array_length($2::uuid[], 1)
  RETURNING board_columns.id, board_columns.board_id, board_columns.name, board_columns.position, board_columns.created_at, board_columns.updated_at, board_columns.deleted_at
)
SELECT * FROM updated ORDER BY position;
