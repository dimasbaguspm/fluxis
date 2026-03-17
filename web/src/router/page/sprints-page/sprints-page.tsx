import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { dateFormat, FormatDate } from "@/lib";
import { useListSprints } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { MenuIcon, PlusIcon } from "@versaur/icons";
import { PageHeader, Table } from "@versaur/react/blocks";
import { ButtonIcon, Text } from "@versaur/react/primitive";

export const SprintsPage = () => {
  const { openDrawer } = useDrawer();
  const [sprints, error, { isLoading }] = useListSprints();

  const handleOnAddButtonClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_SPRINT);
  };

  const handleOnEditSprint = (sprintId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_SPRINT, { sprintId: sprintId });
  };

  if (isLoading) {
    return <div>Loading sprints...</div>;
  }

  if (error) {
    return <div>Error loading sprints</div>;
  }

  const sprintsList = sprints?.items || [];

  return (
    <>
      <PageHeader
        title={
          <PageHeader.Title
            action={
              <ButtonIcon
                aria-label="Add sprint"
                as={PlusIcon}
                variant="ghost"
                onClick={handleOnAddButtonClick}
              />
            }
          >
            Sprints
          </PageHeader.Title>
        }
        subtitle={<PageHeader.Subtitle>Manage your sprints</PageHeader.Subtitle>}
      />
      <Table columns="40px 2fr 1fr 1fr 1fr 100px">
        <Table.Toolbar
          leftContent={(selectedIds) => (
            <Text size="xs">
              {selectedIds.size > 0 ? `${selectedIds.size} selected` : "No selection"}
            </Text>
          )}
          rightContent={(selectedIds) => (
            <Text size="xs">
              {selectedIds.size > 0 ? `${selectedIds.size} selected` : "No selection"}
            </Text>
          )}
        />
        <Table.Header>
          <Table.Col as="th" variant="checkbox">
            <Table.Checkbox isMain />
          </Table.Col>
          <Table.Col as="th">Name</Table.Col>
          <Table.Col as="th">Status</Table.Col>
          <Table.Col as="th">Planned Start</Table.Col>
          <Table.Col as="th">Planned End</Table.Col>
          <Table.Col as="th">Actions</Table.Col>
        </Table.Header>
        <Table.Body>
          {sprintsList.map((sprint) => (
            <Table.Row key={sprint.id}>
              <Table.Col as="td" variant="checkbox">
                <Table.Checkbox rowId={sprint.id} />
              </Table.Col>
              <Table.Col as="td">{sprint.name}</Table.Col>
              <Table.Col as="td">{sprint.status}</Table.Col>
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
                </Table.Action>
              </Table.Col>
            </Table.Row>
          ))}
        </Table.Body>
      </Table>
    </>
  );
};
