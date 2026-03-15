import type { ReactNode } from "react";
import { useCallback, useMemo } from "react";
import { useLocation, useNavigate, useSearchParams } from "react-router";
import { DrawerContext } from "./context";
import { DRAWER_QUERY_KEY, formatDrawerForUrl, parseDrawerFromUrl } from "./helpers";
import type { CloseDrawerFunc, DrawerProviderModel, OpenDrawerFunc } from "./types";

interface DrawerProviderProps {
  children: ReactNode;
}

export function DrawerProvider({ children }: DrawerProviderProps) {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const location = useLocation();

  const drawerValue = searchParams.get(DRAWER_QUERY_KEY);
  const parsed = parseDrawerFromUrl(drawerValue);

  const openDrawer: OpenDrawerFunc = useCallback(
    (id, params, opts) => {
      const newParams = new URLSearchParams(searchParams);
      newParams.set(DRAWER_QUERY_KEY, formatDrawerForUrl(id, params || null));

      navigate(
        {
          pathname: location.pathname,
          search: newParams.toString(),
        },
        { preventScrollReset: true, replace: opts?.replace },
      );
    },
    [navigate, searchParams, location.pathname],
  );

  const closeDrawer: CloseDrawerFunc = useCallback(() => {
    console.log("closedrawer");
    navigate(-1);
  }, [navigate]);

  const value: DrawerProviderModel = useMemo(
    () => ({
      isOpen: parsed.drawerId !== null,
      drawerId: parsed.drawerId,
      params: parsed.params,
      state: null, // TODO: implement state management if needed
      openDrawer,
      closeDrawer,
    }),
    [parsed, openDrawer, closeDrawer],
  );

  return <DrawerContext.Provider value={value}>{children}</DrawerContext.Provider>;
}
