import { useListTickets } from "@/hooks/use-api";
import { dropTargetForElements } from "@atlaskit/pragmatic-drag-and-drop/element/adapter";
import { useEffect, useRef, useState } from "react";
import { TicketGroupHeader } from "./ticket-group-header";
import { TicketsTable } from "./tickets-table";

interface BacklogGroupProps {
  projectId: string;
  onEditTicket: (ticketId: string) => void;
  onMoveTicket?: (ticketId: string, targetSprintId: string | null) => void;
  onCreateTicket?: (sprintId?: string) => void;
}

export const BacklogGroup = ({
  projectId,
  onEditTicket,
  onMoveTicket,
  onCreateTicket,
}: BacklogGroupProps) => {
  const [tickets, error, { isLoading }] = useListTickets({
    projectId: [projectId],
    // @ts-ignore
    sprintId: "",
    pageSize: 50,
  });
  const dropRef = useRef<HTMLDivElement>(null);
  const [isDragOver, setIsDragOver] = useState(false);

  // Filter for tickets without sprintId (backlog items)
  const ticketsList = (tickets?.items || []).filter((ticket) => !ticket.sprintId);
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
        if (data.ticketId && onMoveTicket && data.sourceGroupId !== null) {
          onMoveTicket(data.ticketId, null);
        }
      },
    });
  }, [onMoveTicket]);

  return (
    <div
      ref={dropRef}
      style={{
        border: "1px solid #e0e0e0",
        borderRadius: "6px",
        marginBottom: "1.5rem",
        overflow: "hidden",
        opacity: isDragOver ? 0.8 : 1,
        transition: "opacity 0.2s",
        backgroundColor: isDragOver ? "#f5f5f5" : "transparent",
      }}
    >
      <TicketGroupHeader
        title="Backlog"
        description="Tickets not assigned to any sprint"
        totalStoryPoints={totalStoryPoints}
        onCreateTicket={() => onCreateTicket?.()}
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
          groupId={null}
        />
      )}
    </div>
  );
};
