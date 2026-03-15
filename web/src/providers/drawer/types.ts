import type { DrawerParamsMap, GetDrawerParams } from "@constants/drawer-params";

export type DrawerParams = Record<string, string | number> | null;
export type DrawerState = Record<string, unknown> | null;

export interface OpenDrawerOptions {
  replace?: boolean;
  state?: DrawerState;
}

export type OpenDrawerFunc = <DrawerId extends string, Params extends DrawerParams = null>(
  id: DrawerId,
  params?: Params extends null ? GetDrawerParams<DrawerId> : Params,
  opts?: OpenDrawerOptions,
) => void;

export type CloseDrawerFunc = () => void;

export interface DrawerProviderModel<
  DrawerId extends string = string,
  Params = DrawerId extends keyof DrawerParamsMap ? DrawerParamsMap[DrawerId] : DrawerParams,
  State = DrawerState,
> {
  isOpen: boolean;
  drawerId: DrawerId | null;
  params: Params;
  state: State;
  openDrawer: OpenDrawerFunc;
  closeDrawer: CloseDrawerFunc;
}
