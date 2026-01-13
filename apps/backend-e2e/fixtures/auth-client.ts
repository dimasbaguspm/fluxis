import { APIRequestContext } from "@playwright/test";
import { BaseAPIClient } from "./base-client";
import type { TestContext, APIResponse } from "../types/common";
import type { components } from "../types/openapi";

/**
 * Auth types from OpenAPI operations
 */
export type LoginRequestModel = components["schemas"]["AuthLoginInputModel"];
export type LoginResponseModel = components["schemas"]["AuthLoginOutputModel"];
export type RefreshRequestModel =
  components["schemas"]["AuthRefreshInputModel"];
export type RefreshResponseModel =
  components["schemas"]["AuthRefreshOutputModel"];

/**
 * Auth API client for authentication operations
 */
export class AuthAPIClient extends BaseAPIClient {
  constructor(request: APIRequestContext, context: TestContext) {
    super(request, context);
  }

  /**
   * Login with username and password
   */
  async login(
    username: string,
    password: string
  ): Promise<APIResponse<LoginResponseModel>> {
    const response = await this.post<LoginResponseModel>("/auth/login", {
      username,
      password,
    });

    if (response.data) {
      this.context.accessToken = response.data.accessToken;
      this.context.refreshToken = response.data.refreshToken;
    }

    return response;
  }

  /**
   * Refresh access token
   */
  async refresh(
    refreshToken?: string
  ): Promise<APIResponse<RefreshResponseModel>> {
    const token = refreshToken || this.context.refreshToken;
    if (!token) {
      throw new Error("No refresh token available");
    }

    const response = await this.post<RefreshResponseModel>("/auth/refresh", {
      refreshToken: token,
    });

    if (response.data) {
      this.context.accessToken = response.data.accessToken;
    }

    return response;
  }

  /**
   * Clear stored tokens
   */
  logout(): void {
    this.context.accessToken = undefined;
    this.context.refreshToken = undefined;
  }
}
