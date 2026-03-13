import type { DomainUserModel } from "@interfaces/openapi.generated";
import type { UseQueryOptions } from "@tanstack/react-query";
import { useApiQuery } from "../use-api-query";

/**
 * Get current authenticated user profile
 */
export function useCurrentUser(
  options?: Omit<UseQueryOptions<DomainUserModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["user", "me"],
    {
      method: "GET",
      path: "/users/me",
    },
    options,
  );
}
