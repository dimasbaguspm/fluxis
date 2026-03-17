import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetOrg, useUpdateOrg } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { ButtonGroup, Drawer, NoResults } from "@versaur/react/blocks";
import { Banner, Button, Loader } from "@versaur/react/primitive";
import { When } from "@/lib/when";
import { SearchXIcon } from "@versaur/icons";
import { UPDATE_ORGANISATION_FORM_ID } from "./constants";
import { UpdateOrganisationForm } from "./form";
import type { UpdateOrganisationFormInputs } from "./types";

export const UpdateOrganisationDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_ORGANISATION>();
  const orgId = params?.orgId ?? "";

  const [org, , { isLoading: isLoadingOrg }] = useGetOrg(orgId);
  const [updateOrg, err, { isPending }] = useUpdateOrg();

  const onSubmit = async (data: UpdateOrganisationFormInputs) => {
    if (!orgId) return;
    await updateOrg({
      id: orgId,
      data: { name: data.name },
    });
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Update Organisation</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <When condition={isLoadingOrg}>
          <Loader type="bar" />
        </When>
        <When condition={!isLoadingOrg}>
          <When condition={!org}>
            <NoResults
              icon={SearchXIcon}
              title="Organisation not found"
              subtitle="We couldn't find the organisation you're looking for"
            />
          </When>
          <When condition={!!org}>
            <When condition={err?.message}>
              <Banner variant="warning">{err?.message}</Banner>
            </When>
            <UpdateOrganisationForm organisation={org!} onSubmit={onSubmit} />
          </When>
        </When>
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button
            type="submit"
            form={UPDATE_ORGANISATION_FORM_ID}
            loading={isPending || isLoadingOrg}
            disabled={isPending || isLoadingOrg}
          >
            Update
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
