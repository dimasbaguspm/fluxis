-- Restore position column to boards (rollback)
ALTER TABLE boards ADD COLUMN position INT NOT NULL DEFAULT 0;
CREATE INDEX idx_boards_position ON boards(sprint_id, position);
