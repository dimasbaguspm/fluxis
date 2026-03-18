import { DRAWER_ROUTES } from "./drawer-routes";

export interface CreateOrganisationDrawerParams {
  // No params needed for creating
}

export interface UpdateOrganisationDrawerParams {
  orgId: string;
}

export interface CreateProjectDrawerParams {
  orgId?: string;
}

export interface UpdateProjectDrawerParams {
  projectId: string;
}

export interface CreateSprintDrawerParams {
  projectId?: string;
}

export interface UpdateSprintDrawerParams {
  sprintId: string;
}

export interface CreateBoardDrawerParams {
  sprintId?: string;
}

export interface UpdateBoardDrawerParams {
  boardId: string;
}

export interface CreateTicketDrawerParams {
  projectId?: string;
  sprintId?: string;
}

export interface UpdateTicketDrawerParams {
  ticketId: string;
}

/**
 * Mapping of drawer routes to their specific param types
 */
export type DrawerParamsMap = {
  [DRAWER_ROUTES.CREATE_ORGANISATION]: CreateOrganisationDrawerParams | null;
  [DRAWER_ROUTES.UPDATE_ORGANISATION]: UpdateOrganisationDrawerParams;
  [DRAWER_ROUTES.CREATE_PROJECT]: CreateProjectDrawerParams | null;
  [DRAWER_ROUTES.UPDATE_PROJECT]: UpdateProjectDrawerParams;
  [DRAWER_ROUTES.CREATE_SPRINT]: CreateSprintDrawerParams | null;
  [DRAWER_ROUTES.UPDATE_SPRINT]: UpdateSprintDrawerParams;
  [DRAWER_ROUTES.CREATE_BOARD]: CreateBoardDrawerParams | null;
  [DRAWER_ROUTES.UPDATE_BOARD]: UpdateBoardDrawerParams;
  [DRAWER_ROUTES.CREATE_TICKET]: CreateTicketDrawerParams | null;
  [DRAWER_ROUTES.UPDATE_TICKET]: UpdateTicketDrawerParams;
};

/**
 * Get the param type for a specific drawer route
 */
export type GetDrawerParams<DrawerId extends string> = DrawerId extends keyof DrawerParamsMap
  ? DrawerParamsMap[DrawerId]
  : never;
