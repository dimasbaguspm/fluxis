package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dimasbaguspm/fluxis/internal/board/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/syncx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBoardNotFound = httpx.NotFound("board not found")
)

func toBoardModel(board repository.Board) domain.BoardModel {
	return domain.BoardModel{
		ID:        board.ID,
		SprintID:  board.SprintID,
		Name:      board.Name,
		Position:  board.Position,
		CreatedAt: board.CreatedAt.Time,
		UpdatedAt: board.UpdatedAt.Time,
	}
}

func (s *Service) CreateBoard(ctx context.Context, b domain.BoardCreateModel) (domain.BoardModel, error) {
	sprint, err := s.Sprint.GetSprint(ctx, b.SprintID)
	if err != nil {
		return domain.BoardModel{}, fmt.Errorf("get sprint: %w", err)
	}

	board, err := s.Repo.CreateBoard(ctx, repository.CreateBoardParams{
		SprintID: sprint.ID,
		Name:     b.Name,
	})
	if err != nil {
		return domain.BoardModel{}, fmt.Errorf("create board: %w", err)
	}

	return toBoardModel(board), nil
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

func (s *Service) ListBoardsBySprint(ctx context.Context, sprintID pgtype.UUID) ([]domain.BoardModel, error) {
	boards, err := s.Repo.ListBoardsBySprint(ctx, sprintID)
	if err != nil {
		return nil, fmt.Errorf("list boards: %w", err)
	}

	if boards == nil {
		boards = []repository.Board{}
	}

	result := make([]domain.BoardModel, len(boards))
	for i, board := range boards {
		result[i] = toBoardModel(board)
	}

	return result, nil
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

	return toBoardModel(board), nil
}

func (s *Service) ReorderBoards(ctx context.Context, sprintID pgtype.UUID, reorder domain.BoardReorderModel) ([]domain.BoardModel, error) {
	sprint, err := s.Sprint.GetSprint(ctx, sprintID)
	if err != nil {
		return nil, fmt.Errorf("validate sprint: %w", err)
	}

	boards, err := s.Repo.ReorderBoardsInBatch(ctx, repository.ReorderBoardsInBatchParams{
		SprintID: sprint.ID,
		Column2:  reorder,
	})
	if err != nil {
		return nil, fmt.Errorf("reorder boards: %w", err)
	}

	if len(boards) == 0 {
		if len(reorder) == 0 {
			return nil, httpx.BadRequest("boards array is required and cannot be empty")
		}
		return nil, httpx.BadRequest("some boards not found or don't belong to this sprint, or reorder array must include all boards in the sprint")
	}

	result := make([]domain.BoardModel, 0, len(boards))
	for _, board := range boards {
		result = append(result, domain.BoardModel{
			ID:        board.ID,
			SprintID:  board.SprintID,
			Name:      board.Name,
			Position:  board.Position,
			CreatedAt: board.CreatedAt.Time,
			UpdatedAt: board.UpdatedAt.Time,
		})
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

	return nil
}
