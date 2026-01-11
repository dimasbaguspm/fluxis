-- Down migration for 000001_setup.up.sql
-- Reverses created triggers, functions, tables, and types.
BEGIN;

-- Remove trigger that creates default statuses for new projects
DROP TRIGGER IF EXISTS trg_create_default_statuses ON projects;

-- Remove trigger function
DROP FUNCTION IF EXISTS create_default_statuses_for_project () CASCADE;

-- Drop activity logs first (depends on tasks/statuses/projects)
DROP TABLE IF EXISTS logs;

-- Drop tasks (depends on statuses/projects)
DROP TABLE IF EXISTS tasks;

-- Drop statuses (depends on projects)
DROP TABLE IF EXISTS statuses;

-- Drop projects
DROP TABLE IF EXISTS projects;

-- Drop project status enum
DROP TYPE IF EXISTS project_status;

COMMIT;