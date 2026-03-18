import type { DomainSprintModel } from "@/interfaces/openapi.generated";
import { dateFormat, FormatDate } from "@/lib";
import { ButtonGroup } from "@versaur/react/blocks";
import { Button } from "@versaur/react/primitive";

interface TicketGroupHeaderAction {
  label: string;
  onClick: () => void;
  variant?: "primary" | "secondary";
}

interface TicketGroupHeaderProps {
  title: string;
  description?: string;
  ticketCount: number;
  totalStoryPoints: number;
  sprint?: DomainSprintModel;
  actions?: TicketGroupHeaderAction[];
}

export const TicketGroupHeader = ({
  title,
  description,
  ticketCount,
  totalStoryPoints,
  sprint,
  actions,
}: TicketGroupHeaderProps) => {
  return (
    <div
      style={{
        padding: "1rem",
        backgroundColor: "#f5f5f5",
        borderRadius: "4px 4px 0 0",
        borderBottom: "1px solid #e0e0e0",
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "flex-start",
          marginBottom: "0.5rem",
        }}
      >
        <div>
          <h3 style={{ margin: 0, fontSize: "1rem", fontWeight: 600 }}>{title}</h3>
        </div>
        {actions && actions.length > 0 && (
          <ButtonGroup>
            {actions.map((action, idx) => (
              <Button key={idx} onClick={action.onClick} variant="outline" size="small">
                {action.label}
              </Button>
            ))}
          </ButtonGroup>
        )}
      </div>

      {(description || sprint) && (
        <div style={{ marginBottom: "0.5rem" }}>
          {description && (
            <div style={{ fontSize: "0.875rem", color: "#666", marginTop: "0.25rem" }}>
              {description}
            </div>
          )}
          {sprint && (
            <div style={{ fontSize: "0.875rem", color: "#666", marginTop: "0.25rem" }}>
              Status: <span style={{ fontWeight: 500 }}>{sprint.status}</span>
            </div>
          )}
        </div>
      )}

      {sprint && (sprint.plannedStartedAt || sprint.plannedCompletedAt || sprint.goal) && (
        <div
          style={{ fontSize: "0.875rem", color: "#666", marginTop: "0.75rem", lineHeight: "1.6" }}
        >
          {sprint.plannedStartedAt && (
            <div>
              Planned start:{" "}
              <span style={{ fontWeight: 500 }}>
                {dateFormat(sprint.plannedStartedAt, FormatDate.ShortDate)}
              </span>
            </div>
          )}
          {sprint.plannedCompletedAt && (
            <div>
              Planned end:{" "}
              <span style={{ fontWeight: 500 }}>
                {dateFormat(sprint.plannedCompletedAt, FormatDate.ShortDate)}
              </span>
            </div>
          )}
          {sprint.goal && (
            <div>
              Goal: <span style={{ fontWeight: 500 }}>{sprint.goal}</span>
            </div>
          )}
        </div>
      )}

      <div
        style={{
          fontSize: "0.875rem",
          color: "#666",
          marginTop: "0.75rem",
          paddingTop: "0.75rem",
          borderTop: "1px solid #ddd",
          display: "flex",
          gap: "1.5rem",
        }}
      >
        <div>
          {ticketCount} ticket{ticketCount !== 1 ? "s" : ""}
        </div>
        <div>
          {totalStoryPoints} story point{totalStoryPoints !== 1 ? "s" : ""}
        </div>
      </div>
    </div>
  );
};
