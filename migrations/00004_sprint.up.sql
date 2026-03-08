CREATE TYPE sprint_status AS ENUM ('planned', 'active', 'completed');

CREATE TABLE
    IF NOT EXISTS sprints (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        project_id UUID NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
        name VARCHAR(100) NOT NULL,
        goal TEXT,
        status sprint_status NOT NULL DEFAULT 'planned',
        planned_started_at TIMESTAMPTZ,
        planned_completed_at TIMESTAMPTZ,
        started_at TIMESTAMPTZ,
        completed_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ
    );

CREATE INDEX idx_sprints_project_id ON sprints (project_id);

CREATE INDEX idx_sprints_status ON sprints (status);

CREATE INDEX idx_sprints_deleted_at ON sprints (deleted_at);