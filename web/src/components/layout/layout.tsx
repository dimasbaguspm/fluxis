import { AppLayout } from "@versaur/react/blocks";
import { MobileBreakpoint, TabletOrDesktopBreakpoint } from "@versaur/react/utils";
import { Outlet } from "react-router";
import { SidebarNavigation } from "./sidebar-navigation";
import { BottomBarNavigation } from "./bottom-bar-navigation";

export const Layout = () => {
  return (
    <AppLayout>
      <AppLayout.Body>
        <TabletOrDesktopBreakpoint>
          <AppLayout.SideLeft>
            <SidebarNavigation />
          </AppLayout.SideLeft>
        </TabletOrDesktopBreakpoint>
        <AppLayout.Main>
          <Outlet />
        </AppLayout.Main>
      </AppLayout.Body>
      <MobileBreakpoint>
        <AppLayout.Footer>
          <BottomBarNavigation />
        </AppLayout.Footer>
      </MobileBreakpoint>
    </AppLayout>
  );
};
