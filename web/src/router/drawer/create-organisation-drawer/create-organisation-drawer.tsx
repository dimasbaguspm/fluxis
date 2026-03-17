import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useCreateOrg } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { ButtonGroup, Drawer } from "@versaur/react/blocks";
import { Banner, Button } from "@versaur/react/primitive";
import { CREATE_ORGANISATION_FORM_ID } from "./constants";
import { CreateOrganisationForm } from "./form";
import type { CreateOrganisationFormInputs } from "./types";
import { cx } from "@/lib";
import { vMb4 } from "@versaur/core/utilities";

export const CreateOrganisationDrawer = () => {
  const { closeDrawer } = useDrawer<typeof DRAWER_ROUTES.CREATE_ORGANISATION>();
  const [createOrg, err, { isPending }] = useCreateOrg();

  const onSubmit = async (data: CreateOrganisationFormInputs) => {
    await createOrg(data);
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Create Organisation</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        {err && (
          <Banner variant="warning" className={cx(vMb4)}>
            {err.message}
          </Banner>
        )}
        <CreateOrganisationForm onSubmit={onSubmit} />
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button
            type="submit"
            form={CREATE_ORGANISATION_FORM_ID}
            loading={isPending}
            disabled={isPending}
          >
            Create
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
