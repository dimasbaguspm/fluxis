package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dimasbaguspm/fluxis/internal/board/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
	"github.com/dimasbaguspm/fluxis/pkg/syncx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBoardNotFound           = httpx.NotFound("board not found")
	defaultBoardColumnNames    = []string{
		"To Do",
		"In Progress",
		"Code Review",
		"Ready to Test",
		"Ready for PO Review",
		"Done",
	}
)

func toBoardModel(board repository.Board) domain.BoardModel {
	return domain.BoardModel{
		ID:        board.ID,
		SprintID:  board.SprintID,
		Name:      board.Name,
		CreatedAt: board.CreatedAt.Time,
		UpdatedAt: board.UpdatedAt.Time,
	}
}

func (s *Service) seedDefaultColumns(ctx context.Context, boardID pgtype.UUID) ([]domain.BoardColumnModel, error) {
	positions := make([]int32, len(defaultBoardColumnNames))
	for i := range positions {
		positions[i] = int32(i)
	}

	cols, err := s.Repo.BulkCreateBoardColumns(ctx, repository.BulkCreateBoardColumnsParams{
		Column1: boardID,
		Column2: defaultBoardColumnNames,
		Column3: positions,
	})
	if err != nil {
		return nil, fmt.Errorf("seed default columns: %w", err)
	}

	result := make([]domain.BoardColumnModel, len(cols))
	for i, col := range cols {
		result[i] = domain.BoardColumnModel{
			ID:        col.ID,
			BoardID:   col.BoardID,
			Name:      col.Name,
			Position:  col.Position,
			CreatedAt: col.CreatedAt.Time,
			UpdatedAt: col.UpdatedAt.Time,
		}
	}
	return result, nil
}

func (s *Service) CreateBoard(ctx context.Context, b domain.BoardCreateModel) (domain.BoardModel, error) {
	err := syncx.Run(ctx,
		func(ctx context.Context) error {
			_, err := s.Sprint.GetSprint(ctx, b.SprintID)
			return err // propagates ErrSprintNotFound from sprint service
		},
		func(ctx context.Context) error {
			boards, err := s.Repo.ListBoardsBySprint(ctx, b.SprintID)
			if err != nil {
				return fmt.Errorf("list boards: %w", err)
			}
			if len(boards) > 0 {
				return httpx.Conflict("sprint already has a board")
			}
			return nil
		},
	)
	if err != nil {
		return domain.BoardModel{}, err
	}

	board, err := s.Repo.CreateBoard(ctx, repository.CreateBoardParams{
		SprintID: b.SprintID,
		Name:     b.Name,
	})
	if err != nil {
		return domain.BoardModel{}, fmt.Errorf("create board: %w", err)
	}

	if _, err := s.seedDefaultColumns(ctx, board.ID); err != nil {
		return domain.BoardModel{}, err
	}

	result := toBoardModel(board)
	if err := s.Bus.Publish(ctx, pubsub.BoardCreated, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.BoardCreated), "error", err)
	}

	return result, nil
}

func (s *Service) GetBoard(ctx context.Context, id pgtype.UUID) (domain.BoardModel, error) {
	board, err := s.Repo.GetBoard(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BoardModel{}, ErrBoardNotFound
		}
		return domain.BoardModel{}, fmt.Errorf("get board: %w", err)
	}

	return toBoardModel(board), nil
}

func (s *Service) ListBoards(ctx context.Context, q domain.BoardsSearchModel) (domain.BoardsPagedModel, error) {
	q.ApplyDefaults()

	offset := int32((q.PageNumber - 1) * q.PageSize)
	rows, err := s.Repo.ListBoardsPaged(ctx, repository.ListBoardsPagedParams{
		Column1: q.ID,
		Column2: q.ProjectID,
		Column3: q.SprintID,
		Column4: q.Name,
		Limit:   int32(q.PageSize),
		Offset:  offset,
	})

	if err != nil {
		return domain.BoardsPagedModel{}, fmt.Errorf("list boards: %w", err)
	}

	if len(rows) == 0 {
		return domain.BoardsPagedModel{}.Empty(q.PageNumber, q.PageSize), nil
	}

	totalCount := int(rows[0].TotalCount)
	totalPages := (totalCount + q.PageSize - 1) / q.PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	items := make([]domain.BoardModel, len(rows))
	for i, row := range rows {
		items[i] = domain.BoardModel{
			ID:        row.ID,
			SprintID:  row.SprintID,
			Name:      row.Name,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}

	return domain.BoardsPagedModel{
		Items:      items,
		TotalCount: totalCount,
		TotalPages: totalPages,
		PageNumber: q.PageNumber,
		PageSize:   q.PageSize,
	}, nil
}

func (s *Service) UpdateBoard(ctx context.Context, id pgtype.UUID, b domain.BoardUpdateModel) (domain.BoardModel, error) {
	var existing domain.BoardModel
	var sprint domain.SprintModel

	err := syncx.Run(ctx,
		func(ctx context.Context) error {
			b, err := s.GetBoard(ctx, id)
			if err != nil {
				return fmt.Errorf("validate board: %w", err)
			}
			existing = b
			return nil
		},
		func(ctx context.Context) error {
			if b.SprintID.Valid {
				s, err := s.Sprint.GetSprint(ctx, b.SprintID)
				if err != nil {
					return fmt.Errorf("validate sprint: %w", err)
				}
				sprint = s
			}
			return nil
		},
	)

	if err != nil {
		return domain.BoardModel{}, err
	}

	// Use existing value if not provided
	name := existing.Name
	if b.Name != "" {
		name = b.Name
	}

	sprintID := existing.SprintID
	if b.SprintID.Valid {
		sprintID = sprint.ID
	}

	board, err := s.Repo.UpdateBoard(ctx, repository.UpdateBoardParams{
		ID:       id,
		Name:     name,
		SprintID: sprintID,
	})

	if err != nil {
		return domain.BoardModel{}, fmt.Errorf("update board: %w", err)
	}

	result := toBoardModel(board)
	if err := s.Bus.Publish(ctx, pubsub.BoardUpdated, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.BoardUpdated), "error", err)
	}

	return result, nil
}

func (s *Service) DeleteBoard(ctx context.Context, id pgtype.UUID) error {
	_, err := s.Repo.DeleteBoard(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBoardNotFound
		}
		return fmt.Errorf("delete board: %w", err)
	}

	if err := s.Bus.Publish(ctx, pubsub.BoardDeleted, map[string]string{"id": uuid.UUID(id.Bytes).String()}); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.BoardDeleted), "error", err)
	}

	return nil
}
