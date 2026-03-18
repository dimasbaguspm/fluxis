import { Button } from "@versaur/react/primitive";
import { NoResults, PageContent } from "@versaur/react/blocks";
import { useNavigate, useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";
import { When } from "@/lib/when";
import { SearchXIcon } from "@versaur/icons";
import { DEEP_LINKS } from "@/constants";
import { useDrawer } from "@/providers/drawer";
import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { KanbanBoard } from "./kanban-board";

export const ProjectBoardPage = () => {
  const { activeSprint, activeBoard, projectId } = useOutletContext<ProjectDetailContextType>();
  const { openDrawer } = useDrawer();
  const navigate = useNavigate();

  const handleOnPlanSprintClick = () => {
    navigate(DEEP_LINKS.PROJECT_SPRINTS(projectId));
  };

  const handleOnCreateActiveBoardClick = (activeSprint: string) => {
    openDrawer(DRAWER_ROUTES.CREATE_BOARD, { sprintId: activeSprint });
  };

  return (
    <PageContent>
      <When condition={!activeSprint}>
        <NoResults
          icon={SearchXIcon}
          title="No active sprint"
          subtitle="Plan and start a sprint to view its board."
          action={
            <Button variant="outline" onClick={handleOnPlanSprintClick}>
              Plan Sprint
            </Button>
          }
        />
      </When>
      <When condition={!!activeSprint}>
        <When condition={!activeBoard}>
          <NoResults
            icon={SearchXIcon}
            title="No board found"
            subtitle="No board found for the active sprint."
            action={
              <Button
                variant="outline"
                onClick={() => handleOnCreateActiveBoardClick(activeSprint?.id || "")}
              >
                Create Board
              </Button>
            }
          />
        </When>
        <When condition={!!activeBoard}>
          <KanbanBoard />
        </When>
      </When>
    </PageContent>
  );
};
