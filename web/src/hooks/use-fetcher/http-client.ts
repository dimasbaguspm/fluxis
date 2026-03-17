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

export interface StreamOptions extends RequestOptions {
  onClosed?: () => void;
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
      data = await response.json();
  }

  if (!response.ok) {
    const errorMessage =
      (data as any)?.error?.message || `HTTP ${response.status}: ${response.statusText}`;
    throw createHttpError(response.status, data, errorMessage);
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
 * Open a Server-Sent Events connection with authentication support via fetch + ReadableStream
 * @param path - API endpoint path
 * @param onMessage - Callback when message is received
 * @param onError - Callback when error occurs
 * @param options - Request options including headers for authentication and onClosed callback
 * @returns Function to close the connection
 */
export function streamEvents<TMessage = unknown>(
  path: string,
  onMessage: (data: TMessage) => void,
  onError?: (error: Error) => void,
  options?: StreamOptions,
): () => void {
  const url = `${API_BASE_URL}${path}`;
  const abortController = new AbortController();

  const headers: Record<string, string> = {
    Accept: "text/event-stream",
    ...options?.headers,
  };

  const connectStream = async () => {
    try {
      const response = await fetch(url, {
        method: "GET",
        headers,
        signal: abortController.signal,
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const reader = response.body?.getReader();
      if (!reader) {
        throw new Error("Response body is not readable");
      }

      const decoder = new TextDecoder();
      let buffer = "";
      let eventData = "";

      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          options?.onClosed?.();
          break;
        }

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n");
        buffer = lines.pop() || "";

        for (const line of lines) {
          const trimmed = line.trim();

          // Blank line signals end of event
          if (!trimmed) {
            if (eventData) {
              try {
                const data = JSON.parse(eventData) as TMessage;
                onMessage(data);
              } catch (error) {
                const parseError = new Error(
                  error instanceof Error ? error.message : "Failed to parse SSE message",
                );
                onError?.(parseError);
              }
              eventData = "";
            }
            continue;
          }

          // Skip comments and event type declarations
          if (trimmed.startsWith(":") || trimmed.startsWith("event:")) {
            continue;
          }

          // Accumulate data lines
          if (trimmed.startsWith("data:")) {
            const dataValue = trimmed.slice(5).trim();
            eventData += dataValue;
          }
        }
      }
    } catch (error) {
      if (error instanceof Error && error.name !== "AbortError") {
        onError?.(error);
      }
    }
  };

  connectStream();

  return () => abortController.abort();
}
