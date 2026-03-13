import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { Session, User } from "./types";

interface SessionStore {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  setSession: (session: Session | null) => void;
  setUser: (user: User | null) => void;
  setTokens: (accessToken: string | null, refreshToken: string | null) => void;
  logout: () => void;
}

/**
 * Session store with localStorage persistence
 * Manages user authentication state across the app
 */
export const useSessionStore = create<SessionStore>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,

      setSession: (session) => {
        set({
          user: session?.user ?? null,
          accessToken: session?.accessToken ?? null,
          refreshToken: session?.refreshToken ?? null,
        });
      },

      setUser: (user) => set({ user }),
      setTokens: (accessToken, refreshToken) => set({ accessToken, refreshToken }),

      logout: () => {
        set({ user: null, accessToken: null, refreshToken: null });
      },
    }),
    {
      name: "fluxis:session",
    },
  ),
);
