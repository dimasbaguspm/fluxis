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
	Name     *string     `json:"name" validate:"required,min=1"`
	SprintID pgtype.UUID `json:"sprintId" validate:"required"`
}

type BoardUpdateModel struct {
	Name     *string      `json:"name" validate:"omitempty,min=1"`
	SprintID *pgtype.UUID `json:"sprintId,omitempty"`
}

type BoardReorderModel struct {
	Boards []BoardPositionUpdate `json:"boards" validate:"required"`
}

type BoardPositionUpdate struct {
	ID       pgtype.UUID `json:"id" validate:"required"`
	Position int32       `json:"position" validate:"required,min=0"`
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
	Name     *string `json:"name" validate:"required,min=1"`
	Position *int32  `json:"position" validate:"required,min=0"`
}

type BoardColumnUpdateModel struct {
	Name     *string `json:"name" validate:"omitempty,min=1"`
	Position *int32  `json:"position" validate:"omitempty,min=0"`
}

type BoardReader interface {
	GetBoard(ctx context.Context, id pgtype.UUID) (BoardModel, error)
	ListBoardsBySprint(ctx context.Context, sprintID pgtype.UUID) ([]BoardModel, error)
	GetBoardColumn(ctx context.Context, id pgtype.UUID) (BoardColumnModel, error)
}

type BoardWriter interface {
	CreateBoard(ctx context.Context, sprintID pgtype.UUID, b BoardCreateModel) (BoardModel, error)
	UpdateBoard(ctx context.Context, id pgtype.UUID, b BoardUpdateModel) (BoardModel, error)
	ReorderBoard(ctx context.Context, id pgtype.UUID, position int32) (BoardModel, error)
	DeleteBoard(ctx context.Context, id pgtype.UUID) error
}
