package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dimasbaguspm/fluxis/internal/board/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
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

func (s *Service) CreateBoard(ctx context.Context, sprintID pgtype.UUID, b domain.BoardCreateModel) (domain.BoardModel, error) {
	if b.Name == nil {
		return domain.BoardModel{}, httpx.BadRequest("name is required")
	}

	board, err := s.repo.CreateBoard(ctx, repository.CreateBoardParams{
		SprintID: sprintID,
		Name:     *b.Name,
	})
	if err != nil {
		return domain.BoardModel{}, fmt.Errorf("create board: %w", err)
	}

	return toBoardModel(board), nil
}

func (s *Service) GetBoard(ctx context.Context, id pgtype.UUID) (domain.BoardModel, error) {
	board, err := s.repo.GetBoard(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BoardModel{}, ErrBoardNotFound
		}
		return domain.BoardModel{}, fmt.Errorf("get board: %w", err)
	}

	return toBoardModel(board), nil
}

func (s *Service) ListBoardsBySprint(ctx context.Context, sprintID pgtype.UUID) ([]domain.BoardModel, error) {
	boards, err := s.repo.ListBoardsBySprint(ctx, sprintID)
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
	// Get existing board first
	existing, err := s.repo.GetBoard(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BoardModel{}, ErrBoardNotFound
		}
		return domain.BoardModel{}, fmt.Errorf("get board: %w", err)
	}

	// Use existing value if not provided
	name := existing.Name
	if b.Name != nil {
		name = *b.Name
	}

	board, err := s.repo.UpdateBoard(ctx, repository.UpdateBoardParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return domain.BoardModel{}, fmt.Errorf("update board: %w", err)
	}

	return toBoardModel(board), nil
}

func (s *Service) ReorderBoard(ctx context.Context, id pgtype.UUID, position int32) (domain.BoardModel, error) {
	board, err := s.repo.ReorderBoard(ctx, repository.ReorderBoardParams{
		ID:       id,
		Position: position,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BoardModel{}, ErrBoardNotFound
		}
		return domain.BoardModel{}, fmt.Errorf("reorder board: %w", err)
	}

	return toBoardModel(board), nil
}

func (s *Service) DeleteBoard(ctx context.Context, id pgtype.UUID) error {
	_, err := s.repo.DeleteBoard(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBoardNotFound
		}
		return fmt.Errorf("delete board: %w", err)
	}

	return nil
}
