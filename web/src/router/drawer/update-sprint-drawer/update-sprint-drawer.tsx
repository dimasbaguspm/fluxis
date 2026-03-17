import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetSprint, useUpdateSprint } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import dayjs from "dayjs";
import { ButtonGroup, Drawer, NoResults } from "@versaur/react/blocks";
import { Banner, Button, Loader } from "@versaur/react/primitive";
import { When } from "@/lib/when";
import { SearchXIcon } from "@versaur/icons";
import { UPDATE_SPRINT_FORM_ID } from "./constants";
import { UpdateSprintForm } from "./form";
import type { UpdateSprintFormInputs } from "./types";

export const UpdateSprintDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_SPRINT>();
  const sprintId = params?.sprintId ?? "";

  const [sprint, , { isLoading: isLoadingSprint }] = useGetSprint(sprintId);
  const [updateSprint, err, { isPending }] = useUpdateSprint();

  const onSubmit = async (data: UpdateSprintFormInputs) => {
    if (!sprintId) return;
    await updateSprint({
      sprintId: sprintId,
      data: {
        name: data.name,
        goal: data.goal,
        plannedStartedAt: data.plannedStartedAt
          ? dayjs(data.plannedStartedAt).toISOString()
          : undefined,
        plannedCompletedAt: data.plannedCompletedAt
          ? dayjs(data.plannedCompletedAt).toISOString()
          : undefined,
      },
    });
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Update Sprint</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <When condition={isLoadingSprint}>
          <Loader type="bar" />
        </When>
        <When condition={!isLoadingSprint}>
          <When condition={!sprint}>
            <NoResults
              icon={SearchXIcon}
              title="Sprint not found"
              subtitle="We couldn't find the sprint you're looking for"
            />
          </When>
          <When condition={!!sprint}>
            <When condition={err?.message}>
              <Banner variant="warning">{err?.message}</Banner>
            </When>
            <UpdateSprintForm sprint={sprint!} onSubmit={onSubmit} />
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
            form={UPDATE_SPRINT_FORM_ID}
            loading={isPending || isLoadingSprint}
            disabled={isPending || isLoadingSprint}
          >
            Update
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
