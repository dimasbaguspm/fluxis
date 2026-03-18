import { useListTickets } from "@/hooks/use-api";
import { TicketGroupHeader } from "./ticket-group-header";
import { TicketsTable } from "./tickets-table";

interface BacklogGroupProps {
  projectId: string;
  onEditTicket: (ticketId: string) => void;
}

export const BacklogGroup = ({ projectId, onEditTicket }: BacklogGroupProps) => {
  const [tickets, error, { isLoading }] = useListTickets({
    projectId: [projectId],
    // @ts-ignore
    sprintId: "",
    pageSize: 50,
  });

  // Filter for tickets without sprintId (backlog items)
  const ticketsList = (tickets?.items || []).filter((ticket) => !ticket.sprintId);
  const totalStoryPoints = ticketsList.reduce((sum, ticket) => sum + (ticket.storyPoints || 0), 0);

  if (ticketsList.length === 0) {
    return null;
  }

  return (
    <div style={{ marginBottom: "2rem" }}>
      <TicketGroupHeader
        title="Backlog"
        description="Tickets not assigned to any sprint"
        ticketCount={ticketsList.length}
        totalStoryPoints={totalStoryPoints}
      />
      {error ? (
        <div style={{ padding: "2rem", textAlign: "center", color: "#d32f2f" }}>
          Error loading backlog tickets
        </div>
      ) : (
        <TicketsTable
          tickets={ticketsList}
          onEditTicket={onEditTicket}
          isLoading={isLoading}
        />
      )}
    </div>
  );
};
