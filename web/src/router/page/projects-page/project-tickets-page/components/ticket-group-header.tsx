import type { DomainSprintModel } from "@/interfaces/openapi.generated";
import { When } from "@/lib/when";
import { ButtonGroup } from "@versaur/react/blocks";
import { Badge, Button, Heading, Text } from "@versaur/react/primitive";

interface TicketGroupHeaderProps {
  title: string;
  description?: string;
  totalStoryPoints: number;
  sprint?: DomainSprintModel;
  onCreateTicket?: () => void;
  onStartSprint?: () => void;
  onCompleteSprint?: () => void;
  onEditSprint?: () => void;
}

export const TicketGroupHeader = ({
  title,
  description,
  totalStoryPoints,
  sprint,
  onCreateTicket,
  onStartSprint,
  onCompleteSprint,
  onEditSprint,
}: TicketGroupHeaderProps) => {
  const getSprintStatusVariant = (status?: string) => {
    switch (status) {
      case "active":
        return "primary" as const;
      case "planned":
        return "info" as const;
      case "completed":
        return "success" as const;
      default:
        return "success" as const;
    }
  };

  return (
    <div
      style={{
        borderBottom: "1px solid #e0e0e0",
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          gap: "1rem",
        }}
      >
        {/* Left: Title + Status Badge + Description */}
        <div
          onClick={onEditSprint}
          style={{
            flex: 1,
            minWidth: 0,
            display: "flex",
            alignItems: "center",
            gap: "0.75rem",
            cursor: onEditSprint ? "pointer" : "default",
            borderRadius: "4px",
            padding: "0.5rem 1rem",
            height: "100%",
            transition: "background-color 0.2s",
          }}
          onMouseEnter={(e) => {
            if (onEditSprint) {
              e.currentTarget.style.backgroundColor = "rgba(0, 0, 0, 0.04)";
            }
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = "transparent";
          }}
        >
          <Heading as="h3" size="base">
            {title}
          </Heading>
          <When condition={!!sprint}>
            <Badge size="small" variant={getSprintStatusVariant(sprint?.status)}>
              {sprint?.status.toUpperCase()}
            </Badge>
          </When>
          <When condition={!!description}>
            <Text as="span" intent="gray" size="xs">
              {description}
            </Text>
          </When>
        </div>

        {/* Right: Story Points Badge + Action Buttons */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: "1rem",
            flexShrink: 0,
            padding: "0.5rem 1rem",
          }}
        >
          <Badge variant="info" size="small" shape="pill">
            {totalStoryPoints}
          </Badge>
          <ButtonGroup>
            <When condition={onCreateTicket}>
              <Button onClick={onCreateTicket} variant="outline" size="small">
                Create Ticket
              </Button>
            </When>
            <When condition={[sprint?.status === "planned", onStartSprint]}>
              <Button onClick={onStartSprint} variant="outline" size="small">
                Start Sprint
              </Button>
            </When>
            <When condition={[sprint?.status === "active", onCompleteSprint]}>
              <Button onClick={onCompleteSprint} variant="outline" size="small">
                Complete Sprint
              </Button>
            </When>
          </ButtonGroup>
        </div>
      </div>
    </div>
  );
};
