-- Enforce at most one board per sprint
CREATE UNIQUE INDEX idx_boards_sprint_id_unique
  ON boards(sprint_id)
  WHERE deleted_at IS NULL;
