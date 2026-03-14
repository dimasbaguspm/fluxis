import { useCurrentUser } from "@hooks/use-api/api/use-user";
import { useEffect } from "react";
import { useSessionHandler } from "./use-session-store";

/**
 * Hydrates the session store with the current user from the API.
 *
 * - On success: sets the user in the session store
 * - On error: clears the session store (logout)
 *
 * Should be rendered at the top of protected routes.
 */
export function SessionHydrator({ children }: { children: React.ReactNode }) {
  const [data, err, { isPending }] = useCurrentUser();
  const { setUser, logout } = useSessionHandler();

  useEffect(() => {
    if (data) {
      setUser(data);
    } else if (err) {
      logout();
    }
  }, [data, err, setUser, logout]);

  if (isPending) return null;

  return <>{children}</>;
}
