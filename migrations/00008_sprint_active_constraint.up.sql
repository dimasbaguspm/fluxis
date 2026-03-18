-- Enforce at most one active sprint per project at the database level
CREATE UNIQUE INDEX idx_sprints_one_active_per_project
  ON sprints(project_id)
  WHERE status = 'active' AND deleted_at IS NULL;
