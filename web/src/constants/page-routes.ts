/**
 * Page routes and deep links constants
 * All routes and their variations are centralized here for type-safe navigation
 */

export const PAGES = {
  SIGN_IN: "/sign-in",
  SIGN_UP: "/sign-up",
  DASHBOARD: "/dashboard",
  ORGANIZATIONS: "/organizations",
  PROJECTS: "/projects",
  BOARDS: "/boards",
  SETTINGS: "/settings",
  PROFILE: "/profile",
} as const;

export const DEEP_LINKS = {
  // Auth
  SIGN_IN: PAGES.SIGN_IN,
  SIGN_UP: PAGES.SIGN_UP,

  // Main app
  DASHBOARD: PAGES.DASHBOARD,

  // Organizations
  ORGANIZATIONS: PAGES.ORGANIZATIONS,
  ORG_DETAILS: (orgId: string) => `/organizations/${orgId}`,

  // Projects
  PROJECTS: PAGES.PROJECTS,
  PROJECT_DETAILS: (projectId: string) => `/projects/${projectId}`,
  PROJECT_BOARD: (projectId: string) => `/projects/${projectId}/boards`,
  PROJECT_BOARD_DETAIL: (projectId: string, boardId: string) =>
    `/projects/${projectId}/boards/${boardId}`,

  // Boards
  BOARDS: PAGES.BOARDS,
  BOARD_DETAILS: (boardId: string) => `/boards/${boardId}`,

  // Settings
  SETTINGS: PAGES.SETTINGS,
  PROFILE: PAGES.PROFILE,
} as const;

/**
 * Route groups for organization
 */
export const ROUTE_GROUPS = {
  UNPROTECTED: "unprotected",
  PROTECTED: "protected",
} as const;
