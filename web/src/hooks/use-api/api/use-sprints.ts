import type {
  DomainSprintCreateModel,
  DomainSprintModel,
  DomainSprintsPagedModel,
  DomainSprintUpdateModel,
} from "@interfaces/openapi.generated";
import type { UseQueryOptions } from "@tanstack/react-query";
import { useApiMutation } from "../use-api-mutation";
import { useApiQuery } from "../use-api-query";

interface ListSprintsParams {
  id?: string[];
  name?: string;
  pageNumber?: number;
  pageSize?: number;
  projectId?: string[];
}

/**
 * List sprints with pagination and filtering
 */
export function useListSprints(
  params?: ListSprintsParams,
  options?: Omit<UseQueryOptions<DomainSprintsPagedModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["sprints", params],
    {
      method: "GET",
      path: "/sprints",
      params,
    },
    options,
  );
}

/**
 * Get a single sprint by ID
 */
export function useGetSprint(
  sprintId: string,
  options?: Omit<UseQueryOptions<DomainSprintModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["sprints", sprintId],
    {
      method: "GET",
      path: `/sprints/${sprintId}`,
    },
    { enabled: !!sprintId, ...options },
  );
}

/**
 * Create a new sprint
 */
export function useCreateSprint() {
  return useApiMutation<DomainSprintModel, unknown, DomainSprintCreateModel>((variables) => ({
    method: "POST",
    path: "/sprints",
    body: variables,
  }));
}

/**
 * Update a sprint
 */
export function useUpdateSprint() {
  return useApiMutation<
    DomainSprintModel,
    unknown,
    { sprintId: string; data: DomainSprintUpdateModel }
  >((variables) => ({
    method: "PATCH",
    path: `/sprints/${variables.sprintId}`,
    body: variables.data,
  }));
}

/**
 * Start a sprint
 */
export function useStartSprint() {
  return useApiMutation<DomainSprintModel, unknown, string>((sprintId) => ({
    method: "POST",
    path: `/sprints/${sprintId}/start`,
  }));
}

/**
 * Complete a sprint
 */
export function useCompleteSprint() {
  return useApiMutation<DomainSprintModel, unknown, string>((sprintId) => ({
    method: "POST",
    path: `/sprints/${sprintId}/complete`,
  }));
}
