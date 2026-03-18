-- Remove position column from boards (no longer used with one-board-per-sprint constraint)
DROP INDEX IF EXISTS idx_boards_position;
ALTER TABLE boards DROP COLUMN IF EXISTS position;
