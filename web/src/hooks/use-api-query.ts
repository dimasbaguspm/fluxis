import { useQuery } from "@tanstack/react-query";
import type { QueryKey, UseQueryOptions } from "@tanstack/react-query";
import { request, type RequestOptions } from "@/lib/http-client";
import type { HttpRequest } from "@/lib/http-request";

interface QueryStatus {
  fetchStatus: "idle" | "fetching" | "paused";
  status: "pending" | "error" | "success";
  isLoading: boolean;
  isFetching: boolean;
  isError: boolean;
  isSuccess: boolean;
  isPending: boolean;
}

interface QueryMethods {
  refetch: () => void;
}

/**
 * Base hook for API queries (GET requests)
 * Returns tuple format: [data, error, status, methods]
 *
 * @example
 * const [org, err, { isLoading }, { refetch }] = useApiQuery(
 *   ['orgs', orgId],
 *   { method: 'GET', path: `/orgs/${orgId}` },
 *   { enabled: !!orgId }
 * )
 */
export function useApiQuery<TData = unknown, TError = unknown>(
  queryKey: QueryKey,
  requestConfig: HttpRequest,
  options?: Omit<UseQueryOptions<TData, TError>, "queryKey" | "queryFn"> & { headers?: Record<string, string> },
): [TData | undefined, TError | null, QueryStatus, QueryMethods] {
  const { headers, ...queryOptions } = options || {};

  const query = useQuery({
    queryKey,
    queryFn: async () => {
      const response = await request<TData>(requestConfig, { headers } as RequestOptions);
      return response.data;
    },
    ...queryOptions,
  });

  const status: QueryStatus = {
    fetchStatus: query.fetchStatus,
    status: query.status,
    isLoading: query.isLoading,
    isFetching: query.isFetching,
    isError: query.isError,
    isSuccess: query.isSuccess,
    isPending: query.isPending,
  };

  const methods: QueryMethods = {
    refetch: () => query.refetch(),
  };

  return [query.data, query.error as TError | null, status, methods];
}
