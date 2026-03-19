import { useMoveBoardColumn } from "@/hooks/use-api";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";
import { BoardColumn } from "./components/board-column";

export const KanbanBoard = () => {
  const { activeBoard, projectId } = useOutletContext<ProjectDetailContextType>();
  const [moveBoardColumn] = useMoveBoardColumn();

  const columnsList = activeBoard?.columns ?? [];

  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: `repeat(${columnsList.length}, minmax(300px, 1fr))`,
        overflowX: "auto",
        overflowY: "hidden",
        width: "100%",
        minHeight: "100vh",
      }}
    >
      {columnsList.map((column) => (
        <BoardColumn
          key={column.id}
          column={column}
          projectId={projectId}
          onMoveTicket={moveBoardColumn}
        />
      ))}
    </div>
  );
};
