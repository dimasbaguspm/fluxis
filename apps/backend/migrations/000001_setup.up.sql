CREATE TYPE project_status AS ENUM ('active', 'paused', 'archived');

CREATE TABLE
    projects (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name TEXT NOT NULL,
        description TEXT,
        status project_status NOT NULL DEFAULT 'active',
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ
    );

CREATE TABLE statuses (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        slug TEXT NOT NULL,
        position INT NOT NULL DEFAULT 0,
        is_default BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ
);

CREATE OR REPLACE FUNCTION create_default_statuses_for_project()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO statuses (project_id, name, slug, position, is_default)
    VALUES (NEW.id, 'Todo', 'todo', 0, true),
                 (NEW.id, 'In Progress', 'in_progress', 1, false),
                 (NEW.id, 'Done', 'done', 2, false);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_default_statuses
AFTER INSERT ON projects
FOR EACH ROW
EXECUTE FUNCTION create_default_statuses_for_project();

-- Task priority: integer for ordering (higher = higher priority)
-- Use integers for flexible ordering; default 1 (medium)
CREATE TABLE
    tasks (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        project_id UUID REFERENCES projects (id) ON DELETE CASCADE,
        title TEXT NOT NULL,
        details TEXT,
        -- `status_id` references a row in `statuses` so kanban columns are per-project
        status_id UUID REFERENCES statuses(id) ON DELETE SET NULL,
        priority INT NOT NULL DEFAULT 1 CHECK (priority >= 0),
        due_date TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ
    );

-- Activity log: can optionally link to a task (task_id nullable)
CREATE TABLE
    logs (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        project_id UUID REFERENCES projects (id) ON DELETE CASCADE,
        task_id UUID REFERENCES tasks (id) ON DELETE SET NULL,
        status_id UUID REFERENCES statuses(id) ON DELETE SET NULL,
        entry TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );