import { useSessionHandler, useSessionState } from "@/providers/session";
import type { DomainAuthModel } from "@interfaces/openapi.generated";
import { useCallback } from "react";
import { request } from "./http-client";
import type { HttpRequest } from "./http-request";
import { HttpError } from "./http-request";

export type Fetcher = <TData>(
  requestConfig: HttpRequest,
  extraHeaders?: Record<string, string>,
) => Promise<TData>;

/**
 * Hook that returns a fetcher function with silent token refresh logic.
 *
 * If a request returns 401 and a refresh token exists:
 * 1. Attempts to refresh the access token via POST /auth/refresh
 * 2. On success: stores new tokens and retries the original request
 * 3. On failure: logs out the user and re-throws the original error
 *
 * Consumers see a single pending state even if up to 3 HTTP calls occur internally.
 */
export function useFetcher(): Fetcher {
  const { accessToken, refreshToken } = useSessionState();
  const { setTokens, logout } = useSessionHandler();

  const fetch = useCallback(
    async <TData>(
      requestConfig: HttpRequest,
      extraHeaders?: Record<string, string>,
    ): Promise<TData> => {
      const headers: Record<string, string> = { ...extraHeaders };
      if (accessToken) {
        headers.Authorization = `Bearer ${accessToken}`;
      }

      try {
        const response = await request<TData>(requestConfig, { headers });
        return response.data;
      } catch (err) {
        // Check if this is a 401 and we have a refresh token
        if (err instanceof HttpError && err.status === 401 && refreshToken) {
          try {
            // Attempt silent refresh
            const refreshResponse = await request<DomainAuthModel>({
              method: "POST",
              path: "/auth/refresh",
              body: { accessToken, refreshToken },
            });

            const newAccess = refreshResponse.data.accessToken ?? null;
            const newRefresh = refreshResponse.data.refreshToken ?? null;
            setTokens(newAccess, newRefresh);

            // Retry the original request with the new token
            const retryHeaders: Record<string, string> = { ...extraHeaders };
            if (newAccess) {
              retryHeaders.Authorization = `Bearer ${newAccess}`;
            }
            const retryResponse = await request<TData>(requestConfig, {
              headers: retryHeaders,
            });
            return retryResponse.data;
          } catch {
            // Refresh failed, log out and re-throw the original 401 error
            logout();
            throw err;
          }
        }

        // Not a 401 or no refresh token, re-throw as-is
        throw err;
      }
    },
    [accessToken, refreshToken, setTokens, logout],
  );

  return fetch;
}
