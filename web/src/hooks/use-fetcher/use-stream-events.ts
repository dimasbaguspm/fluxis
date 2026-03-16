import { useSessionState } from "@/providers/session";
import { useCallback } from "react";
import { streamEvents } from "./http-client";

export type StreamEventHandler<TMessage = unknown> = (data: TMessage) => void;
export type StreamErrorHandler = (error: Error) => void;

export interface StreamEventOptions {
  onMessage: StreamEventHandler;
  onError?: StreamErrorHandler;
}

/**
 * Hook that returns a function to open authenticated SSE streams
 * Automatically includes the Authorization Bearer token in headers
 */
export function useStreamEvents() {
  const { accessToken } = useSessionState();

  const connect = useCallback(
    <TMessage = unknown>(
      path: string,
      { onMessage, onError }: StreamEventOptions,
    ): (() => void) => {
      const headers: Record<string, string> = {};
      if (accessToken) {
        headers.Authorization = `Bearer ${accessToken}`;
      }

      return streamEvents<TMessage>(path, onMessage, onError, { headers });
    },
    [accessToken],
  );

  return connect;
}
