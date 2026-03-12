package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type BoardModel struct {
	ID        pgtype.UUID `json:"id"`
	SprintID  pgtype.UUID `json:"sprintId"`
	Name      string      `json:"name" validate:"required,min=1"`
	Position  int32       `json:"position"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type BoardCreateModel struct {
	Name     string      `json:"name" validate:"required,min=1"`
	SprintID pgtype.UUID `json:"sprintId" validate:"required"`
}

type BoardUpdateModel struct {
	Name     string      `json:"name,omitempty" validate:"omitempty,min=1"`
	SprintID pgtype.UUID `json:"sprintId,omitempty"`
}

type BoardReorderModel []pgtype.UUID

type BoardsSearchModel struct {
	SprintID   pgtype.UUID `json:"sprintId" validate:"required"`
	Name       string      `json:"name"`
	PageNumber int         `json:"pageNumber" validate:"min=1"`
	PageSize   int         `json:"pageSize" validate:"min=1,max=100"`
}

func (b *BoardsSearchModel) ApplyDefaults() {
	const (
		defaultPageNumber = 1
		defaultPageSize   = 25
	)
	if b.PageNumber == 0 {
		b.PageNumber = defaultPageNumber
	}
	if b.PageSize == 0 {
		b.PageSize = defaultPageSize
	}
}

type BoardsPagedModel struct {
	Items      []BoardModel `json:"items"`
	TotalCount int          `json:"totalCount"`
	TotalPages int          `json:"totalPages"`
	PageNumber int          `json:"pageNumber"`
	PageSize   int          `json:"pageSize"`
}

func (m BoardsPagedModel) Empty(pageNumber, pageSize int) BoardsPagedModel {
	return BoardsPagedModel{
		Items:      []BoardModel{},
		TotalCount: 0,
		TotalPages: 0,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}

type BoardColumnModel struct {
	ID        pgtype.UUID `json:"id"`
	BoardID   pgtype.UUID `json:"boardId"`
	Name      string      `json:"name" validate:"required,min=1"`
	Position  int32       `json:"position"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type BoardColumnCreateModel struct {
	Name string `json:"name" validate:"required,min=1"`
}

type BoardColumnUpdateModel struct {
	Name string `json:"name,omitempty" validate:"omitempty,min=1"`
}

type BoardColumnReorderModel []pgtype.UUID

type BoardColumnsSearchModel struct {
	BoardID    pgtype.UUID `json:"boardId" validate:"required"`
	Name       string      `json:"name"`
	PageNumber int         `json:"pageNumber" validate:"min=1"`
	PageSize   int         `json:"pageSize" validate:"min=1,max=100"`
}

func (b *BoardColumnsSearchModel) ApplyDefaults() {
	const (
		defaultPageNumber = 1
		defaultPageSize   = 25
	)
	if b.PageNumber == 0 {
		b.PageNumber = defaultPageNumber
	}
	if b.PageSize == 0 {
		b.PageSize = defaultPageSize
	}
}

type BoardColumnsPagedModel struct {
	Items      []BoardColumnModel `json:"items"`
	TotalCount int                `json:"totalCount"`
	TotalPages int                `json:"totalPages"`
	PageNumber int                `json:"pageNumber"`
	PageSize   int                `json:"pageSize"`
}

func (m BoardColumnsPagedModel) Empty(pageNumber, pageSize int) BoardColumnsPagedModel {
	return BoardColumnsPagedModel{
		Items:      []BoardColumnModel{},
		TotalCount: 0,
		TotalPages: 0,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}

type BoardReader interface {
	GetBoard(ctx context.Context, id pgtype.UUID) (BoardModel, error)
	ListBoards(ctx context.Context, q BoardsSearchModel) (BoardsPagedModel, error)
	GetBoardColumn(ctx context.Context, id pgtype.UUID) (BoardColumnModel, error)
	ListBoardColumns(ctx context.Context, q BoardColumnsSearchModel) (BoardColumnsPagedModel, error)
}

type BoardWriter interface {
	CreateBoard(ctx context.Context, b BoardCreateModel) (BoardModel, error)
	UpdateBoard(ctx context.Context, id pgtype.UUID, b BoardUpdateModel) (BoardModel, error)
	ReorderBoards(ctx context.Context, sprintID pgtype.UUID, reorder BoardReorderModel) ([]BoardModel, error)
	DeleteBoard(ctx context.Context, id pgtype.UUID) error
	CreateBoardColumn(ctx context.Context, boardID pgtype.UUID, b BoardColumnCreateModel) (BoardColumnModel, error)
	UpdateBoardColumn(ctx context.Context, boardID, columnID pgtype.UUID, b BoardColumnUpdateModel) (BoardColumnModel, error)
	ReorderBoardColumns(ctx context.Context, boardID pgtype.UUID, reorder BoardColumnReorderModel) ([]BoardColumnModel, error)
	DeleteBoardColumn(ctx context.Context, boardID, columnID pgtype.UUID) error
}
