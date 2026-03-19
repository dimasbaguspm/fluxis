import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useListTickets } from "@/hooks/use-api";
import type {
  DomainBoardColumnModel,
  DomainTicketMoveBoardColumnModel,
} from "@/interfaces/openapi.generated";
import { useDrawer } from "@/providers/drawer";
import { dropTargetForElements } from "@atlaskit/pragmatic-drag-and-drop/element/adapter";
import { Button, Text } from "@versaur/react/primitive";
import { useEffect, useRef, useState } from "react";
import { TicketCard } from "./ticket-card";

interface BoardColumnProps {
  column: DomainBoardColumnModel;
  projectId: string;
  onMoveTicket: (variables: {
    ticketId: string;
    move: DomainTicketMoveBoardColumnModel;
  }) => void;
}

export const BoardColumn = ({
  column,
  projectId,
  onMoveTicket,
}: BoardColumnProps) => {
  const { openDrawer } = useDrawer();
  const dropRef = useRef<HTMLDivElement>(null);
  const [isDragOver, setIsDragOver] = useState(false);

  const handleCreateTicket = () => {
    openDrawer(DRAWER_ROUTES.CREATE_TICKET, { projectId, boardColumnId: column.id });
  };

  const [ticketsData] = useListTickets({
    boardColumnId: [column.id],
    projectId: [projectId],
    pageSize: 50,
  });

  const tickets = ticketsData?.items ?? [];
  const hasTickets = tickets.length > 0;

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
        const data = source.data as {
          ticketId: string;
          sourceColumnId?: string;
        };
        if (data.ticketId && data.sourceColumnId !== column.id) {
          onMoveTicket({
            ticketId: data.ticketId,
            move: { boardColumnId: column.id },
          });
        }
      },
    });
  }, [column.id, onMoveTicket]);

  return (
    <div
      ref={dropRef}
      style={{
        display: "flex",
        flexDirection: "column",
        backgroundColor: isDragOver ? "#e8f0fe" : "#f9fafb",
        height: "100%",
        minWidth: "300px",
        transition: "background-color 0.2s",
      }}
    >
      {/* Column Header */}
      <div
        style={{
          padding: "1rem",
          borderBottom: "1px solid #e5e7eb",
          backgroundColor: "#fff",
          borderRadius: "0.75rem 0.75rem 0 0",
          flexShrink: 0,
        }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <Text size="sm" style={{ fontWeight: "600", color: "#1f2937" }}>
            {column.name}
          </Text>
          <div
            style={{
              padding: "0.25rem 0.75rem",
              backgroundColor: "#f3f4f6",
              borderRadius: "0.375rem",
              fontSize: "0.75rem",
              color: "#6b7280",
              fontWeight: "500",
            }}
          >
            {tickets.length}
          </div>
        </div>
      </div>

      {/* Column Content */}
      <div
        style={{
          flex: 1,
          padding: "0.5rem",
          overflowY: "auto",
          display: "flex",
          flexDirection: "column",
          gap: "0.75rem",
        }}
      >
        {!hasTickets ? (
          <div
            style={{
              padding: "2rem 1rem",
              textAlign: "center",
              color: "#9ca3af",
            }}
          >
            <Text size="xs">No tickets yet</Text>
          </div>
        ) : (
          tickets.map((ticket) => (
            <TicketCard
              key={ticket.id}
              ticket={ticket}
              columnId={column.id}
            />
          ))
        )}
      </div>

      {/* Create Ticket Button */}
      <div
        style={{
          padding: "1rem",
          borderTop: "1px solid #e5e7eb",
          backgroundColor: "#fff",
          borderRadius: "0 0 0.75rem 0.75rem",
          flexShrink: 0,
        }}
      >
        <Button
          variant="outline"
          onClick={handleCreateTicket}
          style={{
            width: "100%",
            fontSize: "0.875rem",
          }}
        >
          Create Ticket
        </Button>
      </div>
    </div>
  );
};
