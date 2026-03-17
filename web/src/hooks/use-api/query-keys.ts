/**
 * Centralized query keys for all API hooks
 * Used to enable consistent cache invalidation across the application
 */

export const queryKeys = {
  // Auth
  auth: {
    all: ["auth"] as const,
    login: ["auth", "login"] as const,
    register: ["auth", "register"] as const,
  },

  // User
  user: {
    all: ["user"] as const,
    me: ["user", "me"] as const,
  },

  // Organisations
  orgs: {
    all: ["orgs"] as const,
    list: (params?: unknown) => ["orgs", params] as const,
    detail: (orgId: string) => ["orgs", orgId] as const,
    members: {
      all: (orgId: string) => ["orgs", orgId, "members"] as const,
      list: (orgId: string, params?: unknown) => ["orgs", orgId, "members", params] as const,
    },
  },

  // Projects
  projects: {
    all: ["projects"] as const,
    list: (params?: unknown) => ["projects", params] as const,
    detail: (projectId: string) => ["projects", projectId] as const,
  },

  // Sprints
  sprints: {
    all: ["sprints"] as const,
    list: (params?: unknown) => ["sprints", params] as const,
    detail: (sprintId: string) => ["sprints", sprintId] as const,
  },

  // Boards
  boards: {
    all: ["boards"] as const,
    list: (params?: unknown) => ["boards", params] as const,
    detail: (boardId: string) => ["boards", boardId] as const,
    columns: {
      all: (boardId: string) => ["boards", boardId, "columns"] as const,
      list: (boardId: string, params?: unknown) => ["boards", boardId, "columns", params] as const,
    },
  },

  // Tickets
  tickets: {
    all: ["tickets"] as const,
    list: (params?: unknown) => ["tickets", params] as const,
    detail: (ticketId: string) => ["tickets", ticketId] as const,
  },
} as const;
