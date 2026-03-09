-- name: GenerateTicketKey :one
SELECT generate_ticket_key($1);

-- name: CreateTicket :one
INSERT INTO tickets (
    project_id,
    ticket_number,
    key,
    type,
    priority,
    title,
    description,
    reporter_id,
    assignee_id,
    story_points,
    due_date
)
VALUES (
    $1,
    (SELECT next_number - 1 FROM ticket_counters WHERE project_id = $1),
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
)
RETURNING id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at;

-- name: GetTicket :one
SELECT id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at
FROM tickets
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetTicketByKey :one
SELECT id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at
FROM tickets
WHERE project_id = $1 AND key = $2 AND deleted_at IS NULL;

-- name: ListTicketsByProject :many
SELECT id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at
FROM tickets
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY ticket_number DESC;

-- name: ListTicketsBySprint :many
SELECT id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at
FROM tickets
WHERE project_id = $1 AND sprint_id = $2 AND deleted_at IS NULL
ORDER BY ticket_number DESC;

-- name: ListTicketsByBoard :many
SELECT id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at
FROM tickets
WHERE board_id = $1 AND deleted_at IS NULL
ORDER BY ticket_number DESC;

-- name: ListTicketsByBoardColumn :many
SELECT id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at
FROM tickets
WHERE board_column_id = $1 AND deleted_at IS NULL
ORDER BY ticket_number DESC;

-- name: UpdateTicketBoard :one
UPDATE tickets
SET board_id = $2, board_column_id = $3, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at;

-- name: UpdateTicketSprint :one
UPDATE tickets
SET sprint_id = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at;

-- name: UpdateTicketDetails :one
UPDATE tickets
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    type = COALESCE($4, type),
    priority = COALESCE($5, priority),
    assignee_id = COALESCE($6, assignee_id),
    story_points = COALESCE($7, story_points),
    due_date = COALESCE($8, due_date),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at;

-- name: DeleteTicket :one
UPDATE tickets
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, project_id, ticket_number, key, sprint_id, board_id, board_column_id, type, priority, title, description, assignee_id, reporter_id, epic_id, parent_id, story_points, due_date, created_at, updated_at, deleted_at;

-- name: HardDeleteTicket :exec
DELETE FROM tickets
WHERE id = $1;
