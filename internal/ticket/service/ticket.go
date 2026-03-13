package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dimasbaguspm/fluxis/internal/ticket/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
	"github.com/dimasbaguspm/fluxis/pkg/syncx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrTicketNotFound = httpx.NotFound("ticket not found")
)

func (s *Service) ListTickets(ctx context.Context, q domain.TicketSearchModel) (domain.TicketsPagedModel, error) {
	q.ApplyDefaults()

	// Require at least projectId for listing
	if len(q.ProjectID) == 0 {
		return domain.TicketsPagedModel{}, httpx.BadRequest("projectId is required")
	}

	offset := int32((q.PageNumber - 1) * q.PageSize)
	rows, err := s.Repo.ListTicketsPaged(ctx, repository.ListTicketsPagedParams{
		Column1: q.ProjectID,
		Column2: q.ID,
		Column3: q.SprintID,
		Column4: q.BoardID,
		Limit:   int32(q.PageSize),
		Offset:  offset,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TicketsPagedModel{}.Empty(q.PageNumber, q.PageSize), nil
		}
		return domain.TicketsPagedModel{}, fmt.Errorf("list tickets: %w", err)
	}

	if len(rows) == 0 {
		return domain.TicketsPagedModel{}.Empty(q.PageNumber, q.PageSize), nil
	}

	totalCount := int(rows[0].TotalCount)
	items := make([]domain.TicketModel, len(rows))
	for i, row := range rows {
		items[i] = domain.TicketModel{
			ID:            row.ID,
			ProjectID:     row.ProjectID,
			TicketNumber:  row.TicketNumber,
			Key:           row.Key,
			Type:          string(row.Type),
			Priority:      string(row.Priority),
			Title:         row.Title,
			Description:   row.Description.String,
			SprintID:      row.SprintID,
			BoardID:       row.BoardID,
			BoardColumnID: row.BoardColumnID,
			AssigneeID:    row.AssigneeID,
			ReporterID:    row.ReporterID,
			EpicID:        row.EpicID,
			ParentID:      row.ParentID,
			StoryPoints:   row.StoryPoints.Int32,
			DueDate:       row.DueDate.Time,
			CreatedAt:     row.CreatedAt.Time,
			UpdatedAt:     row.UpdatedAt.Time,
		}
	}

	totalPages := (totalCount + q.PageSize - 1) / q.PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return domain.TicketsPagedModel{
		Items:      items,
		TotalCount: totalCount,
		TotalPages: totalPages,
		PageNumber: q.PageNumber,
		PageSize:   q.PageSize,
	}, nil
}

func (s *Service) GetTicket(ctx context.Context, id pgtype.UUID) (domain.TicketModel, error) {
	ticket, err := s.Repo.GetTicket(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TicketModel{}, ErrTicketNotFound
		}
		return domain.TicketModel{}, fmt.Errorf("get ticket: %w", err)
	}

	return s.ticketToModel(ticket), nil
}

func (s *Service) GetTicketByKey(ctx context.Context, projectID pgtype.UUID, key string) (domain.TicketModel, error) {
	ticket, err := s.Repo.GetTicketByKey(ctx, repository.GetTicketByKeyParams{
		ProjectID: projectID,
		Key:       key,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TicketModel{}, ErrTicketNotFound
		}
		return domain.TicketModel{}, fmt.Errorf("get ticket by key: %w", err)
	}

	return s.ticketToModel(ticket), nil
}

