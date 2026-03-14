import type { QueryKey, UseInfiniteQueryOptions } from "@tanstack/react-query";
import { useInfiniteQuery } from "@tanstack/react-query";
import type { HttpRequest } from "../use-fetcher";
import { useFetcher } from "../use-fetcher";

interface InfiniteQueryStatus {
  fetchStatus: "idle" | "fetching" | "paused";
  status: "pending" | "error" | "success";
  isLoading: boolean;
  isFetching: boolean;
  isFetchingNextPage: boolean;
  isFetchingPreviousPage: boolean;
  isError: boolean;
  isSuccess: boolean;
  isPending: boolean;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
}

interface InfiniteQueryMethods {
  fetchNextPage: () => void;
  fetchPreviousPage: () => void;
  refetch: () => void;
}

/**
 * Base hook for paginated API queries
 * Automatically injects access token from session
 * Returns tuple format: [pages, error, status, methods]
 *
 * @example
 * const [pages, err, { hasNextPage, isFetchingNextPage }, { fetchNextPage }] = useApiInfiniteQuery(
 *   ['orgs'],
 *   (pageParam) => ({ method: 'GET', path: '/orgs', params: { pageNumber: pageParam } }),
 *   {
 *     initialPageParam: 1,
 *     getNextPageParam: (lastPage, _, lastPageParam) =>
 *       lastPage.totalPages > lastPageParam ? lastPageParam + 1 : undefined,
 *   }
 * )
 */
export function useApiInfiniteQuery<TData = unknown, TError = unknown>(
  queryKey: QueryKey,
  requestConfigFactory: (pageParam: number) => HttpRequest,
  options: Omit<
    UseInfiniteQueryOptions<TData, TError, TData, QueryKey, number>,
    "queryKey" | "queryFn"
  > & { headers?: Record<string, string> },
): [TData[] | undefined, TError | null, InfiniteQueryStatus, InfiniteQueryMethods] {
  const { headers, ...infiniteOptions } = options || {};
  const fetcher = useFetcher();

  const query = useInfiniteQuery({
    queryKey,
    queryFn: async ({ pageParam }: { pageParam: number }) => {
      const requestConfig = requestConfigFactory(pageParam as number);
      return fetcher<TData>(requestConfig, headers);
    },
    ...infiniteOptions,
  });

  const status: InfiniteQueryStatus = {
    fetchStatus: query.fetchStatus,
    status: query.status,
    isLoading: query.isLoading,
    isFetching: query.isFetching,
    isFetchingNextPage: query.isFetchingNextPage,
    isFetchingPreviousPage: query.isFetchingPreviousPage,
    isError: query.isError,
    isSuccess: query.isSuccess,
    isPending: query.isPending,
    hasNextPage: query.hasNextPage ?? false,
    hasPreviousPage: query.hasPreviousPage ?? false,
  };

  const methods: InfiniteQueryMethods = {
    fetchNextPage: () => query.fetchNextPage(),
    fetchPreviousPage: () => query.fetchPreviousPage(),
    refetch: () => query.refetch(),
  };

  return [
    (query.data as any)?.pages as TData[] | undefined,
    query.error as TError | null,
    status,
    methods,
  ];
}
