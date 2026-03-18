import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetTicket, useMoveToSprint, useUpdateTicket } from "@/hooks/use-api";
import { dateFormat, FormatDate } from "@/lib";
import { When } from "@/lib/when";
import { useDrawer } from "@/providers/drawer";
import { SearchXIcon } from "@versaur/icons";
import { ButtonGroup, Drawer, NoResults } from "@versaur/react/blocks";
import { Banner, Button, Loader } from "@versaur/react/primitive";
import { UPDATE_TICKET_FORM_ID } from "./constants";
import { UpdateTicketForm } from "./form";
import type { UpdateTicketFormInputs } from "./types";

export const UpdateTicketDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_TICKET>();
  const ticketId = params?.ticketId ?? "";

  const [ticket, , { isLoading: isLoadingTicket }] = useGetTicket(ticketId);
  const [updateTicket, err, { isPending: isUpdateTicketPending }] = useUpdateTicket();
  const [moveToSprint, errMoveToSprint, { isPending: isMoveToSprintPending }] = useMoveToSprint();

  const onSubmit = async (data: UpdateTicketFormInputs) => {
    if (!ticketId) return;
    await updateTicket({
      ticketId: ticketId,
      data: {
        title: data.title,
        type: data.type as "bug" | "story" | "task" | "epic" | undefined,
        priority: data.priority as "low" | "medium" | "high" | "critical" | undefined,
        description: data.description,
        sprintId: data?.sprintId ? data.sprintId : undefined,
        storyPoints: data?.storyPoints ? data.storyPoints : undefined,
        dueDate: data?.dueDate ? dateFormat(data.dueDate, FormatDate.ISO) : undefined,
      },
    });

    if (data.sprintId && data.sprintId !== ticket?.sprintId) {
      await moveToSprint({ ticketId, move: { sprintId: data.sprintId } });
    } else if (!data.sprintId && ticket?.sprintId) {
      await moveToSprint({ ticketId, move: { sprintId: undefined } });
    }
    closeDrawer();
  };

  const errMessage = err?.message || errMoveToSprint?.message;
  const isPending = isUpdateTicketPending || isMoveToSprintPending;

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Update Ticket</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <When condition={isLoadingTicket}>
          <Loader type="bar" />
        </When>
        <When condition={!isLoadingTicket}>
          <When condition={!ticket}>
            <NoResults
              icon={SearchXIcon}
              title="Ticket not found"
              subtitle="We couldn't find the ticket you're looking for"
            />
          </When>
          <When condition={!!ticket}>
            <When condition={errMessage}>
              <Banner variant="warning">{errMessage}</Banner>
            </When>
            <UpdateTicketForm ticket={ticket!} onSubmit={onSubmit} />
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
            form={UPDATE_TICKET_FORM_ID}
            loading={isPending || isLoadingTicket}
            disabled={isPending || isLoadingTicket}
          >
            Update
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
