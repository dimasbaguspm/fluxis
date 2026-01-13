import { APIRequestContext } from "@playwright/test";
import { BaseAPIClient } from "./base-client";
import type { TestContext, APIResponse } from "../types/common";
import type { components } from "../types/openapi";

/**
 * Task types from OpenAPI operations
 */
export type TaskCreateRequest = components["schemas"]["TaskCreateModel"];
export type TaskUpdateRequest = components["schemas"]["TaskUpdateModel"];
export type TaskResponse = components["schemas"]["TaskModel"];
export type TaskPaginatedResponse = components["schemas"]["TaskPaginatedModel"];
export type LogPaginatedResponse = components["schemas"]["LogPaginatedModel"];

/**
 * Task API client for task operations
 */
export class TaskAPIClient extends BaseAPIClient {
  constructor(request: APIRequestContext, context: TestContext) {
    super(request, context);
  }

  /**
   * Create a new task
   */
  async create(data: TaskCreateRequest): Promise<APIResponse<TaskResponse>> {
    return this.post<TaskResponse>("/tasks", data);
  }

  /**
   * Get paginated tasks with filters
   */
  async getPaginated(params?: {
    id?: string[];
    projectId?: string[];
    statusId?: string[];
    query?: string;
    pageNumber?: number;
    pageSize?: number;
    sortBy?: "dueDate" | "createdAt" | "updatedAt" | "priority";
    sortOrder?: "asc" | "desc";
  }): Promise<APIResponse<TaskPaginatedResponse>> {
    return this.get<TaskPaginatedResponse>("/tasks", params);
  }

  /**
   * Get a single task by ID
   */
  async getById(taskId: string): Promise<APIResponse<TaskResponse>> {
    return this.get<TaskResponse>(`/tasks/${taskId}`);
  }

  /**
   * Update a task
   */
  async update(
    taskId: string,
    data: TaskUpdateRequest
  ): Promise<APIResponse<TaskResponse>> {
    return this.patch<TaskResponse>(`/tasks/${taskId}`, data);
  }

  /**
   * Delete a task
   */
  async remove(taskId: string): Promise<APIResponse<void>> {
    return this.delete<void>(`/tasks/${taskId}`);
  }

  /**
   * Get logs for a task
   */
  async getLogs(
    taskId: string,
    params?: {
      taskId?: string[];
      statusId?: string[];
      query?: string;
      pageNumber?: number;
      pageSize?: number;
    }
  ): Promise<APIResponse<LogPaginatedResponse>> {
    return this.get<LogPaginatedResponse>(`/tasks/${taskId}/logs`, params);
  }
}
