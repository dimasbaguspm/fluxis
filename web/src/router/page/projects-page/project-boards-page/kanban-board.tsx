import { Button, Text, Badge, Heading } from "@versaur/react/primitive";
import { Card, BadgeGroup } from "@versaur/react/blocks";
import type { DomainBoardColumnModel, DomainTicketModel } from "@/interfaces/openapi.generated";
import { useOutletContext } from "react-router";
import { useDrawer } from "@/providers/drawer";
import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import type { ProjectDetailContextType } from "../project-detail-layout";

// Mock tickets for demonstration
const MOCK_TICKETS: DomainTicketModel[] = [
  {
    id: "550e8400-e29b-41d4-a716-446655440001",
    key: "PROJ-1",
    title: "Implement user authentication",
    type: "story",
    priority: "high",
    description: "Add JWT-based authentication",
    storyPoints: 8,
    projectId: "proj-1",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440002",
    key: "PROJ-2",
    title: "Fix login button styling",
    type: "bug",
    priority: "medium",
    description: "Button not aligned properly on mobile",
    storyPoints: 3,
    projectId: "proj-1",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440003",
    key: "PROJ-3",
    title: "Setup database migrations",
    type: "task",
    priority: "critical",
    description: "Create initial schema",
    storyPoints: 5,
    projectId: "proj-1",
  },
];

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

export const KanbanBoard = () => {
  const { activeBoard, projectId } = useOutletContext<ProjectDetailContextType>();

  const columnsList = activeBoard?.columns ?? [];

  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: `repeat(${columnsList.length}, minmax(300px, 1fr))`,
        gap: "1rem",
        overflowX: "auto",
        overflowY: "hidden",
        width: "100%",
        minHeight: "100vh",
      }}
    >
      {columnsList.map((column) => (
        <BoardColumn key={column.id} column={column} projectId={projectId} />
      ))}
    </div>
  );
};

interface BoardColumnProps {
  column: DomainBoardColumnModel;
  projectId: string;
}

function BoardColumn({ column, projectId }: BoardColumnProps) {
  const { openDrawer } = useDrawer();

  const handleCreateTicket = () => {
    openDrawer(DRAWER_ROUTES.CREATE_TICKET, { projectId });
  };

  const tickets = MOCK_TICKETS;
  const hasTickets = tickets.length > 0;

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        backgroundColor: "#f9fafb",
        borderRadius: "0.75rem",
        border: "1px solid #e5e7eb",
        height: "100%",
        minWidth: "300px",
        boxShadow: "0 1px 2px 0 rgb(0 0 0 / 0.05)",
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
            <Card key={ticket.id} size="sm">
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
}
