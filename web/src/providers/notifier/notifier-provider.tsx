import { createContext, useContext, useEffect, useState } from "react";
import { useStreamEvents } from "@/hooks/use-fetcher";
import type { NotificationEvent, NotifierContextType } from "./types";

const NotifierContext = createContext<NotifierContextType | null>(null);
NotifierContext.displayName = "NotifierContext";

interface NotifierProviderProps {
  children: React.ReactNode;
}

export function NotifierProvider({ children }: NotifierProviderProps) {
  const [event, setEvent] = useState<NotificationEvent | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const connect = useStreamEvents();

  useEffect(() => {
    console.log("[Notifier] Setting up stream connection");
    const unsubscribe = connect<NotificationEvent>("/notifications", {
      onMessage: (notification: unknown) => {
        const event = notification as NotificationEvent;
        console.log("[Notifier] Received event:", event.type);
        setEvent(event);
      },
      onError: (error) => {
        console.error("[Notifier] Stream error:", error);
        setIsConnected(false);
      },
    });

    setIsConnected(true);
    console.log("[Notifier] Stream connected");

    return () => {
      console.log("[Notifier] Cleaning up stream");
      unsubscribe();
      setIsConnected(false);
    };
  }, [connect]);

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
