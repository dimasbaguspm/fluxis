-- name: GetUser :one
SELECT
    id, email, display_name, created_at, updated_at
FROM
    users
WHERE
    id = $1
    AND deleted_at IS NULL
LIMIT
    1;

-- name: GetUserByEmail :one
SELECT
    id, email, display_name, created_at, updated_at
FROM
    users
WHERE
    email = $1
    AND deleted_at IS NULL
LIMIT
    1;

-- name: ListUsers :many
SELECT
    id, email, display_name, created_at, updated_at
FROM
    users
WHERE
    deleted_at IS NULL
    AND ($1::text = '' OR email ILIKE '%' || $1 || '%')
    AND ($2::text = '' OR display_name ILIKE '%' || $2 || '%')
ORDER BY
    created_at DESC
LIMIT  $3
OFFSET $4;

-- name: CountUsers :one
SELECT
    COUNT(*)
FROM
    users
WHERE
    deleted_at IS NULL
    AND ($1::text = '' OR email ILIKE '%' || $1 || '%')
    AND ($2::text = '' OR display_name ILIKE '%' || $2 || '%');

-- name: CreateUser :one
INSERT INTO
    users (email, display_name, password_hash)
VALUES
    ($1, $2, $3)
RETURNING
    id, email, display_name, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users
SET
    display_name = COALESCE(NULLIF($1, ''), display_name),
    password_hash = COALESCE(NULLIF($2, ''), password_hash),
    updated_at = NOW()
WHERE
    id = $3
    AND deleted_at IS NULL
RETURNING
    id, email, display_name, created_at, updated_at;

-- name: DeleteUser :exec
UPDATE users
SET
    deleted_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;