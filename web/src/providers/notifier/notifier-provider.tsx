import { createContext, useContext, useEffect, useState } from "react";
import { useStreamEvents } from "@/hooks/use-fetcher";
import type { NotifierContextType, NotifierEvent } from "./types";
import { translateEvent } from "./helpers";

const NotifierContext = createContext<NotifierContextType | null>(null);
NotifierContext.displayName = "NotifierContext";

interface NotifierProviderProps {
  children: React.ReactNode;
  onNotified?: (event: NotifierEvent) => void;
}

export function NotifierProvider({ children, onNotified }: NotifierProviderProps) {
  const [event, setEvent] = useState<NotifierEvent | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const connect = useStreamEvents<{ type: string; payload: unknown }>({
    serverTimeout: 300_000,
    threshold: 30_000,
  });

  useEffect(() => {
    const unsubscribe = connect("/notifications", {
      onMessage: (notification) => {
        const translatedEvent = translateEvent(notification);
        setEvent(translatedEvent);
        onNotified?.(translatedEvent);
      },
      onError: (error) => {
        console.error("[Notifier] Stream error:", error);
        setIsConnected(false);
      },
      onClosed: () => setIsConnected(false),
    });

    setIsConnected(true);

    return () => {
      unsubscribe();
      setIsConnected(false);
    };
  }, [connect, onNotified]);

  const value: NotifierContextType = {
    event,
    isConnected,
  };

  return <NotifierContext.Provider value={value}>{children}</NotifierContext.Provider>;
}

export function useNotifier(): NotifierContextType {
  const context = useContext(NotifierContext);
  if (!context) {
    throw new Error("useNotifier must be used within a NotifierProvider");
  }
  return context;
}
