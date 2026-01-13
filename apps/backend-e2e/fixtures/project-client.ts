import { APIRequestContext } from "@playwright/test";
import { BaseAPIClient } from "./base-client";
import type { TestContext, APIResponse } from "../types/common";
import type { components } from "../types/openapi";

/**
 * Project types from OpenAPI operations
 */
export type ProjectCreateRequest = components["schemas"]["ProjectCreateModel"];
export type ProjectUpdateRequest = components["schemas"]["ProjectUpdateModel"];
export type ProjectResponse = components["schemas"]["ProjectModel"];
export type ProjectPaginatedResponse =
  components["schemas"]["ProjectPaginatedModel"];
export type LogPaginatedResponse = components["schemas"]["LogPaginatedModel"];

/**
 * Project API client for project operations
 */
export class ProjectAPIClient extends BaseAPIClient {
  constructor(request: APIRequestContext, context: TestContext) {
    super(request, context);
  }

  /**
   * Create a new project
   */
  async create(
    data: ProjectCreateRequest
  ): Promise<APIResponse<ProjectResponse>> {
    return this.post<ProjectResponse>("/projects", data);
  }

  /**
   * Get paginated list of projects
   */
  async getPaginated(params?: {
    id?: string[];
    query?: string;
    status?: ("active" | "paused" | "archived")[];
    pageNumber?: number;
    pageSize?: number;
    sortBy?: "createdAt" | "updatedAt" | "status";
    sortOrder?: "asc" | "desc";
  }): Promise<APIResponse<ProjectPaginatedResponse>> {
    return this.get<ProjectPaginatedResponse>("/projects", params);
  }

  /**
   * Get a single project by ID
   */
  async getById(projectId: string): Promise<APIResponse<ProjectResponse>> {
    return this.get<ProjectResponse>(`/projects/${projectId}`);
  }

  /**
   * Update a project
   */
  async update(
    projectId: string,
    data: ProjectUpdateRequest
  ): Promise<APIResponse<ProjectResponse>> {
    return this.patch<ProjectResponse>(`/projects/${projectId}`, data);
  }

  /**
   * Delete a project
   */
  async remove(projectId: string): Promise<APIResponse<void>> {
    return this.delete<void>(`/projects/${projectId}`);
  }

  /**
   * Get logs for a project
   */
  async getLogs(
    projectId: string,
    params?: {
      taskId?: string[];
      statusId?: string[];
      query?: string;
      pageNumber?: number;
      pageSize?: number;
    }
  ): Promise<APIResponse<LogPaginatedResponse>> {
    return this.get<LogPaginatedResponse>(
      `/projects/${projectId}/logs`,
      params
    );
  }
}
