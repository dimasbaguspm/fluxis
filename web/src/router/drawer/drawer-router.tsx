import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useDrawer } from "@providers/drawer";
import { Drawer } from "@versaur/react/blocks";
import { CreateOrganisationDrawer } from "./create-organisation-drawer";

export function DrawerRouter() {
  const { isOpen, drawerId, closeDrawer } = useDrawer();

  function renderContent() {
    if (drawerId === DRAWER_ROUTES.CREATE_ORGANISATION) return <CreateOrganisationDrawer />;
    return null;
  }

  return (
    <Drawer
      open={isOpen}
      placement="right"
      onOpenChange={(open) => {
        if (isOpen) {
          closeDrawer();
        }
      }}
    >
      {renderContent()}
    </Drawer>
  );
}