func (s *Service) CreateTicket(ctx context.Context, projectID pgtype.UUID, p domain.TicketCreateModel) (domain.TicketModel, error) {
	userID := httpx.MustUserID(ctx)

	// Validate project exists before creating ticket
	_, err := s.Project.GetProjectById(ctx, projectID)
	if err != nil {
		return domain.TicketModel{}, err
	}

	// Generate ticket key
	key, err := s.Repo.GenerateTicketKey(ctx, projectID)
	if err != nil {
		return domain.TicketModel{}, fmt.Errorf("generate ticket key: %w", err)
	}

	// Convert DueDate to pgtype.Date if provided
	var dueDate pgtype.Date
	if !p.DueDate.IsZero() {
		dueDate = pgtype.Date{
			Time:  p.DueDate,
			Valid: true,
		}
	}

	// Convert AssigneeID
	assigneeID := pgtype.UUID{Valid: false}
	if p.AssigneeID.Valid {
		assigneeID = p.AssigneeID
	}

	ticket, err := s.Repo.CreateTicket(ctx, repository.CreateTicketParams{
		ProjectID:   projectID,
		Key:         key,
		Type:        repository.TicketType(p.Type),
		Priority:    repository.TicketPriority(p.Priority),
		Title:       p.Title,
		Description: pgtype.Text{String: p.Description, Valid: p.Description != ""},
		ReporterID:  userID,
		AssigneeID:  assigneeID,
		StoryPoints: pgtype.Int4{Int32: p.StoryPoints, Valid: p.StoryPoints > 0},
		DueDate:     dueDate,
	})
	if err != nil {
		return domain.TicketModel{}, fmt.Errorf("create ticket: %w", err)
	}

	result := s.ticketToModel(ticket)
	if err := s.Bus.Publish(ctx, pubsub.TicketCreated, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.TicketCreated), "error", err)
	}

	return result, nil
}

func (s *Service) UpdateTicket(ctx context.Context, id pgtype.UUID, p domain.TicketUpdateModel) (domain.TicketModel, error) {
	// Fetch current ticket to preserve values for optional fields
	currentTicket, err := s.Repo.GetTicket(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TicketModel{}, ErrTicketNotFound
		}
		return domain.TicketModel{}, fmt.Errorf("get ticket: %w", err)
	}

	// Convert DueDate to pgtype.Date if provided
	var dueDate pgtype.Date
	if !p.DueDate.IsZero() {
		dueDate = pgtype.Date{
			Time:  p.DueDate,
			Valid: true,
		}
	}

	// Convert AssigneeID
	assigneeID := pgtype.UUID{Valid: false}
	if p.AssigneeID.Valid {
		assigneeID = p.AssigneeID
	}

	// Use current values for empty optional enum fields
	ticketType := p.Type
	if ticketType == "" {
		ticketType = string(currentTicket.Type)
	}

	priority := p.Priority
	if priority == "" {
		priority = string(currentTicket.Priority)
	}

	ticket, err := s.Repo.UpdateTicketDetails(ctx, repository.UpdateTicketDetailsParams{
		ID:          id,
		Title:       p.Title,
		Description: pgtype.Text{String: p.Description, Valid: p.Description != ""},
		Type:        repository.TicketType(ticketType),
		Priority:    repository.TicketPriority(priority),
		AssigneeID:  assigneeID,
		StoryPoints: pgtype.Int4{Int32: p.StoryPoints, Valid: p.StoryPoints > 0},
		DueDate:     dueDate,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TicketModel{}, ErrTicketNotFound
		}
		return domain.TicketModel{}, fmt.Errorf("update ticket: %w", err)
	}

	result := s.ticketToModel(ticket)
	if err := s.Bus.Publish(ctx, pubsub.TicketUpdated, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.TicketUpdated), "error", err)
	}

	return result, nil
}

func (s *Service) MoveTicketToBoard(ctx context.Context, id pgtype.UUID, p domain.TicketBoardMoveModel) (domain.TicketModel, error) {
	var board domain.BoardModel
	var boardColumn domain.BoardColumnModel

	err := syncx.Run(ctx,
		func(ctx context.Context) error {
			b, err := s.Board.GetBoard(ctx, p.BoardID)
			if err != nil {
				return fmt.Errorf("validate board: %w", err)
			}
			board = b
			return nil
		},
		func(ctx context.Context) error {
			bc, err := s.Board.GetBoardColumn(ctx, p.BoardColumnID)
			if err != nil {
				return fmt.Errorf("validate board column: %w", err)
			}
			boardColumn = bc
			return nil
		},
	)
	if err != nil {
		return domain.TicketModel{}, err
	}

	if boardColumn.BoardID != board.ID {
		return domain.TicketModel{}, httpx.BadRequest("board column does not belong to the board")
	}

	ticket, err := s.Repo.UpdateTicketBoard(ctx, repository.UpdateTicketBoardParams{
		ID:            id,
		BoardID:       board.ID,
		BoardColumnID: boardColumn.ID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TicketModel{}, ErrTicketNotFound
		}
		return domain.TicketModel{}, fmt.Errorf("move ticket to board: %w", err)
	}

	result := s.ticketToModel(ticket)
	if err := s.Bus.Publish(ctx, pubsub.TicketMovedToBoard, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.TicketMovedToBoard), "error", err)
	}

	return result, nil
}

