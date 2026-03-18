import type { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useCreateTicket, useMoveToSprint } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { ButtonGroup, Drawer } from "@versaur/react/blocks";
import { Banner, Button } from "@versaur/react/primitive";
import { CREATE_TICKET_FORM_ID } from "./constants";
import { CreateTicketForm } from "./form";
import type { CreateTicketFormInputs } from "./types";

export const CreateTicketDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.CREATE_TICKET>();

  const [createTicket, errCreateTicket, { isPending: isCreateTicketPending }] = useCreateTicket();
  const [moveToSprint, errMoveToSprint, { isPending: isMoveToSprintPending }] = useMoveToSprint();

  const onSubmit = async (data: CreateTicketFormInputs) => {
    const ticket = await createTicket({
      projectId: data.projectId,
      data: {
        projectId: data.projectId,
        title: data.title,
        type: data.type as "bug" | "story" | "task" | "epic",
        priority: data.priority as "low" | "medium" | "high" | "critical",
        description: data.description,
        sprintId: data.sprintId ? data.sprintId : undefined,
        storyPoints: data.storyPoints ? data.storyPoints : undefined,
        dueDate: data.dueDate ? data.dueDate : undefined,
      },
    });

    if (data.sprintId) {
      await moveToSprint({ ticketId: ticket.id, move: { sprintId: data.sprintId } });
    }
    closeDrawer();
  };

  const errMessage = errCreateTicket?.message || errMoveToSprint?.message;
  const isPending = isCreateTicketPending || isMoveToSprintPending;

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Create Ticket</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        {errMessage && <Banner variant="warning">{errMessage}</Banner>}
        <CreateTicketForm
          projectId={params?.projectId}
          sprintId={params?.sprintId}
          onSubmit={onSubmit}
        />
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button
            type="submit"
            form={CREATE_TICKET_FORM_ID}
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
