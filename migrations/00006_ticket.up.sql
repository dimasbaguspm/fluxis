CREATE TYPE ticket_type AS ENUM ('bug', 'story', 'task', 'epic');
CREATE TYPE ticket_priority AS ENUM ('low', 'medium', 'high', 'critical');

CREATE TABLE ticket_counters (
    project_id UUID PRIMARY KEY REFERENCES projects (id) ON DELETE CASCADE,
    next_number INT NOT NULL DEFAULT 1
);

CREATE OR REPLACE FUNCTION generate_ticket_key(p_project_id UUID)
RETURNS VARCHAR AS $$
DECLARE
    v_next_number INT;
    v_project_key VARCHAR;
BEGIN
    -- Initialize counter if needed
    INSERT INTO ticket_counters (project_id, next_number)
    VALUES (p_project_id, 1)
    ON CONFLICT (project_id) DO NOTHING;

    -- Increment and get the next number
    UPDATE ticket_counters
    SET next_number = next_number + 1
    WHERE project_id = p_project_id
    RETURNING next_number - 1 INTO v_next_number;

    -- Get the project key
    SELECT key INTO v_project_key
    FROM projects
    WHERE id = p_project_id;

    -- Return the formatted key
    RETURN v_project_key || '-' || v_next_number;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE
   IF NOT EXISTS tickets (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
       project_id UUID NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
       ticket_number INT NOT NULL,
       key VARCHAR(20) UNIQUE NOT NULL,
       sprint_id UUID REFERENCES sprints (id) ON DELETE SET NULL,
       board_id UUID REFERENCES boards (id) ON DELETE SET NULL,
       board_column_id UUID REFERENCES board_columns (id) ON DELETE SET NULL,
       type ticket_type NOT NULL,
       priority ticket_priority NOT NULL,
       title VARCHAR(255) NOT NULL,
       description TEXT,
       assignee_id UUID REFERENCES users (id) ON DELETE SET NULL,
       reporter_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
       epic_id UUID REFERENCES tickets (id) ON DELETE SET NULL,
       parent_id UUID REFERENCES tickets (id) ON DELETE SET NULL,
       story_points INT,
       due_date DATE,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
       deleted_at TIMESTAMPTZ
   );

CREATE INDEX idx_tickets_project_id ON tickets (project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tickets_sprint_id ON tickets (sprint_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tickets_board_id ON tickets (board_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tickets_board_column_id ON tickets (board_column_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tickets_assignee_id ON tickets (assignee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tickets_reporter_id ON tickets (reporter_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tickets_epic_id ON tickets (epic_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tickets_parent_id ON tickets (parent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tickets_deleted_at ON tickets (deleted_at);
