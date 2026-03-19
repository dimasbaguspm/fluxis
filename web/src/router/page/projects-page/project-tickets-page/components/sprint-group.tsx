import { useListTickets } from "@/hooks/use-api";
import type { DomainSprintModel } from "@/interfaces/openapi.generated";
import { dropTargetForElements } from "@atlaskit/pragmatic-drag-and-drop/element/adapter";
import { useEffect, useRef, useState } from "react";
import { TicketGroupHeader } from "./ticket-group-header";
import { TicketsTable } from "./tickets-table";

interface SprintGroupProps {
  sprint: DomainSprintModel;
  onEditTicket: (ticketId: string) => void;
  onEditSprint?: (sprintId: string) => void;
  onMoveTicket?: (ticketId: string, targetSprintId: string | null) => void;
  onCreateTicket?: (sprintId?: string) => void;
  onStartSprint?: (sprintId: string) => void;
  onCompleteSprint?: (sprintId: string) => void;
}

export const SprintGroup = ({
  sprint,
  onEditTicket,
  onEditSprint,
  onMoveTicket,
  onCreateTicket,
  onStartSprint,
  onCompleteSprint,
}: SprintGroupProps) => {
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
        border: "1px solid var(--color-border)",
        marginBottom: "1rem",
        overflow: "hidden",
        ...(sprint.status === "active" && {
          borderLeft: "3px solid var(--color-primary, #1976d2)",
        }),
        opacity: isDragOver ? 0.8 : 1,
        transition: "opacity 0.2s",
        backgroundColor: isDragOver ? "#f5f5f5" : "var(--color-background)",
      }}
    >
      <TicketGroupHeader
        title={sprint.name}
        description={sprint.goal}
        totalStoryPoints={totalStoryPoints}
        sprint={sprint}
        onCreateTicket={onCreateTicket ? () => onCreateTicket(sprint.id) : undefined}
        onEditSprint={onEditSprint ? () => onEditSprint(sprint.id) : undefined}
        onStartSprint={onStartSprint ? () => onStartSprint(sprint.id) : undefined}
        onCompleteSprint={onCompleteSprint ? () => onCompleteSprint(sprint.id) : undefined}
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
