/**
 * HTTP Request abstraction for type-safe API calls
 * Decouples request building from HTTP client implementation
 */

export type HttpRequestMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export interface HttpRequest {
  method: HttpRequestMethod;
  path: string;
  body?: unknown;
  params?: unknown;
  responseType?: "json" | "blob" | "text";
}

export interface HttpResponse<T = unknown> {
  status: number;
  data: T;
  headers: Record<string, string>;
}

export function createHttpError(status: number, data: unknown, message: string): HttpError {
  return new HttpError(status, data, message);
}

class HttpError extends Error {
  status: number;
  data: unknown;

  constructor(status: number, data: unknown, message: string) {
    super(message);
    this.name = "HttpError";
    this.status = status;
    this.data = data;
  }
}
