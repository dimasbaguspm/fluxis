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

func (s *Service) CreateBoardColumn(ctx context.Context, boardID pgtype.UUID, b domain.BoardColumnCreateModel) (domain.BoardColumnModel, error) {
	if _, err := s.GetBoard(ctx, boardID); err != nil {
		return domain.BoardColumnModel{}, fmt.Errorf("validate board: %w", err)
	}

	if b.Name == nil {
		return domain.BoardColumnModel{}, httpx.BadRequest("name is required")
	}

	if b.Position == nil {
		return domain.BoardColumnModel{}, httpx.BadRequest("position is required")
	}

	col, err := s.Repo.CreateBoardColumn(ctx, repository.CreateBoardColumnParams{
		BoardID:  boardID,
		Name:     *b.Name,
		Position: *b.Position,
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
