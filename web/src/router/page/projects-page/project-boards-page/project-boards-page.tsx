import { Text } from "@versaur/react/primitive";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";
import { KanbanBoard } from "./kanban-board";

export const ProjectBoardPage = () => {
  const { projectId: _projectId } = useOutletContext<ProjectDetailContextType>();

  // TODO: Get active board from context or route params
  // For now, using a hardcoded board ID
  const boardId = "board-1";

  if (!boardId) {
    return (
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          height: "100%",
          padding: "2rem",
        }}
      >
        <Text size="sm" style={{ color: "#6b7280" }}>
          No board selected.
        </Text>
      </div>
    );
  }

  return <KanbanBoard boardId={boardId} />;
}
