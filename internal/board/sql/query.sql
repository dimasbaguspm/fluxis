-- name: CreateBoard :one
INSERT INTO boards (sprint_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetBoard :one
SELECT * FROM boards WHERE id = $1 AND deleted_at IS NULL;

-- name: ListBoardsBySprint :many
SELECT * FROM boards WHERE sprint_id = $1 AND deleted_at IS NULL;

-- name: ListBoardsPaged :many
WITH applicable_sprints AS (
  -- If projectId is provided, get all sprints for those projects
  -- If projectId is not provided, this CTE returns all non-deleted sprints
  SELECT id
  FROM sprints
  WHERE
    deleted_at IS NULL
    AND (array_length($2::uuid[], 1) IS NULL OR project_id = ANY($2::uuid[]))
),
filtered_boards AS (
  SELECT
    b.id, b.sprint_id, b.name, b.created_at, b.updated_at, b.deleted_at,
    COUNT(*) OVER () as total_count
  FROM
    boards b
  WHERE
    b.deleted_at IS NULL
    -- Filter by board IDs if provided
    AND (array_length($1::uuid[], 1) IS NULL OR b.id = ANY($1::uuid[]))
    -- Filter by sprint ID if provided, otherwise use applicable sprints from project filter
    AND (
      CASE
        WHEN array_length($3::uuid[], 1) IS NOT NULL THEN b.sprint_id = ANY($3::uuid[])
        ELSE b.sprint_id IN (SELECT id FROM applicable_sprints)
      END
    )
    -- Filter by name if provided
    AND ($4::text = '' OR b.name ILIKE '%' || $4 || '%')
)
SELECT
  id, sprint_id, name, created_at, updated_at, deleted_at, total_count
FROM
  filtered_boards
LIMIT $5
OFFSET $6;

-- name: UpdateBoard :one
UPDATE boards
SET name = $2, sprint_id = $3, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

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

-- name: ListBoardColumnsPaged :many
WITH filtered_columns AS (
  SELECT
    id, board_id, name, position, created_at, updated_at, deleted_at,
    COUNT(*) OVER () as total_count
  FROM
    board_columns
  WHERE
    deleted_at IS NULL
    AND (array_length($1::uuid[], 1) IS NULL OR id = ANY($1::uuid[]))
    AND (array_length($2::uuid[], 1) IS NULL OR board_id = ANY($2::uuid[]))
    AND ($3::text = '' OR name ILIKE '%' || $3 || '%')
)
SELECT
  id, board_id, name, position, created_at, updated_at, deleted_at, total_count
FROM
  filtered_columns
ORDER BY
  position ASC
LIMIT $4
OFFSET $5;

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
