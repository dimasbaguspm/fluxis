import { BACKEND_EVENT_MAP, NotifierAction, NotifierEventType } from "./constants";
import type { NotifierEvent } from "./types";

/**
 * Translates raw SSE event to typed NotifierEvent
 * Maps backend event string to type-safe type and action values
 */
export function translateEvent(raw: {
  type: string;
  payload: unknown;
}): NotifierEvent {
  const mapped =
    BACKEND_EVENT_MAP[raw.type] ??
    ({
      type: NotifierEventType.Unknown,
      action: NotifierAction.Unknown,
    } as const);

  return {
    type: mapped.type,
    action: mapped.action,
    payload: raw.payload,
  };
}
