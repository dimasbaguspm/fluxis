import { useContext } from "react";
import { DrawerContext } from "./context";
import type { DrawerParams, DrawerProviderModel, DrawerState } from "./types";

export function useDrawer<
  Params extends DrawerParams = DrawerParams,
  State extends DrawerState = DrawerState,
>(): DrawerProviderModel<Params, State> {
  const context = useContext(DrawerContext);
  if (!context) {
    throw new Error("useDrawer must be used within a DrawerProvider");
  }
  return context as DrawerProviderModel<Params, State>;
}
