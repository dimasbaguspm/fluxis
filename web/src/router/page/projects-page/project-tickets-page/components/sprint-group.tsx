import { useRef, useEffect, useState } from "react";
import { useListTickets } from "@/hooks/use-api";
import type { DomainSprintModel } from "@/interfaces/openapi.generated";
import { dropTargetForElements } from "@atlaskit/pragmatic-drag-and-drop/element/adapter";
import { TicketGroupHeader } from "./ticket-group-header";
import { TicketsTable } from "./tickets-table";

interface SprintGroupProps {
  sprint: DomainSprintModel;
  onEditTicket: (ticketId: string) => void;
  onMoveTicket?: (ticketId: string, targetSprintId: string | null) => void;
}

export const SprintGroup = ({ sprint, onEditTicket, onMoveTicket }: SprintGroupProps) => {
  const [tickets, error, { isLoading }] = useListTickets({
    projectId: [sprint.projectId],
    sprintId: [sprint.id],
    pageSize: 50,
  });
  const dropRef = useRef<HTMLDivElement>(null);
  const [isDragOver, setIsDragOver] = useState(false);

  const ticketsList = tickets?.items || [];
  const totalStoryPoints = ticketsList.reduce((sum, ticket) => sum + (ticket.storyPoints || 0), 0);

  useEffect(() => {
    if (!dropRef.current) return;

    return dropTargetForElements({
      element: dropRef.current,
      onDrag: () => {
        setIsDragOver(true);
      },
      onDragLeave: () => {
        setIsDragOver(false);
      },
      onDrop: ({ source }: any) => {
        setIsDragOver(false);
        const data = source.data as { ticketId: string; sourceGroupId?: string };
        if (data.ticketId && onMoveTicket && data.sourceGroupId !== sprint.id) {
          onMoveTicket(data.ticketId, sprint.id);
        }
      },
    });
  }, [sprint.id, onMoveTicket]);

  return (
    <div
      ref={dropRef}
      style={{
        marginBottom: "2rem",
        opacity: isDragOver ? 0.8 : 1,
        transition: "opacity 0.2s",
        backgroundColor: isDragOver ? "#f5f5f5" : "transparent",
      }}
    >
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
          groupId={sprint.id}
        />
      )}
    </div>
  );
};
