import type { UseMutationOptions } from "@tanstack/react-query";
import type {
  DomainAuthLoginModel,
  DomainAuthModel,
  DomainAuthRefreshModel,
  DomainAuthRegisterModel,
  HttpxErrBlock,
} from "@interfaces/openapi.generated";
import { useApiMutation } from "../use-api-mutation";

/**
 * Login with email and password
 */
export function useLogin(
  options?: Omit<
    UseMutationOptions<DomainAuthModel, HttpxErrBlock, DomainAuthLoginModel>,
    "mutationFn"
  >,
) {
  return useApiMutation<DomainAuthModel, HttpxErrBlock, DomainAuthLoginModel>(
    (variables) => ({
      method: "POST",
      path: "/auth/login",
      body: variables,
    }),
    options,
  );
}

/**
 * Register a new user account
 */
export function useRegister(
  options?: Omit<
    UseMutationOptions<DomainAuthModel, HttpxErrBlock, DomainAuthRegisterModel>,
    "mutationFn"
  >,
) {
  return useApiMutation<DomainAuthModel, HttpxErrBlock, DomainAuthRegisterModel>(
    (variables) => ({
      method: "POST",
      path: "/auth/register",
      body: variables,
    }),
    options,
  );
}

/**
 * Refresh access token using refresh token
 */
export function useRefreshToken(
  options?: Omit<
    UseMutationOptions<DomainAuthModel, HttpxErrBlock, DomainAuthRefreshModel>,
    "mutationFn"
  >,
) {
  return useApiMutation<DomainAuthModel, HttpxErrBlock, DomainAuthRefreshModel>(
    (variables) => ({
      method: "POST",
      path: "/auth/refresh",
      body: variables,
    }),
    options,
  );
}