func (s *Service) MoveTicketToSprint(ctx context.Context, id pgtype.UUID, sprintID pgtype.UUID) (domain.TicketModel, error) {
	// Validate sprint exists
	if _, err := s.Sprint.GetSprint(ctx, sprintID); err != nil {
		return domain.TicketModel{}, fmt.Errorf("validate sprint: %w", err)
	}

	ticket, err := s.Repo.UpdateTicketSprint(ctx, repository.UpdateTicketSprintParams{
		ID:       id,
		SprintID: sprintID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TicketModel{}, ErrTicketNotFound
		}
		return domain.TicketModel{}, fmt.Errorf("move ticket to sprint: %w", err)
	}

	result := s.ticketToModel(ticket)
	if err := s.Bus.Publish(ctx, pubsub.TicketMovedToSprint, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.TicketMovedToSprint), "error", err)
	}

	return result, nil
}

func (s *Service) MoveTicketToBoardColumn(ctx context.Context, id pgtype.UUID, p domain.TicketBoardMoveModel) (domain.TicketModel, error) {
	var board domain.BoardModel
	var boardColumn domain.BoardColumnModel

	err := syncx.Run(ctx,
		func(ctx context.Context) error {
			b, err := s.Board.GetBoard(ctx, p.BoardID)
			if err != nil {
				return fmt.Errorf("validate board: %w", err)
			}
			board = b
			return nil
		},
		func(ctx context.Context) error {
			bc, err := s.Board.GetBoardColumn(ctx, p.BoardColumnID)
			if err != nil {
				return fmt.Errorf("validate board column: %w", err)
			}
			boardColumn = bc
			return nil
		},
	)
	if err != nil {
		return domain.TicketModel{}, err
	}

	if boardColumn.BoardID != board.ID {
		return domain.TicketModel{}, httpx.BadRequest("board column does not belong to the board")
	}

	ticket, err := s.Repo.UpdateTicketBoard(ctx, repository.UpdateTicketBoardParams{
		ID:            id,
		BoardID:       board.ID,
		BoardColumnID: boardColumn.ID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TicketModel{}, ErrTicketNotFound
		}
		return domain.TicketModel{}, fmt.Errorf("move ticket to board column: %w", err)
	}

	result := s.ticketToModel(ticket)
	if err := s.Bus.Publish(ctx, pubsub.TicketMovedToBoardColumn, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.TicketMovedToBoardColumn), "error", err)
	}

	return result, nil
}

func (s *Service) DeleteTicket(ctx context.Context, id pgtype.UUID) error {
	_, err := s.Repo.DeleteTicket(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTicketNotFound
		}
		return fmt.Errorf("delete ticket: %w", err)
	}

	if err := s.Bus.Publish(ctx, pubsub.TicketDeleted, map[string]string{"id": fmt.Sprintf("%v", id)}); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.TicketDeleted), "error", err)
	}

	return nil
}

// Helper function to convert repository model to domain model
func (s *Service) ticketToModel(t repository.Ticket) domain.TicketModel {
	return domain.TicketModel{
		ID:            t.ID,
		ProjectID:     t.ProjectID,
		TicketNumber:  t.TicketNumber,
		Key:           t.Key,
		Type:          string(t.Type),
		Priority:      string(t.Priority),
		Title:         t.Title,
		Description:   t.Description.String,
		SprintID:      t.SprintID,
		BoardID:       t.BoardID,
		BoardColumnID: t.BoardColumnID,
		AssigneeID:    t.AssigneeID,
		ReporterID:    t.ReporterID,
		EpicID:        t.EpicID,
		ParentID:      t.ParentID,
		StoryPoints:   t.StoryPoints.Int32,
		DueDate:       t.DueDate.Time,
		CreatedAt:     t.CreatedAt.Time,
		UpdatedAt:     t.UpdatedAt.Time,
	}
}
