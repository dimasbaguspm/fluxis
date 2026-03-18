import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useListSprints } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { PageContent } from "@versaur/react/blocks";
import { useOutletContext } from "react-router";
import { SprintGroup } from "./components/sprint-group";
import { BacklogGroup } from "./components/backlog-group";
import type { ProjectDetailContextType } from "../project-detail-layout";

export const ProjectTicketsPage = () => {
  const { projectId } = useOutletContext<ProjectDetailContextType>();
  const { openDrawer } = useDrawer();
  const [sprints, sprintsError, { isLoading: isLoadingSprints }] = useListSprints({
    projectId: [projectId],
    pageSize: 50,
  });

  const handleOnEditTicket = (ticketId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_TICKET, { ticketId });
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
      {sprintsList.length === 0 ? (
        <div style={{ textAlign: "center", padding: "2rem", color: "#999" }}>
          No sprints found for this project
        </div>
      ) : (
        <>
          {sprintsList.map((sprint) => (
            <SprintGroup key={sprint.id} sprint={sprint} onEditTicket={handleOnEditTicket} />
          ))}
          <BacklogGroup projectId={projectId} onEditTicket={handleOnEditTicket} />
        </>
      )}
    </PageContent>
  );
};
