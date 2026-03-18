import { Text } from "@versaur/react/primitive";
import type { DomainBoardColumnModel } from "@/interfaces/openapi.generated";
import { useListBoardColumns } from "@/hooks/use-api";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";

export const KanbanBoard = () => {
  const { activeBoard } = useOutletContext<ProjectDetailContextType>();
  const [columns] = useListBoardColumns(activeBoard?.id || "");

  const columnsList = columns?.items || [];

  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: `repeat(${columnsList.length}, minmax(300px, 1fr))`,
        gap: "1rem",
        padding: "1.5rem",
        overflowX: "auto",
        overflowY: "hidden",
        width: "100%",
        height: "100%",
      }}
    >
      {columnsList.map((column) => (
        <BoardColumn key={column.id} column={column} />
      ))}
    </div>
  );
};

interface BoardColumnProps {
  column: DomainBoardColumnModel;
}

function BoardColumn({ column }: BoardColumnProps) {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        backgroundColor: "#f3f4f6",
        borderRadius: "0.5rem",
        border: "1px solid #e5e7eb",
        height: "100%",
        minWidth: "300px",
      }}
    >
      {/* Column Header */}
      <div
        style={{
          padding: "1rem",
          borderBottom: "1px solid #e5e7eb",
          backgroundColor: "#fff",
          borderRadius: "0.5rem 0.5rem 0 0",
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
              padding: "0.25rem 0.5rem",
              backgroundColor: "#e5e7eb",
              borderRadius: "0.25rem",
              fontSize: "0.75rem",
              color: "#4b5563",
              fontWeight: "500",
            }}
          >
            0
          </div>
        </div>
      </div>

      {/* Column Content */}
      <div
        style={{
          flex: 1,
          padding: "1rem",
          overflowY: "auto",
          display: "flex",
          flexDirection: "column",
          gap: "0.75rem",
        }}
      >
        {/* Tickets will go here */}
        <div
          style={{
            padding: "2rem 1rem",
            textAlign: "center",
            color: "#9ca3af",
          }}
        >
          <Text size="xs">No tickets yet</Text>
        </div>
      </div>
    </div>
  );
}
