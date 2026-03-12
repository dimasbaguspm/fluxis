package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type TicketSearchModel struct {
	ProjectID pgtype.UUID `json:"projectId" validate:"uuid4"`
	SprintID  pgtype.UUID `json:"sprintId" validate:"omitempty,uuid4"`
	BoardID   pgtype.UUID `json:"boardId" validate:"omitempty,uuid4"`
}

type TicketModel struct {
	ID            pgtype.UUID `json:"id" validate:"required,uuid4"`
	ProjectID     pgtype.UUID `json:"projectId" validate:"required,uuid4"`
	TicketNumber  int32       `json:"ticketNumber"`
	Key           string      `json:"key"`
	Type          string      `json:"type"`
	Priority      string      `json:"priority"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	SprintID      pgtype.UUID `json:"sprintId"`
	BoardID       pgtype.UUID `json:"boardId"`
	BoardColumnID pgtype.UUID `json:"boardColumnId"`
	AssigneeID    pgtype.UUID `json:"assigneeId"`
	ReporterID    pgtype.UUID `json:"reporterId"`
	EpicID        pgtype.UUID `json:"epicId"`
	ParentID      pgtype.UUID `json:"parentId"`
	StoryPoints   int32       `json:"storyPoints"`
	DueDate       time.Time   `json:"dueDate"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}

type TicketCreateModel struct {
	Type        string      `json:"type" validate:"required,oneof=bug story task epic"`
	Priority    string      `json:"priority" validate:"required,oneof=low medium high critical"`
	Title       string      `json:"title" validate:"required,min=1,max=255"`
	Description string      `json:"description"`
	AssigneeID  pgtype.UUID `json:"assigneeId" validate:"omitempty,uuid4"`
	SprintID    pgtype.UUID `json:"sprintId" validate:"omitempty,uuid4"`
	StoryPoints int32       `json:"storyPoints" validate:"omitempty,min=0"`
	DueDate     time.Time   `json:"dueDate,omitempty"`
}

type TicketUpdateModel struct {
	Title       string      `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description string      `json:"description,omitempty"`
	Type        string      `json:"type,omitempty" validate:"omitempty,oneof=bug story task epic"`
	Priority    string      `json:"priority,omitempty" validate:"omitempty,oneof=low medium high critical"`
	AssigneeID  pgtype.UUID `json:"assigneeId,omitempty" validate:"omitempty,uuid4"`
	SprintID    pgtype.UUID `json:"sprintId,omitempty" validate:"omitempty,uuid4"`
	StoryPoints int32       `json:"storyPoints,omitempty" validate:"omitempty,min=0"`
	DueDate     time.Time   `json:"dueDate,omitempty"`
}

type TicketBoardMoveModel struct {
	BoardID       pgtype.UUID `json:"boardId" validate:"required"`
	BoardColumnID pgtype.UUID `json:"boardColumnId" validate:"required"`
}

type TicketReader interface {
	ListTickets(ctx context.Context, q TicketSearchModel) ([]TicketModel, error)
	GetTicket(ctx context.Context, id pgtype.UUID) (TicketModel, error)
	GetTicketByKey(ctx context.Context, projectID pgtype.UUID, key string) (TicketModel, error)
}

type TicketWriter interface {
	CreateTicket(ctx context.Context, projectID pgtype.UUID, p TicketCreateModel) (TicketModel, error)
	UpdateTicket(ctx context.Context, id pgtype.UUID, p TicketUpdateModel) (TicketModel, error)
	MoveTicketToBoard(ctx context.Context, id pgtype.UUID, p TicketBoardMoveModel) (TicketModel, error)
	MoveTicketToSprint(ctx context.Context, id pgtype.UUID, sprintID pgtype.UUID) (TicketModel, error)
	MoveTicketToBoardColumn(ctx context.Context, id pgtype.UUID, p TicketBoardMoveModel) (TicketModel, error)
	DeleteTicket(ctx context.Context, id pgtype.UUID) error
}
