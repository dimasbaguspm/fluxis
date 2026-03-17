export const DRAWER_ROUTES = {
  CREATE_ORGANISATION: "create-org",
  UPDATE_ORGANISATION: "update-org",
  CREATE_PROJECT: "create-project",
  UPDATE_PROJECT: "update-project",
  CREATE_SPRINT: "create-sprint",
  UPDATE_SPRINT: "update-sprint",
  CREATE_BOARD: "create-board",
  UPDATE_BOARD: "update-board",
} as const;

export type DrawerRouteId = (typeof DRAWER_ROUTES)[keyof typeof DRAWER_ROUTES];
