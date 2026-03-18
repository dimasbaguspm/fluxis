import { useListTickets } from "@/hooks/use-api";
import type { DomainSprintModel } from "@/interfaces/openapi.generated";
import { TicketGroupHeader } from "./ticket-group-header";
import { TicketsTable } from "./tickets-table";

interface SprintGroupProps {
  sprint: DomainSprintModel;
  onEditTicket: (ticketId: string) => void;
}

export const SprintGroup = ({ sprint, onEditTicket }: SprintGroupProps) => {
  const [tickets, error, { isLoading }] = useListTickets({
    projectId: [sprint.projectId],
    sprintId: [sprint.id],
    pageSize: 50,
  });

  const ticketsList = tickets?.items || [];
  const totalStoryPoints = ticketsList.reduce((sum, ticket) => sum + (ticket.storyPoints || 0), 0);

  return (
    <div style={{ marginBottom: "2rem" }}>
      <TicketGroupHeader
        title={sprint.name}
        ticketCount={ticketsList.length}
        totalStoryPoints={totalStoryPoints}
        sprint={sprint}
      />
      {error ? (
        <div style={{ padding: "2rem", textAlign: "center", color: "#d32f2f" }}>
          Error loading tickets for this sprint
        </div>
      ) : (
        <TicketsTable
          tickets={ticketsList}
          onEditTicket={onEditTicket}
          isLoading={isLoading}
          sprint={sprint}
        />
      )}
    </div>
  );
};
