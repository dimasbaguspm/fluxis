/**
 * API Response wrapper
 */
export interface APIResponse<T = unknown> {
  data?: T;
  error?: ErrorModel;
  status: number;
  headers: Record<string, string>;
}

/**
 * Error model from API
 */
export interface ErrorModel {
  title?: string;
  status?: number;
  detail?: string;
  errors?: ErrorDetail[];
}

export interface ErrorDetail {
  location?: string;
  message?: string;
  value?: unknown;
}

/**
 * Test context that carries auth and other shared state
 */
export interface TestContext {
  accessToken?: string;
  refreshToken?: string;
  baseURL: string;
}
