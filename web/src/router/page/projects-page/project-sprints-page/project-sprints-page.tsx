import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { dateFormat, FormatDate } from "@/lib";
import { useListSprints, useStartSprint, useCompleteSprint } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { MenuIcon, SearchXIcon } from "@versaur/icons";
import { NoResults, PageContent, Table } from "@versaur/react/blocks";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";
import { When } from "@/lib/when";
import { Badge, Banner, Loader } from "@versaur/react/primitive";

export const ProjectSprintsPage = () => {
  const { projectId } = useOutletContext<ProjectDetailContextType>();
  const { openDrawer } = useDrawer();
  const [sprints, error, { isLoading }] = useListSprints({ projectId: [projectId] });
  const [startSprint] = useStartSprint();
  const [completeSprint] = useCompleteSprint();

  const handleOnEditSprint = (sprintId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_SPRINT, { sprintId });
  };

  const handleStartSprint = async (sprintId: string) => {
    await startSprint(sprintId);
  };

  const handleCompleteSprint = async (sprintId: string) => {
    await completeSprint(sprintId);
  };

  if (isLoading) {
    return <PageContent>Loading sprints...</PageContent>;
  }

  if (error) {
    return <PageContent>Error loading sprints</PageContent>;
  }

  const sprintsList = sprints?.items || [];

  return (
    <PageContent>
      <When condition={isLoading}>
        <Loader type="bar" />
      </When>
      <When condition={!isLoading}>
        <When condition={!!error}>
          <Banner variant="warning">Failed to load sprints</Banner>
        </When>
        <When condition={sprintsList.length === 0}>
          <NoResults
            icon={SearchXIcon}
            title="No sprints found"
            subtitle="There are no sprints for this project yet."
          />
        </When>
        <When condition={sprintsList.length > 0}>
          <Table columns="2fr 1fr 1fr 1fr 100px">
            <Table.Header>
              <Table.Col as="th">Name</Table.Col>
              <Table.Col as="th">Status</Table.Col>
              <Table.Col as="th">Planned Start</Table.Col>
              <Table.Col as="th">Planned End</Table.Col>
              <Table.Col as="th">Actions</Table.Col>
            </Table.Header>
            <Table.Body>
              {sprintsList.map((sprint) => (
                <Table.Row key={sprint.id}>
                  <Table.Col as="td">{sprint.name}</Table.Col>
                  <Table.Col as="td">
                    <Badge
                      size="small"
                      variant={
                        sprint.status === "active"
                          ? "primary"
                          : sprint.status === "planned"
                            ? "info"
                            : "success"
                      }
                    >
                      {sprint.status.toUpperCase()}
                    </Badge>
                  </Table.Col>
                  <Table.Col as="td">
                    {sprint.plannedStartedAt
                      ? dateFormat(sprint.plannedStartedAt, FormatDate.ShortDate)
                      : "-"}
                  </Table.Col>
                  <Table.Col as="td">
                    {sprint.plannedCompletedAt
                      ? dateFormat(sprint.plannedCompletedAt, FormatDate.ShortDate)
                      : "-"}
                  </Table.Col>
                  <Table.Col as="td" variant="action">
                    <Table.Action icon={MenuIcon}>
                      <Table.ActionItem onClick={() => handleOnEditSprint(sprint.id)}>
                        Edit
                      </Table.ActionItem>
                      {sprint.status === "planned" && (
                        <Table.ActionItem onClick={() => handleStartSprint(sprint.id)}>
                          Start
                        </Table.ActionItem>
                      )}
                      {sprint.status === "active" && (
                        <Table.ActionItem onClick={() => handleCompleteSprint(sprint.id)}>
                          Complete
                        </Table.ActionItem>
                      )}
                    </Table.Action>
                  </Table.Col>
                </Table.Row>
              ))}
            </Table.Body>
          </Table>
        </When>
      </When>
    </PageContent>
  );
};
