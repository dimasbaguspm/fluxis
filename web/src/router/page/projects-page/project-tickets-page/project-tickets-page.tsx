import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import {
  useCompleteSprint,
  useListSprints,
  useMoveToSprint,
  useStartSprint,
} from "@/hooks/use-api";
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
  const [startSprint] = useStartSprint();
  const [completeSprint] = useCompleteSprint();

  const handleOnEditTicket = (ticketId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_TICKET, { ticketId });
  };

  const handleEditSprint = (sprintId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_SPRINT, { sprintId });
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
    startSprint(sprintId);
  };

  const handleCompleteSprint = (sprintId: string) => {
    completeSprint(sprintId);
  };

  if (isLoadingSprints) {
    return <PageContent>Loading sprints...</PageContent>;
  }

  if (sprintsError) {
    return <PageContent>Error loading sprints</PageContent>;
  }

  const sprintsList = sprints?.items || [];

  // Sort: active first, then others by createdAt desc, backlog always last
  const activeSprints = sprintsList.filter((s) => s.status === "active");
  const otherSprints = sprintsList
    .filter((s) => s.status !== "active")
    .sort((a, b) => new Date(b.createdAt ?? 0).getTime() - new Date(a.createdAt ?? 0).getTime());
  const orderedSprints = [...activeSprints, ...otherSprints];

  return (
    <PageContent>
      {orderedSprints.map((sprint) => (
        <SprintGroup
          key={sprint.id}
          sprint={sprint}
          onEditTicket={handleOnEditTicket}
          onEditSprint={handleEditSprint}
          onMoveTicket={handleMoveTicket}
          onCreateTicket={handleCreateTicket}
          onStartSprint={handleStartSprint}
          onCompleteSprint={handleCompleteSprint}
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
