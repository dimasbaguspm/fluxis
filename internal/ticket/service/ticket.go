package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dimasbaguspm/fluxis/internal/ticket/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/syncx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrTicketNotFound = httpx.NotFound("ticket not found")
)

func (s *Service) ListTickets(ctx context.Context, q domain.TicketSearchModel) ([]domain.TicketModel, error) {
	var tickets []repository.Ticket
	var err error

	if q.SprintID.Valid {
		tickets, err = s.Repo.ListTicketsBySprint(ctx, repository.ListTicketsBySprintParams{
			ProjectID: q.ProjectID,
			SprintID:  q.SprintID,
		})
	} else if q.BoardID.Valid {
		tickets, err = s.Repo.ListTicketsByBoard(ctx, q.BoardID)
	} else {
		tickets, err = s.Repo.ListTicketsByProject(ctx, q.ProjectID)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.TicketModel{}, nil
		}
		return []domain.TicketModel{}, fmt.Errorf("list tickets: %w", err)
	}

	data := make([]domain.TicketModel, 0, len(tickets))
	for _, t := range tickets {
		data = append(data, s.ticketToModel(t))
	}

	return data, nil
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

	// Generate ticket key
	key, err := s.Repo.GenerateTicketKey(ctx, projectID)
	if err != nil {
		return domain.TicketModel{}, fmt.Errorf("generate ticket key: %w", err)
	}

	// Convert DueDate to pgtype.Date if provided
	var dueDate pgtype.Date
	if p.DueDate != nil {
		dueDate = pgtype.Date{
			Time:  *p.DueDate,
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

	return s.ticketToModel(ticket), nil
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
	if p.DueDate != nil {
		dueDate = pgtype.Date{
			Time:  *p.DueDate,
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

	return s.ticketToModel(ticket), nil
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

	return s.ticketToModel(ticket), nil
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

	return s.ticketToModel(ticket), nil
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

	return s.ticketToModel(ticket), nil
}

func (s *Service) DeleteTicket(ctx context.Context, id pgtype.UUID) error {
	_, err := s.Repo.DeleteTicket(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTicketNotFound
		}
		return fmt.Errorf("delete ticket: %w", err)
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
