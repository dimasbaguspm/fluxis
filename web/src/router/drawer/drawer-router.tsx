import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useDrawer } from "@providers/drawer";
import { Drawer } from "@versaur/react/blocks";
import { CreateOrganisationDrawer } from "./create-organisation-drawer";
import { UpdateOrganisationDrawer } from "./update-organisation-drawer";
import { CreateProjectDrawer } from "./create-project-drawer";
import { UpdateProjectDrawer } from "./update-project-drawer";
import { CreateSprintDrawer } from "./create-sprint-drawer";
import { UpdateSprintDrawer } from "./update-sprint-drawer";
import { CreateBoardDrawer } from "./create-board-drawer";
import { UpdateBoardDrawer } from "./update-board-drawer";

export function DrawerRouter() {
  const { isOpen, drawerId, closeDrawer } = useDrawer();

  function renderContent() {
    if (drawerId === DRAWER_ROUTES.CREATE_ORGANISATION) return <CreateOrganisationDrawer />;
    if (drawerId === DRAWER_ROUTES.UPDATE_ORGANISATION) return <UpdateOrganisationDrawer />;
    if (drawerId === DRAWER_ROUTES.CREATE_PROJECT) return <CreateProjectDrawer />;
    if (drawerId === DRAWER_ROUTES.UPDATE_PROJECT) return <UpdateProjectDrawer />;
    if (drawerId === DRAWER_ROUTES.CREATE_SPRINT) return <CreateSprintDrawer />;
    if (drawerId === DRAWER_ROUTES.UPDATE_SPRINT) return <UpdateSprintDrawer />;
    if (drawerId === DRAWER_ROUTES.CREATE_BOARD) return <CreateBoardDrawer />;
    if (drawerId === DRAWER_ROUTES.UPDATE_BOARD) return <UpdateBoardDrawer />;
    return null;
  }

  return (
    <Drawer
      open={isOpen}
      placement="right"
      onOpenChange={() => {
        if (isOpen) {
          closeDrawer();
        }
      }}
    >
      {renderContent()}
    </Drawer>
  );
}
