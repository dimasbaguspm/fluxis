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

func (s *Service) listBoardColumnsUnpaginated(ctx context.Context, boardID pgtype.UUID) ([]domain.BoardColumnModel, error) {
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

func (s *Service) ListBoardColumns(ctx context.Context, q domain.BoardColumnsSearchModel) (domain.BoardColumnsPagedModel, error) {
	q.ApplyDefaults()

	if _, err := s.GetBoard(ctx, q.BoardID); err != nil {
		return domain.BoardColumnsPagedModel{}, fmt.Errorf("validate board: %w", err)
	}

	offset := int32((q.PageNumber - 1) * q.PageSize)
	rows, err := s.Repo.ListBoardColumnsPaged(ctx, repository.ListBoardColumnsPagedParams{
		BoardID: q.BoardID,
		Column2: q.Name,
		Limit:   int32(q.PageSize),
		Offset:  offset,
	})
	if err != nil {
		return domain.BoardColumnsPagedModel{}, fmt.Errorf("list board columns: %w", err)
	}

	if len(rows) == 0 {
		return domain.BoardColumnsPagedModel{}.Empty(q.PageNumber, q.PageSize), nil
	}

	totalCount := int(rows[0].TotalCount)
	totalPages := (totalCount + q.PageSize - 1) / q.PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	items := make([]domain.BoardColumnModel, len(rows))
	for i, row := range rows {
		items[i] = domain.BoardColumnModel{
			ID:        row.ID,
			BoardID:   row.BoardID,
			Name:      row.Name,
			Position:  row.Position,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}

	return domain.BoardColumnsPagedModel{
		Items:      items,
		TotalCount: totalCount,
		TotalPages: totalPages,
		PageNumber: q.PageNumber,
		PageSize:   q.PageSize,
	}, nil
}

func (s *Service) CreateBoardColumn(ctx context.Context, boardID pgtype.UUID, b domain.BoardColumnCreateModel) (domain.BoardColumnModel, error) {
	if _, err := s.GetBoard(ctx, boardID); err != nil {
		return domain.BoardColumnModel{}, fmt.Errorf("validate board: %w", err)
	}

	col, err := s.Repo.CreateBoardColumn(ctx, repository.CreateBoardColumnParams{
		BoardID: boardID,
		Name:    b.Name,
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

	colUpdated, err := s.Repo.UpdateBoardColumn(ctx, repository.UpdateBoardColumnParams{
		ID:   columnID,
		Name: b.Name,
	})
	if err != nil {
		return domain.BoardColumnModel{}, fmt.Errorf("update board column: %w", err)
	}

	return domain.BoardColumnModel{
		ID:        colUpdated.ID,
		BoardID:   colUpdated.BoardID,
		Name:      colUpdated.Name,
		Position:  colUpdated.Position,
		CreatedAt: colUpdated.CreatedAt.Time,
		UpdatedAt: colUpdated.UpdatedAt.Time,
	}, nil
}

func (s *Service) ReorderBoardColumns(ctx context.Context, boardID pgtype.UUID, reorder domain.BoardColumnReorderModel) ([]domain.BoardColumnModel, error) {
	if _, err := s.GetBoard(ctx, boardID); err != nil {
		return nil, fmt.Errorf("validate board: %w", err)
	}

	cols, err := s.Repo.ReorderBoardColumnsInBatch(ctx, repository.ReorderBoardColumnsInBatchParams{
		BoardID: boardID,
		Column2: reorder,
	})
	if err != nil {
		return nil, fmt.Errorf("reorder board columns: %w", err)
	}

	if len(cols) == 0 {
		if len(reorder) == 0 {
			return nil, httpx.BadRequest("columns array is required and cannot be empty")
		}
		return nil, httpx.BadRequest("some board columns not found or don't belong to this board, or reorder array must include all board columns")
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
