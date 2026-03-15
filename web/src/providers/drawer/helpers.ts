import type { DrawerParams } from "./types";

export const DRAWER_QUERY_KEY = "drawer";

export interface ParsedDrawer {
  drawerId: string | null;
  params: DrawerParams;
}

/**
 * Format a drawer ID and optional params into a URL-safe string.
 * Format: "drawerId" or "drawerId:base64(JSON.stringify(params))"
 */
export function formatDrawerForUrl(drawerId: string, params?: DrawerParams): string {
  if (!params || Object.keys(params).length === 0) {
    return drawerId;
  }
  const encoded = btoa(JSON.stringify(params));
  return `${drawerId}:${encoded}`;
}

/**
 * Parse a drawer URL value back into drawerId and params.
 * Returns { drawerId: null, params: null } if the value is empty or malformed.
 */
export function parseDrawerFromUrl(value: string | null): ParsedDrawer {
  if (!value) {
    return { drawerId: null, params: null };
  }

  const [drawerId, encoded] = value.split(":");

  if (!drawerId) {
    return { drawerId: null, params: null };
  }

  if (!encoded) {
    return { drawerId, params: null };
  }

  try {
    const params = JSON.parse(atob(encoded));
    return { drawerId, params };
  } catch {
    // Malformed encoding, treat as closed
    return { drawerId: null, params: null };
  }
}
