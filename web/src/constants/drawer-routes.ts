export const DRAWER_ROUTES = {
  CREATE_ORGANISATION: "create-org",
  UPDATE_ORGANISATION: "update-org",
} as const;

export type DrawerRouteId = (typeof DRAWER_ROUTES)[keyof typeof DRAWER_ROUTES];
