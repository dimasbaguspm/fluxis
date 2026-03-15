import { DRAWER_ROUTES } from "./drawer-routes";

export interface CreateOrganisationDrawerParams {
  // No params needed for creating
}

export interface UpdateOrganisationDrawerParams {
  orgId: string;
}

/**
 * Mapping of drawer routes to their specific param types
 */
export type DrawerParamsMap = {
  [DRAWER_ROUTES.CREATE_ORGANISATION]: CreateOrganisationDrawerParams | null;
  [DRAWER_ROUTES.UPDATE_ORGANISATION]: UpdateOrganisationDrawerParams;
};

/**
 * Get the param type for a specific drawer route
 */
export type GetDrawerParams<DrawerId extends string> = DrawerId extends keyof DrawerParamsMap
  ? DrawerParamsMap[DrawerId]
  : never;
