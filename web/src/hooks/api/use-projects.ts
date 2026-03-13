import type {
  DomainProjectCreateModel,
  DomainProjectModel,
  DomainProjectsPagedModel,
  DomainProjectUpdateModel,
} from "@interfaces/openapi.generated";
import type { UseQueryOptions } from "@tanstack/react-query";
import { useApiMutation } from "../use-api-mutation";
import { useApiQuery } from "../use-api-query";

interface ListProjectsParams {
  id?: string[];
  name?: string;
  orgId?: string[];
  pageNumber?: number;
  pageSize?: number;
}

/**
 * List projects with pagination and filtering
 */
export function useListProjects(
  params?: ListProjectsParams,
  options?: Omit<UseQueryOptions<DomainProjectsPagedModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["projects", params],
    {
      method: "GET",
      path: "/projects",
      params,
    },
    options,
  );
}

/**
 * Get a single project by ID
 */
export function useGetProject(
  projectId: string,
  options?: Omit<UseQueryOptions<DomainProjectModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["projects", projectId],
    {
      method: "GET",
      path: `/projects/${projectId}`,
    },
    { ...options, enabled: !!projectId },
  );
}

/**
 * Create a new project
 */
export function useCreateProject() {
  return useApiMutation<
    DomainProjectModel,
    unknown,
    { orgId: string; data: DomainProjectCreateModel }
  >((variables) => ({
    method: "POST",
    path: "/projects",
    body: { ...variables.data, orgId: variables.orgId },
  }));
}

/**
 * Update a project
 */
export function useUpdateProject() {
  return useApiMutation<
    DomainProjectModel,
    unknown,
    { projectId: string; data: DomainProjectUpdateModel }
  >((variables) => ({
    method: "PATCH",
    path: `/projects/${variables.projectId}`,
    body: variables.data,
  }));
}

/**
 * Update project visibility
 */
export function useUpdateProjectVisibility() {
  return useApiMutation<
    DomainProjectModel,
    unknown,
    { projectId: string; visibility: "public" | "private" }
  >((variables) => ({
    method: "PATCH",
    path: `/projects/${variables.projectId}/visibility`,
    body: { visibility: variables.visibility },
  }));
}

/**
 * Delete a project
 */
export function useDeleteProject() {
  return useApiMutation<void, unknown, string>((projectId) => ({
    method: "DELETE",
    path: `/projects/${projectId}`,
  }));
}
