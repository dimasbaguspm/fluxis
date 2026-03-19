import type { DomainTicketModel } from "@/interfaces/openapi.generated";
import { draggable } from "@atlaskit/pragmatic-drag-and-drop/element/adapter";
import { Card } from "@versaur/react/blocks";
import { Badge, Heading, Text } from "@versaur/react/primitive";
import { BadgeGroup } from "@versaur/react/blocks";
import { useEffect, useRef } from "react";

const getPriorityVariant = (priority?: string) => {
  switch (priority) {
    case "critical":
      return "warning";
    case "high":
      return "warning";
    case "medium":
      return "info";
    case "low":
      return "success";
    default:
      return "info";
  }
};

const getTypeVariant = (type?: string) => {
  switch (type) {
    case "bug":
      return "warning";
    case "story":
      return "info";
    case "task":
      return "warning";
    case "epic":
      return "success";
    default:
      return "info";
  }
};

interface TicketCardProps {
  ticket: DomainTicketModel;
  columnId: string;
}

export const TicketCard = ({ ticket, columnId }: TicketCardProps) => {
  const cardRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!cardRef.current) return;

    return draggable({
      element: cardRef.current,
      getInitialData: () => ({
        ticketId: ticket.id,
        sourceColumnId: columnId,
        type: "ticket",
      }),
    });
  }, [ticket.id, columnId]);

  return (
    <Card ref={cardRef} size="sm" style={{ cursor: "grab" }}>
      <Card.Header
        style={{
          gap: "0.5rem",
          flexDirection: "column",
          alignItems: "flex-start",
        }}
      >
        <div style={{ width: "100%" }}>
          <Heading
            as="h3"
            size="sm"
            style={{
              margin: 0,
              fontSize: "0.875rem",
              fontWeight: "600",
              color: "#1f2937",
              lineHeight: "1.4",
            }}
          >
            {ticket.title}
          </Heading>
          <Text
            size="xs"
            style={{
              marginTop: "0.25rem",
              color: "#6b7280",
            }}
          >
            {ticket.key}
          </Text>
        </div>
      </Card.Header>
      <Card.Footer
        style={{
          gap: "0.5rem",
          justifyContent: "space-between",
          flexWrap: "wrap",
        }}
      >
        <BadgeGroup>
          <Badge size="small" variant={getTypeVariant(ticket.type)}>
            {ticket.type}
          </Badge>
          <Badge size="small" variant={getPriorityVariant(ticket.priority)}>
            {ticket.priority}
          </Badge>
        </BadgeGroup>
        {ticket.storyPoints && (
          <Text size="xs" style={{ color: "#6b7280" }}>
            {ticket.storyPoints}pts
          </Text>
        )}
      </Card.Footer>
    </Card>
  );
};
