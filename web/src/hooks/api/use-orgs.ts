import type {
  DomainOrganisationCreateModel,
  DomainOrganisationMemberCreateModel,
  DomainOrganisationMembersPagedModel,
  DomainOrganisationMemberUpdateModel,
  DomainOrganisationModel,
  DomainOrganisationPagedModel,
  DomainOrganisationUpdateModel,
} from "@interfaces/openapi.generated";
import type { UseQueryOptions } from "@tanstack/react-query";
import { useApiMutation } from "../use-api-mutation";
import { useApiQuery } from "../use-api-query";

interface ListOrgsParams {
  id?: string[];
  name?: string[];
  pageNumber?: number;
  pageSize?: number;
  sortBy?: "name" | "createdAt" | "updatedAt";
  sortOrder?: "asc" | "desc";
}

/**
 * List organisations with pagination and filtering
 */
export function useListOrgs(
  params?: ListOrgsParams,
  options?: Omit<UseQueryOptions<DomainOrganisationPagedModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["orgs", params],
    {
      method: "GET",
      path: "/orgs",
      params,
    },
    options,
  );
}

/**
 * Get a single organisation by ID
 */
export function useGetOrg(
  orgId: string,
  options?: Omit<UseQueryOptions<DomainOrganisationModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["orgs", orgId],
    {
      method: "GET",
      path: `/orgs/${orgId}`,
    },
    { ...options, enabled: !!orgId },
  );
}

/**
 * Create a new organisation
 */
export function useCreateOrg() {
  return useApiMutation<DomainOrganisationModel, unknown, DomainOrganisationCreateModel>(
    (variables) => ({
      method: "POST",
      path: "/orgs",
      body: variables,
    }),
  );
}

/**
 * Update an organisation
 */
export function useUpdateOrg() {
  return useApiMutation<
    DomainOrganisationModel,
    unknown,
    { id: string; data: DomainOrganisationUpdateModel }
  >((variables) => ({
    method: "PATCH",
    path: `/orgs/${variables.id}`,
    body: variables.data,
  }));
}

/**
 * Delete an organisation
 */
export function useDeleteOrg() {
  return useApiMutation<void, unknown, string>((orgId) => ({
    method: "DELETE",
    path: `/orgs/${orgId}`,
  }));
}

interface ListOrgMembersParams {
  displayName?: string;
  email?: string;
  pageNumber?: number;
  pageSize?: number;
  userId?: string[];
}

/**
 * List organisation members
 */
export function useListOrgMembers(
  orgId: string,
  params?: ListOrgMembersParams,
  options?: Omit<UseQueryOptions<DomainOrganisationMembersPagedModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["orgs", orgId, "members", params],
    {
      method: "GET",
      path: `/orgs/${orgId}/members`,
      params,
    },
    { ...options, enabled: !!orgId },
  );
}

/**
 * Add a member to an organisation
 */
export function useAddOrgMember() {
  return useApiMutation<
    void,
    unknown,
    { orgId: string; data: DomainOrganisationMemberCreateModel }
  >((variables) => ({
    method: "POST",
    path: `/orgs/${variables.orgId}/members`,
    body: variables.data,
  }));
}

/**
 * Update an organisation member's role
 */
export function useUpdateOrgMember() {
  return useApiMutation<
    void,
    unknown,
    { orgId: string; userId: string; data: DomainOrganisationMemberUpdateModel }
  >((variables) => ({
    method: "PATCH",
    path: `/orgs/${variables.orgId}/members/${variables.userId}`,
    body: variables.data,
  }));
}

/**
 * Remove a member from an organisation
 */
export function useRemoveOrgMember() {
  return useApiMutation<void, unknown, { orgId: string; userId: string }>((variables) => ({
    method: "DELETE",
    path: `/orgs/${variables.orgId}/members/${variables.userId}`,
  }));
}
