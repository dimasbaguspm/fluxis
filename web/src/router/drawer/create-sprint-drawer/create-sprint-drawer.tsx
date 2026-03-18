import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useCreateSprint, useCreateBoard } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import dayjs from "dayjs";
import { ButtonGroup, Drawer } from "@versaur/react/blocks";
import { Banner, Button } from "@versaur/react/primitive";
import { CREATE_SPRINT_FORM_ID } from "./constants";
import { CreateSprintForm } from "./form";
import type { CreateSprintFormInputs } from "./types";
import { cx } from "@/lib";
import { vMb4 } from "@versaur/core/utilities";

export const CreateSprintDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.CREATE_SPRINT>();
  const [createSprint, sprintErr, { isPending: isSprintPending }] = useCreateSprint();
  const [createBoard, boardErr, { isPending: isBoardPending }] = useCreateBoard();

  const onSubmit = async (data: CreateSprintFormInputs) => {
    const sprintResult = await createSprint({
      projectId: data.projectId,
      name: data.name,
      goal: data.goal,
      plannedStartedAt: data.plannedStartedAt
        ? dayjs(data.plannedStartedAt).toISOString()
        : undefined,
      plannedCompletedAt: data.plannedCompletedAt
        ? dayjs(data.plannedCompletedAt).toISOString()
        : undefined,
    });

    // Create associated board if sprint was created successfully
    if (sprintResult?.data?.id) {
      await createBoard({
        sprintId: sprintResult.data.id,
        name: `${data.name} - Board`,
      });
    }

    closeDrawer();
  };

  const isPending = isSprintPending || isBoardPending;
  const err = sprintErr || boardErr;

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Create Sprint</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        {err && (
          <Banner variant="warning" className={cx(vMb4)}>
            {err?.message}
          </Banner>
        )}
        <CreateSprintForm projectId={params?.projectId ?? ""} onSubmit={onSubmit} />
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button
            type="submit"
            form={CREATE_SPRINT_FORM_ID}
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
