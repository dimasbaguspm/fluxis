import { useContext } from "react";
import type { GetDrawerParams } from "@constants/drawer-params";
import { DrawerContext } from "./context";
import type { DrawerProviderModel, DrawerState } from "./types";

export function useDrawer<DrawerId extends string = string>(
): DrawerProviderModel<DrawerId, GetDrawerParams<DrawerId>, DrawerState> {
  const context = useContext(DrawerContext);
  if (!context) {
    throw new Error("useDrawer must be used within a DrawerProvider");
  }
  return context as DrawerProviderModel<DrawerId, GetDrawerParams<DrawerId>, DrawerState>;
}
