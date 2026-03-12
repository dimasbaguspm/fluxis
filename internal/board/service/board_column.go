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

func (s *Service) GetBoardColumn(ctx context.Context, id pgtype.UUID) (domain.BoardColumnModel, error) {
	col, err := s.Repo.GetBoardColumn(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BoardColumnModel{}, httpx.NotFound("board column not found")
		}
		return domain.BoardColumnModel{}, fmt.Errorf("get board column: %w", err)
	}

	return domain.BoardColumnModel{
		ID:        col.ID,
		BoardID:   col.BoardID,
		Name:      col.Name,
		Position:  col.Position,
		CreatedAt: col.CreatedAt.Time,
		UpdatedAt: col.UpdatedAt.Time,
	}, nil
}

func (s *Service) ListBoardColumns(ctx context.Context, boardID pgtype.UUID) ([]domain.BoardColumnModel, error) {
	if _, err := s.GetBoard(ctx, boardID); err != nil {
		return nil, fmt.Errorf("validate board: %w", err)
	}

	cols, err := s.Repo.ListBoardColumns(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("list board columns: %w", err)
	}

	result := make([]domain.BoardColumnModel, 0, len(cols))
	for _, col := range cols {
		result = append(result, domain.BoardColumnModel{
			ID:        col.ID,
			BoardID:   col.BoardID,
			Name:      col.Name,
			Position:  col.Position,
			CreatedAt: col.CreatedAt.Time,
			UpdatedAt: col.UpdatedAt.Time,
		})
	}

	return result, nil
}

func (s *Service) CreateBoardColumn(ctx context.Context, boardID pgtype.UUID, b domain.BoardColumnCreateModel) (domain.BoardColumnModel, error) {
	if _, err := s.GetBoard(ctx, boardID); err != nil {
		return domain.BoardColumnModel{}, fmt.Errorf("validate board: %w", err)
	}

	col, err := s.Repo.CreateBoardColumn(ctx, repository.CreateBoardColumnParams{
		BoardID:  boardID,
		Name:     b.Name,
		Position: b.Position,
	})
	if err != nil {
		return domain.BoardColumnModel{}, fmt.Errorf("create board column: %w", err)
	}

	return domain.BoardColumnModel{
		ID:        col.ID,
		BoardID:   col.BoardID,
		Name:      col.Name,
		Position:  col.Position,
		CreatedAt: col.CreatedAt.Time,
		UpdatedAt: col.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateBoardColumn(ctx context.Context, boardID, columnID pgtype.UUID, b domain.BoardColumnUpdateModel) (domain.BoardColumnModel, error) {
	col, err := s.GetBoardColumn(ctx, columnID)
	if err != nil {
		return domain.BoardColumnModel{}, err
	}

	if col.BoardID != boardID {
		return domain.BoardColumnModel{}, httpx.NotFound("board column not found in this board")
	}

	if b.Name != "" {
		colUpdated, err := s.Repo.UpdateBoardColumn(ctx, repository.UpdateBoardColumnParams{
			ID:   columnID,
			Name: b.Name,
		})
		if err != nil {
			return domain.BoardColumnModel{}, fmt.Errorf("update board column: %w", err)
		}
		col.Name = colUpdated.Name
		col.UpdatedAt = colUpdated.UpdatedAt.Time
	}

	if b.Position != 0 {
		colUpdated, err := s.Repo.ReorderBoardColumn(ctx, repository.ReorderBoardColumnParams{
			ID:       columnID,
			Position: b.Position,
		})
		if err != nil {
			return domain.BoardColumnModel{}, fmt.Errorf("update board column position: %w", err)
		}
		col.Position = colUpdated.Position
		col.UpdatedAt = colUpdated.UpdatedAt.Time
	}

	return col, nil
}

func (s *Service) DeleteBoardColumn(ctx context.Context, boardID, columnID pgtype.UUID) error {
	col, err := s.GetBoardColumn(ctx, columnID)
	if err != nil {
		return err
	}

	if col.BoardID != boardID {
		return httpx.NotFound("board column not found in this board")
	}

	_, err = s.Repo.DeleteBoardColumn(ctx, columnID)
	if err != nil {
		return fmt.Errorf("delete board column: %w", err)
	}

	return nil
}
