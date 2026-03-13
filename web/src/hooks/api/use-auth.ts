import type {
  DomainAuthLoginModel,
  DomainAuthModel,
  DomainAuthRefreshModel,
  DomainAuthRegisterModel,
} from "@interfaces/openapi.generated";
import { useApiMutation } from "../use-api-mutation";

/**
 * Login with email and password
 */
export function useLogin() {
  return useApiMutation<DomainAuthModel, unknown, DomainAuthLoginModel>(
    (variables) => ({
      method: "POST",
      path: "/auth/login",
      body: variables,
    }),
  );
}

/**
 * Register a new user account
 */
export function useRegister() {
  return useApiMutation<DomainAuthModel, unknown, DomainAuthRegisterModel>(
    (variables) => ({
      method: "POST",
      path: "/auth/register",
      body: variables,
    }),
  );
}

/**
 * Refresh access token using refresh token
 */
export function useRefreshToken() {
  return useApiMutation<DomainAuthModel, unknown, DomainAuthRefreshModel>(
    (variables) => ({
      method: "POST",
      path: "/auth/refresh",
      body: variables,
    }),
  );
}
