import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetTicket, useUpdateTicket } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { ButtonGroup, Drawer, NoResults } from "@versaur/react/blocks";
import { Banner, Button, Loader } from "@versaur/react/primitive";
import { When } from "@/lib/when";
import { SearchXIcon } from "@versaur/icons";
import { UPDATE_TICKET_FORM_ID } from "./constants";
import { UpdateTicketForm } from "./form";
import type { UpdateTicketFormInputs } from "./types";

export const UpdateTicketDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_TICKET>();
  const ticketId = params?.ticketId ?? "";

  const [ticket, , { isLoading: isLoadingTicket }] = useGetTicket(ticketId);
  const [updateTicket, err, { isPending }] = useUpdateTicket();

  const onSubmit = async (data: UpdateTicketFormInputs) => {
    if (!ticketId) return;
    await updateTicket({
      ticketId: ticketId,
      data: {
        title: data.title,
        type: data.type as "bug" | "story" | "task" | "epic" | undefined,
        priority: data.priority as "low" | "medium" | "high" | "critical" | undefined,
        description: data.description,
        sprintId: data.sprintId,
        storyPoints: data.storyPoints,
        dueDate: data.dueDate,
      },
    });
    closeDrawer();
  };

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
            <When condition={err?.message}>
              <Banner variant="warning">{err?.message}</Banner>
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
