package pubsub

import (
	"context"
	"strings"
)

type EventType string

type Event struct {
	Type    EventType
	Payload map[string]string
}

func Channel(et EventType) string {
	s := string(et)
	if i := strings.Index(s, "."); i >= 0 {
		s = s[:i]
	}
	return "events:" + s
}

type Publisher interface {
	Publish(ctx context.Context, et EventType, payload map[string]string) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, channel string, handler func(context.Context, Event) error)
}

type Bus interface {
	Publisher
	Subscriber
	Close() error
}

// Base channel constants for subscriptions
const (
	Auth    EventType = "auth"
	User    EventType = "user"
	Org     EventType = "org"
	Project EventType = "project"
	Sprint  EventType = "sprint"
	Board   EventType = "board"
	Ticket  EventType = "ticket"
)

// Event variant constants for publishing
const (
	AuthLogin   EventType = "auth.auth.login"
	AuthLogout  EventType = "auth.auth.logout"
	AuthRefresh EventType = "auth.auth.refresh"

	UserCreated EventType = "user.user.created"
	UserUpdated EventType = "user.user.updated"
	UserDeleted EventType = "user.user.deleted"
)

const (
	OrgCreated EventType = "org.org.created"
	OrgUpdated EventType = "org.org.updated"
	OrgDeleted EventType = "org.org.deleted"

	OrgMemberAdded   EventType = "org.orgmember.added"
	OrgMemberUpdated EventType = "org.orgmember.updated"
	OrgMemberRemoved EventType = "org.orgmember.removed"
)

const (
	ProjectCreated           EventType = "project.project.created"
	ProjectUpdated           EventType = "project.project.updated"
	ProjectDeleted           EventType = "project.project.deleted"
	ProjectVisibilityUpdated EventType = "project.project.visibility_updated"
)

const (
	SprintCreated   EventType = "sprint.sprint.created"
	SprintUpdated   EventType = "sprint.sprint.updated"
	SprintStarted   EventType = "sprint.sprint.started"
	SprintCompleted EventType = "sprint.sprint.completed"
)

const (
	BoardCreated   EventType = "board.board.created"
	BoardUpdated   EventType = "board.board.updated"
	BoardDeleted   EventType = "board.board.deleted"
	BoardReordered EventType = "board.board.reordered"

	BoardColumnCreated   EventType = "board.boardcolumn.created"
	BoardColumnUpdated   EventType = "board.boardcolumn.updated"
	BoardColumnDeleted   EventType = "board.boardcolumn.deleted"
	BoardColumnReordered EventType = "board.boardcolumn.reordered"
)

const (
	TicketCreated EventType = "ticket.ticket.created"
	TicketUpdated EventType = "ticket.ticket.updated"
	TicketDeleted EventType = "ticket.ticket.deleted"

	TicketMovedToBoard       EventType = "ticket.ticket.moved_to_board"
	TicketMovedToBoardColumn EventType = "ticket.ticket.moved_to_board_column"
	TicketMovedToSprint      EventType = "ticket.ticket.moved_to_sprint"
)
