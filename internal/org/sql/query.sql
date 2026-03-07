-- name: CreateOrg :one
INSERT INTO
    orgs (name, slug)
VALUES
    ($1, $2)
RETURNING
    id, name, slug, created_at, updated_at;

-- name: GetOrgById :one
SELECT
    id, name, slug, created_at, updated_at
FROM
    orgs
WHERE
    id = $1
    AND deleted_at IS NULL
LIMIT
    1;

-- name: GetOrgBySlug :one
SELECT
    id, name, slug, created_at, updated_at
FROM
    orgs
WHERE
    slug = $1
    AND deleted_at IS NULL
LIMIT
    1;

-- name: SlugExists :one
SELECT EXISTS (
    SELECT 1 FROM orgs WHERE slug = $1 AND deleted_at IS NULL
);

-- name: UpdateOrg :one
UPDATE orgs
SET
    name = COALESCE(NULLIF($1, ''), name),
    slug = COALESCE(NULLIF($2, ''), slug),
    updated_at = NOW()
WHERE
    id = $3
    AND deleted_at IS NULL
RETURNING
    id, name, slug, created_at, updated_at;

-- name: DeleteOrg :exec
UPDATE orgs
SET
    deleted_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: ListOrg :many
SELECT
    id, name, slug, created_at, updated_at
FROM
    orgs
WHERE
    deleted_at IS NULL
    AND ($1::uuid IS NULL OR id = $1)
    AND ($2::uuid IS NULL OR id IN (
        SELECT org_id FROM org_members WHERE user_id = $2
    ))
ORDER BY
    created_at DESC;

-- name: CreateOrgMember :one
INSERT INTO
    org_members (org_id, user_id, role)
VALUES
    ($1, $2, $3)
RETURNING
    org_id, user_id, role, joined_at;

-- name: GetOrgMember :one
SELECT
    om.org_id, om.user_id, om.role, om.joined_at,
    u.email, u.display_name
FROM
    org_members om
    JOIN users u ON u.id = om.user_id
WHERE
    om.org_id = $1
    AND om.user_id = $2
LIMIT
    1;

-- name: UpdateOrgMemberRole :one
UPDATE org_members
SET
    role = $3
WHERE
    org_id = $1
    AND user_id = $2
RETURNING
    org_id, user_id, role, joined_at;

-- name: DeleteOrgMember :exec
DELETE FROM org_members
WHERE
    org_id = $1
    AND user_id = $2;

-- name: ListOrgMembers :many
SELECT
    om.org_id, om.user_id, om.role, om.joined_at,
    u.email, u.display_name
FROM
    org_members om
    JOIN users u ON u.id = om.user_id
WHERE
    om.org_id = $1
    AND ($2::text = '' OR u.email ILIKE '%' || $2 || '%')
    AND ($3::text = '' OR u.display_name ILIKE '%' || $3 || '%')
ORDER BY
    om.joined_at DESC
LIMIT  $4
OFFSET $5;

-- name: CountOrgMembers :one
SELECT
    COUNT(*)
FROM
    org_members om
    JOIN users u ON u.id = om.user_id
WHERE
    om.org_id = $1
    AND ($2::text = '' OR u.email ILIKE '%' || $2 || '%')
    AND ($3::text = '' OR u.display_name ILIKE '%' || $3 || '%');
