export interface NotificationEvent {
  type: string;
  payload: unknown;
}

export interface NotifierContextType {
  event: NotificationEvent | null;
  isConnected: boolean;
}
