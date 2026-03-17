import type { DomainUserModel } from "@interfaces/openapi.generated";
import type { UseQueryOptions } from "@tanstack/react-query";
import { useApiQuery } from "../use-api-query";
import { queryKeys } from "../query-keys";

/**
 * Get current authenticated user profile
 */
export function useCurrentUser(
  options?: Omit<UseQueryOptions<DomainUserModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    queryKeys.user.me,
    {
      method: "GET",
      path: "/users/me",
    },
    options,
  );
}
