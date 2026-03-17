/**
 * Frontend-friendly event type object, grouped by entity
 */
export const NotifierEventType = {
  Org: "org",
  OrgMember: "orgmember",
  Project: "project",
  Sprint: "sprint",
  Board: "board",
  BoardColumn: "boardcolumn",
  Ticket: "ticket",
  User: "user",
  Auth: "auth",
  Unknown: "unknown",
} as const;

export type NotifierEventTypeValue = (typeof NotifierEventType)[keyof typeof NotifierEventType];

/**
 * Type-safe action constants
 */
export const NotifierAction = {
  Created: "created",
  Updated: "updated",
  Deleted: "deleted",
  Added: "added",
  Removed: "removed",
  Started: "started",
  Completed: "completed",
  Reordered: "reordered",
  VisibilityUpdated: "visibility_updated",
  MovedToBoard: "moved_to_board",
  MovedToBoardColumn: "moved_to_board_column",
  MovedToSprint: "moved_to_sprint",
  Login: "login",
  Logout: "logout",
  Refresh: "refresh",
  Unknown: "unknown",
} as const;

export type NotifierActionValue = (typeof NotifierAction)[keyof typeof NotifierAction];

/**
 * Maps backend event string to type-safe { type, action } pair
 * Avoids direct string extraction from backend event names
 */
export const BACKEND_EVENT_MAP: Record<
  string,
  { type: NotifierEventTypeValue; action: NotifierActionValue }
> = {
  // Auth events
  "auth.auth.login": {
    type: NotifierEventType.Auth,
    action: NotifierAction.Login,
  },
  "auth.auth.logout": {
    type: NotifierEventType.Auth,
    action: NotifierAction.Logout,
  },
  "auth.auth.refresh": {
    type: NotifierEventType.Auth,
    action: NotifierAction.Refresh,
  },

  // User events
  "user.user.created": {
    type: NotifierEventType.User,
    action: NotifierAction.Created,
  },
  "user.user.updated": {
    type: NotifierEventType.User,
    action: NotifierAction.Updated,
  },
  "user.user.deleted": {
    type: NotifierEventType.User,
    action: NotifierAction.Deleted,
  },

  // Org events
  "org.org.created": {
    type: NotifierEventType.Org,
    action: NotifierAction.Created,
  },
  "org.org.updated": {
    type: NotifierEventType.Org,
    action: NotifierAction.Updated,
  },
  "org.org.deleted": {
    type: NotifierEventType.Org,
    action: NotifierAction.Deleted,
  },

  // OrgMember events
  "org.orgmember.added": {
    type: NotifierEventType.OrgMember,
    action: NotifierAction.Added,
  },
  "org.orgmember.updated": {
    type: NotifierEventType.OrgMember,
    action: NotifierAction.Updated,
  },
  "org.orgmember.removed": {
    type: NotifierEventType.OrgMember,
    action: NotifierAction.Removed,
  },

  // Project events
  "project.project.created": {
    type: NotifierEventType.Project,
    action: NotifierAction.Created,
  },
  "project.project.updated": {
    type: NotifierEventType.Project,
    action: NotifierAction.Updated,
  },
  "project.project.deleted": {
    type: NotifierEventType.Project,
    action: NotifierAction.Deleted,
  },
  "project.project.visibility_updated": {
    type: NotifierEventType.Project,
    action: NotifierAction.VisibilityUpdated,
  },

  // Sprint events
  "sprint.sprint.created": {
    type: NotifierEventType.Sprint,
    action: NotifierAction.Created,
  },
  "sprint.sprint.updated": {
    type: NotifierEventType.Sprint,
    action: NotifierAction.Updated,
  },
  "sprint.sprint.started": {
    type: NotifierEventType.Sprint,
    action: NotifierAction.Started,
  },
  "sprint.sprint.completed": {
    type: NotifierEventType.Sprint,
    action: NotifierAction.Completed,
  },

  // Board events
  "board.board.created": {
    type: NotifierEventType.Board,
    action: NotifierAction.Created,
  },
  "board.board.updated": {
    type: NotifierEventType.Board,
    action: NotifierAction.Updated,
  },
  "board.board.deleted": {
    type: NotifierEventType.Board,
    action: NotifierAction.Deleted,
  },
  "board.board.reordered": {
    type: NotifierEventType.Board,
    action: NotifierAction.Reordered,
  },

  // BoardColumn events
  "board.boardcolumn.created": {
    type: NotifierEventType.BoardColumn,
    action: NotifierAction.Created,
  },
  "board.boardcolumn.updated": {
    type: NotifierEventType.BoardColumn,
    action: NotifierAction.Updated,
  },
  "board.boardcolumn.deleted": {
    type: NotifierEventType.BoardColumn,
    action: NotifierAction.Deleted,
  },
  "board.boardcolumn.reordered": {
    type: NotifierEventType.BoardColumn,
    action: NotifierAction.Reordered,
  },

  // Ticket events
  "ticket.ticket.created": {
    type: NotifierEventType.Ticket,
    action: NotifierAction.Created,
  },
  "ticket.ticket.updated": {
    type: NotifierEventType.Ticket,
    action: NotifierAction.Updated,
  },
  "ticket.ticket.deleted": {
    type: NotifierEventType.Ticket,
    action: NotifierAction.Deleted,
  },
  "ticket.ticket.moved_to_board": {
    type: NotifierEventType.Ticket,
    action: NotifierAction.MovedToBoard,
  },
  "ticket.ticket.moved_to_board_column": {
    type: NotifierEventType.Ticket,
    action: NotifierAction.MovedToBoardColumn,
  },
  "ticket.ticket.moved_to_sprint": {
    type: NotifierEventType.Ticket,
    action: NotifierAction.MovedToSprint,
  },
};
