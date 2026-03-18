import { Text } from "@versaur/react/primitive";
import { PageContent } from "@versaur/react/blocks";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";
import { KanbanBoard } from "./kanban-board";

export const ProjectBoardPage = () => {
  const { activeSprint, activeBoard } = useOutletContext<ProjectDetailContextType>();

  if (!activeSprint) {
    return (
      <PageContent>
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
            No active sprint found. Create and start a sprint to view its board.
          </Text>
        </div>
      </PageContent>
    );
  }

  if (!activeBoard) {
    return (
      <PageContent>
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
            No board found for the active sprint.
          </Text>
        </div>
      </PageContent>
    );
  }

  return <KanbanBoard boardId={activeBoard.id} />;
};
