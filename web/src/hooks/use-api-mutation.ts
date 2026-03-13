import type { UseMutationOptions } from "@tanstack/react-query";
import { useMutation } from "@tanstack/react-query";
import { request, type RequestOptions } from "@/lib/http-client";
import type { HttpRequest } from "@/lib/http-request";
import { useSessionStore } from "@providers/session";

interface MutationStatus {
  status: "idle" | "pending" | "success" | "error";
  isPending: boolean;
  isSuccess: boolean;
  isError: boolean;
  isIdle: boolean;
}

interface MutationMethods<TVariables> {
  mutate: (variables: TVariables) => void;
  mutateAsync: (variables: TVariables) => Promise<any>;
}

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
  options?: Omit<UseMutationOptions<TData, TError, TVariables>, "mutationFn"> & { headers?: Record<string, string> },
): [TData | undefined, TError | null, MutationStatus, MutationMethods<TVariables>] {
  const { headers, ...mutationOptions } = options || {};
  const accessToken = useSessionStore((state) => state.accessToken);

  const mutation = useMutation({
    mutationFn: async (variables: TVariables) => {
      const requestConfig = requestConfigFactory(variables);
      const requestHeaders: Record<string, string> = { ...headers };
      if (accessToken) {
        requestHeaders.Authorization = `Bearer ${accessToken}`;
      }
      const response = await request<TData>(requestConfig, { headers: requestHeaders } as RequestOptions);
      return response.data;
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

  const methods: MutationMethods<TVariables> = {
    mutate: mutation.mutate,
    mutateAsync: mutation.mutateAsync,
  };

  return [mutation.data, mutation.error as TError | null, status, methods];
}
