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

-- name: SearchOrganisations :many
-- Searches organisations with pagination support
-- Parameters: $1=idArray, $2=nameArray, $3=sortBy (name/createdAt/updatedAt), $4=sortOrder (asc/desc), $5=pageSize, $6=pageNumber
-- Defaults should be applied in service layer: sortBy=updatedAt, sortOrder=desc, pageSize=25, pageNumber=1
WITH filtered_orgs AS (
  SELECT
    id, name, slug, created_at, updated_at,
    COUNT(*) OVER () as total_count
  FROM orgs
  WHERE
    deleted_at IS NULL
    AND (array_length($1::uuid[], 1) IS NULL OR id = ANY($1::uuid[]))
    AND (array_length($2::text[], 1) IS NULL OR name ILIKE ANY((SELECT '%' || unnest($2::text[]) || '%')))
)
SELECT
    id, name, slug, created_at, updated_at, total_count
FROM
    filtered_orgs
ORDER BY
    CASE WHEN $3 = 'name' AND $4 = 'asc' THEN name END ASC,
    CASE WHEN $3 = 'name' AND $4 = 'desc' THEN name END DESC,
    CASE WHEN $3 = 'createdAt' AND $4 = 'asc' THEN created_at END ASC,
    CASE WHEN $3 = 'createdAt' AND $4 = 'desc' THEN created_at END DESC,
    CASE WHEN $3 = 'updatedAt' AND $4 = 'asc' THEN updated_at END ASC,
    CASE WHEN $3 = 'updatedAt' AND $4 = 'desc' THEN updated_at END DESC
LIMIT $5
OFFSET (($6 - 1) * $5);

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
WITH filtered_members AS (
  SELECT
    om.org_id, om.user_id, om.role, om.joined_at,
    u.email, u.display_name,
    COUNT(*) OVER () as total_count
  FROM
    org_members om
    JOIN users u ON u.id = om.user_id
  WHERE
    om.org_id = $1
    AND ($2::text = '' OR u.email ILIKE '%' || $2 || '%')
    AND ($3::text = '' OR u.display_name ILIKE '%' || $3 || '%')
)
SELECT
    org_id, user_id, role, joined_at, email, display_name, total_count
FROM
    filtered_members
ORDER BY
    joined_at DESC
LIMIT $4
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
