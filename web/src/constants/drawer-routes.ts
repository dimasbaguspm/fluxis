export const DRAWER_ROUTES = {
  CREATE_ORGANISATION: "create-org",
} as const;

export type DrawerRouteId = (typeof DRAWER_ROUTES)[keyof typeof DRAWER_ROUTES];
