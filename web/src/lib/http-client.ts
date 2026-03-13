import { createHttpError, type HttpRequest, type HttpResponse } from "./http-request";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

function serializeParams(params: Record<string, any>): string {
  const searchParams = new URLSearchParams();

  for (const [key, value] of Object.entries(params)) {
    if (value === null || value === undefined) {
      continue;
    }

    if (Array.isArray(value)) {
      // For arrays, add multiple entries with the same key
      for (const item of value) {
        if (item !== null && item !== undefined) {
          searchParams.append(key, String(item));
        }
      }
    } else {
      searchParams.set(key, String(value));
    }
  }

  return searchParams.toString();
}

export interface RequestOptions {
  headers?: Record<string, string>;
}

export async function request<TResponse = unknown>(
  req: HttpRequest,
  options?: RequestOptions,
): Promise<HttpResponse<TResponse>> {
  let url = `${API_BASE_URL}${req.path}`;
  if (req.params) {
    const query = serializeParams(req.params);
    if (query) {
      url += `?${query}`;
    }
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...options?.headers,
  };

  const fetchOptions: RequestInit = {
    method: req.method,
    headers,
    credentials: "include",
  };

  if (req.body && (req.method === "POST" || req.method === "PUT" || req.method === "PATCH")) {
    fetchOptions.body = JSON.stringify(req.body);
  }

  // Execute request
  const response = await fetch(url, fetchOptions);

  // Parse response
  let data: any;
  const responseType = req.responseType || "json";

  switch (responseType) {
    case "blob":
      data = await response.blob();
      break;
    case "text":
      data = await response.text();
      break;
    case "json":
    default:
      data = response.ok ? await response.json() : null;
  }

  // Handle errors
  if (!response.ok) {
    throw createHttpError(response.status, data, `HTTP ${response.status}: ${response.statusText}`);
  }

  // Build response headers
  const responseHeaders: Record<string, string> = {};
  response.headers.forEach((value, key) => {
    responseHeaders[key] = value;
  });

  return {
    status: response.status,
    data,
    headers: responseHeaders,
  };
}

/**
 * Open a Server-Sent Events connection
 * @param path - API endpoint path
 * @param onMessage - Callback when message is received
 * @param onError - Callback when error occurs
 * @param options - Request options including headers for authentication
 * @returns Function to close the connection
 */
export function streamEvents<TMessage = unknown>(
  path: string,
  onMessage: (data: TMessage) => void,
  onError?: (error: Error) => void,
): () => void {
  const url = `${API_BASE_URL}${path}`;

  const eventSource = new EventSource(url);

  eventSource.addEventListener("message", (event) => {
    try {
      const data = JSON.parse(event.data) as TMessage;
      onMessage(data);
    } catch (error) {
      const parseError = new Error(
        error instanceof Error ? error.message : "Failed to parse SSE message",
      );
      onError?.(parseError);
    }
  });

  eventSource.addEventListener("error", () => {
    const error = new Error("SSE connection error");
    onError?.(error);
    eventSource.close();
  });

  return () => eventSource.close();
}
