import { DEEP_LINKS } from "@/constants/page-routes";
import { useGetProject } from "@/hooks/use-api";
import type { DomainProjectModel } from "@interfaces/openapi.generated";
import { SettingsIcon, VerticalMenuIcon } from "@versaur/icons";
import { ButtonGroup, Menu, PageHeader, Tabs } from "@versaur/react/blocks";
import { ButtonIcon } from "@versaur/react/primitive";
import { useDrawer } from "@/providers/drawer";
import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useLocation, useNavigate, useParams, Outlet } from "react-router";
import { vWFull } from "@versaur/core/utilities";

export interface ProjectDetailContextType {
  projectId: string;
  project: DomainProjectModel | null;
}

export const ProjectDetailLayout = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const { openDrawer } = useDrawer();
  const [project] = useGetProject(projectId || "");

  if (!projectId) {
    return <div>Project not found</div>;
  }

  const contextValue: ProjectDetailContextType = {
    projectId,
    project: project || null,
  };

  const activeTab = location.pathname.endsWith("/sprints")
    ? "sprints"
    : location.pathname.endsWith("/boards")
      ? "boards"
      : location.pathname.endsWith("/tickets")
        ? "tickets"
        : "overview";

  const handleTabChange = (tab: string) => {
    switch (tab) {
      case "overview":
        navigate(DEEP_LINKS.PROJECT_OVERVIEW(projectId));
        break;
      case "sprints":
        navigate(DEEP_LINKS.PROJECT_SPRINTS(projectId));
        break;
      case "boards":
        navigate(DEEP_LINKS.PROJECT_BOARDS(projectId));
        break;
      case "tickets":
        navigate(DEEP_LINKS.PROJECT_TICKETS(projectId));
        break;
    }
  };

  const handleEditProject = () => {
    openDrawer(DRAWER_ROUTES.UPDATE_PROJECT, { projectId });
  };

  const handleOnAddSprintClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_SPRINT, { projectId });
  };

  const handleOnAddBoardClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_BOARD, { projectId });
  };

  const handleOnAddTicketClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_TICKET, { projectId });
  };

  return (
    <>
      <PageHeader
        title={
          <PageHeader.Title
            action={
              <ButtonGroup>
                <ButtonIcon
                  aria-label="Edit project"
                  as={SettingsIcon}
                  variant="ghost"
                  onClick={handleEditProject}
                />
                <ButtonIcon
                  aria-label="More"
                  as={VerticalMenuIcon}
                  variant="ghost"
                  {...Menu.getTriggerProps({ id: "project-actions-menu" })}
                />

                <Menu id="project-actions-menu" placement="bottom">
                  <Menu.Item onClick={handleOnAddTicketClick}>Add Ticket</Menu.Item>
                  <Menu.Item onClick={handleOnAddBoardClick}>Add Board</Menu.Item>
                  <Menu.Item onClick={handleOnAddSprintClick}>Add Sprint</Menu.Item>
                </Menu>
              </ButtonGroup>
            }
          >
            {project?.name ?? "..."}
          </PageHeader.Title>
        }
        subtitle={<PageHeader.Subtitle>{project?.key}</PageHeader.Subtitle>}
        supplementary={
          <Tabs value={activeTab} onChange={handleTabChange} className={vWFull}>
            <Tabs.Item value="overview">Overview</Tabs.Item>
            <Tabs.Item value="sprints">Sprints</Tabs.Item>
            <Tabs.Item value="boards">Boards</Tabs.Item>
            <Tabs.Item value="tickets">Tickets</Tabs.Item>
          </Tabs>
        }
      />

      <Outlet context={contextValue} />
    </>
  );
};
