import { useSessionState } from "@/providers/session";
import { useCallback, useRef } from "react";
import { streamEvents } from "./http-client";

export type StreamEventHandler<TMessage = unknown> = (data: TMessage) => void;
export type StreamErrorHandler = (error: Error) => void;

export interface UseStreamEventsOptions {
  serverTimeout?: number; // default: 300_000 (5 min, matches SERVER_WRITE_TIMEOUT)
  threshold?: number; // default: 30_000 (fire renewal 30s before timeout)
}

export interface StreamEventOptions<TMessage = unknown> {
  onMessage: StreamEventHandler<TMessage>;
  onError?: StreamErrorHandler; // pure observation — hook does NOT reconnect on error
  onClosed?: () => void; // pure observation — fires when stream naturally ends
  onNearlyClosed?: () => void; // pure observation — fires when renewal is triggered
}

/**
 * Hook that returns a function to open authenticated SSE streams with automatic renewal
 * The hook manages connection lifecycle and automatically renews the connection
 * before the server timeout is reached.
 * @param hookOpts - Hook-level options including serverTimeout and threshold
 */
export function useStreamEvents<TMessage = unknown>(hookOpts?: UseStreamEventsOptions) {
  const { accessToken } = useSessionState();
  const abortCurrentRef = useRef<(() => void) | null>(null);
  const renewTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const cancelledRef = useRef(false);

  const connect = useCallback(
    (path: string, options: StreamEventOptions<TMessage>): (() => void) => {
      const serverTimeout = hookOpts?.serverTimeout ?? 300_000;
      const threshold = hookOpts?.threshold ?? 30_000;
      const renewDelay = serverTimeout - threshold;

      const headers: Record<string, string> = {};
      if (accessToken) {
        headers.Authorization = `Bearer ${accessToken}`;
      }

      cancelledRef.current = false;

      const start = () => {
        abortCurrentRef.current = streamEvents<TMessage>(path, options.onMessage, options.onError, {
          headers,
          onClosed: options.onClosed,
        });

        renewTimerRef.current = setTimeout(() => {
          if (cancelledRef.current) return;
          options.onNearlyClosed?.();
          const oldAbort = abortCurrentRef.current;
          start(); // new connection first (no gap)
          oldAbort?.(); // then abort old
        }, renewDelay);
      };

      start();

      return () => {
        cancelledRef.current = true;
        if (renewTimerRef.current !== null) clearTimeout(renewTimerRef.current);
        abortCurrentRef.current?.();
      };
    },
    [accessToken, hookOpts?.serverTimeout, hookOpts?.threshold],
  );

  return connect;
}
