CREATE TABLE
    IF NOT EXISTS boards (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        sprint_id UUID NOT NULL REFERENCES sprints (id) ON DELETE CASCADE,
        name VARCHAR(100) NOT NULL,
        position INT NOT NULL DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ
    );

CREATE INDEX idx_boards_sprint_id ON boards(sprint_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_boards_position ON boards(sprint_id, position);

CREATE TABLE
    IF NOT EXISTS board_columns (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        board_id UUID NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
        name VARCHAR(100) NOT NULL,
        position INT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ
    );

CREATE INDEX idx_board_columns_board_id ON board_columns(board_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_board_columns_position ON board_columns(board_id, position);