import type { UseMutationOptions } from "@tanstack/react-query";
import { useMutation } from "@tanstack/react-query";
import type { HttpRequest } from "../use-fetcher";
import { useFetcher } from "../use-fetcher";

interface MutationStatus {
  status: "idle" | "pending" | "success" | "error";
  isPending: boolean;
  isSuccess: boolean;
  isError: boolean;
  isIdle: boolean;
}

type MutationMethod<TVariables> = (variables: TVariables) => Promise<any>;

/**
 * Base hook for API mutations (POST, PUT, PATCH, DELETE requests)
 * Automatically injects access token from session
 * Returns tuple format: [data, error, status, methods]
 *
 * @example
 * const [result, err, { isPending }, { mutate }] = useApiMutation(
 *   (variables: LoginPayload) => ({
 *     method: 'POST',
 *     path: '/auth/login',
 *     body: variables,
 *   }),
 *   { onSuccess: (data) => setSession(data) }
 * )
 *
 * mutate({ email: 'user@example.com', password: 'pass' })
 */
export function useApiMutation<TData = unknown, TError = unknown, TVariables = void>(
  requestConfigFactory: (variables: TVariables) => HttpRequest,
  options?: Omit<UseMutationOptions<TData, TError, TVariables>, "mutationFn"> & {
    headers?: Record<string, string>;
  },
): [MutationMethod<TVariables>, TError | null, MutationStatus] {
  const { headers, ...mutationOptions } = options || {};
  const fetcher = useFetcher();

  const mutation = useMutation({
    mutationFn: async (variables: TVariables) => {
      const requestConfig = requestConfigFactory(variables);
      return fetcher<TData>(requestConfig, headers);
    },
    ...mutationOptions,
  });

  const status: MutationStatus = {
    status: mutation.status,
    isPending: mutation.isPending,
    isSuccess: mutation.isSuccess,
    isError: mutation.isError,
    isIdle: mutation.isIdle,
  };

  return [mutation.mutateAsync, mutation.error as TError | null, status];
}
