import type { DomainUserModel } from "@/interfaces/openapi.generated";
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { useShallow } from "zustand/shallow";

interface SessionStore {
  user: DomainUserModel | null;
  accessToken: string | null;
  refreshToken: string | null;
  setUser: (user: DomainUserModel | null) => void;
  setTokens: (accessToken: string | null, refreshToken: string | null) => void;
  logout: () => void;
}

export const useSessionStore = create<SessionStore>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,

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

type SessionState = Pick<SessionStore, "accessToken" | "refreshToken" | "user">;

export const useSessionState = (): SessionState =>
  useSessionStore(
    useShallow((state: SessionStore) => ({
      accessToken: state.accessToken,
      refreshToken: state.refreshToken,
      user: state.user,
    })),
  );

type SessionHandler = Pick<SessionStore, "setTokens" | "setUser" | "logout">;

export const useSessionHandler = (): SessionHandler =>
  useSessionStore(
    useShallow((state: SessionStore) => ({
      setTokens: state.setTokens,
      setUser: state.setUser,
      logout: state.logout,
    })),
  );
