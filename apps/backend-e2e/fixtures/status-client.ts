import { APIRequestContext } from "@playwright/test";
import { BaseAPIClient } from "./base-client";
import type { TestContext, APIResponse } from "../types/common";
import type { components } from "../types/openapi";

/**
 * Status types from OpenAPI operations
 */
export type StatusCreateRequest = components["schemas"]["StatusCreateModel"];
export type StatusUpdateRequest = components["schemas"]["StatusUpdateModel"];
export type StatusReorderRequest = components["schemas"]["StatusReorderModel"];
export type StatusResponse = components["schemas"]["StatusModel"];
export type LogPaginatedResponse = components["schemas"]["LogPaginatedModel"];

/**
 * Status API client for status operations
 */
export class StatusAPIClient extends BaseAPIClient {
  constructor(request: APIRequestContext, context: TestContext) {
    super(request, context);
  }

  /**
   * Create a new status
   */
  async create(
    data: StatusCreateRequest
  ): Promise<APIResponse<StatusResponse>> {
    return this.post<StatusResponse>("/statuses", data);
  }

  /**
   * Get statuses by project ID (returns array, not paginated)
   */
  async getByProject(
    projectId: string
  ): Promise<APIResponse<StatusResponse[]>> {
    return this.get<StatusResponse[]>("/statuses", { projectId });
  }

  /**
   * Get a single status by ID
   */
  async getById(statusId: string): Promise<APIResponse<StatusResponse>> {
    return this.get<StatusResponse>(`/statuses/${statusId}`);
  }

  /**
   * Update a status
   */
  async update(
    statusId: string,
    data: StatusUpdateRequest
  ): Promise<APIResponse<StatusResponse>> {
    return this.patch<StatusResponse>(`/statuses/${statusId}`, data);
  }

  /**
   * Delete a status
   */
  async remove(statusId: string): Promise<APIResponse<void>> {
    return this.delete<void>(`/statuses/${statusId}`);
  }

  /**
   * Reorder statuses for a project
   */
  async reorder(
    data: StatusReorderRequest
  ): Promise<APIResponse<StatusResponse[]>> {
    return this.post<StatusResponse[]>("/statuses/reorder", data);
  }

  /**
   * Get logs for a status
   */
  async getLogs(
    statusId: string,
    params?: {
      taskId?: string[];
      statusId?: string[];
      query?: string;
      pageNumber?: number;
      pageSize?: number;
    }
  ): Promise<APIResponse<LogPaginatedResponse>> {
    return this.get<LogPaginatedResponse>(`/statuses/${statusId}/logs`, params);
  }
}
