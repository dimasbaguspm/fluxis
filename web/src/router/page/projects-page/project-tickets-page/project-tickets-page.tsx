import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useListSprints, useMoveToSprint } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { PageContent } from "@versaur/react/blocks";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";
import { BacklogGroup } from "./components/backlog-group";
import { SprintGroup } from "./components/sprint-group";

export const ProjectTicketsPage = () => {
  const { projectId } = useOutletContext<ProjectDetailContextType>();
  const { openDrawer } = useDrawer();
  const [sprints, sprintsError, { isLoading: isLoadingSprints }] = useListSprints({
    projectId: [projectId],
    pageSize: 50,
  });
  const [moveToSprint] = useMoveToSprint();

  const handleOnEditTicket = (ticketId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_TICKET, { ticketId });
  };

  const handleMoveTicket = (ticketId: string, targetSprintId: string | null) => {
    if (targetSprintId) {
      moveToSprint({ ticketId, move: { sprintId: targetSprintId } });
    } else {
      moveToSprint({ ticketId, move: { sprintId: undefined } });
    }
  };

  const handleCreateTicket = (sprintId: string | undefined) => {
    openDrawer(DRAWER_ROUTES.CREATE_TICKET, { projectId, sprintId });
  };

  const handleStartSprint = (sprintId: string) => {
    // TODO: Implement start sprint functionality
    console.log("Start sprint:", sprintId);
  };

  if (isLoadingSprints) {
    return <PageContent>Loading sprints...</PageContent>;
  }

  if (sprintsError) {
    return <PageContent>Error loading sprints</PageContent>;
  }

  const sprintsList = sprints?.items || [];

  return (
    <PageContent>
      {sprintsList.map((sprint) => (
        <SprintGroup
          key={sprint.id}
          sprint={sprint}
          onEditTicket={handleOnEditTicket}
          onMoveTicket={handleMoveTicket}
          onCreateTicket={handleCreateTicket}
          onStartSprint={handleStartSprint}
        />
      ))}
      <BacklogGroup
        projectId={projectId}
        onEditTicket={handleOnEditTicket}
        onMoveTicket={handleMoveTicket}
        onCreateTicket={handleCreateTicket}
      />
    </PageContent>
  );
};
