import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { dateFormat, FormatDate } from "@/lib";
import { useListProjects } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { MenuIcon, PlusIcon } from "@versaur/icons";
import { PageContent, PageHeader, Table } from "@versaur/react/blocks";
import { ButtonIcon, Text } from "@versaur/react/primitive";

export const ProjectsPage = () => {
  const { openDrawer } = useDrawer();
  const [projects, error, { isLoading }] = useListProjects();

  const handleOnAddButtonClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_PROJECT);
  };

  const handleOnEditProject = (projectId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_PROJECT, { projectId: projectId });
  };

  if (isLoading) {
    return <div>Loading projects...</div>;
  }

  if (error) {
    return <div>Error loading projects</div>;
  }

  const projectsList = projects?.items || [];

  return (
    <>
      <PageHeader
        title={
          <PageHeader.Title
            action={
              <ButtonIcon
                aria-label="Add project"
                as={PlusIcon}
                variant="ghost"
                onClick={handleOnAddButtonClick}
              />
            }
          >
            Projects
          </PageHeader.Title>
        }
        subtitle={<PageHeader.Subtitle>Manage your projects</PageHeader.Subtitle>}
      />
      <PageContent>
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
            <Table.Col as="th">Key</Table.Col>
            <Table.Col as="th">Visibility</Table.Col>
            <Table.Col as="th">Updated</Table.Col>
            <Table.Col as="th">Actions</Table.Col>
          </Table.Header>
          <Table.Body>
            {projectsList.map((project) => (
              <Table.Row key={project.id}>
                <Table.Col as="td" variant="checkbox">
                  <Table.Checkbox rowId={project.id} />
                </Table.Col>
                <Table.Col as="td">{project.name}</Table.Col>
                <Table.Col as="td">{project.key}</Table.Col>
                <Table.Col as="td">{project.visibility}</Table.Col>
                <Table.Col as="td">{dateFormat(project.updatedAt, FormatDate.ShortDate)}</Table.Col>
                <Table.Col as="td" variant="action">
                  <Table.Action icon={MenuIcon}>
                    <Table.ActionItem onClick={() => handleOnEditProject(project.id)}>
                      Edit
                    </Table.ActionItem>
                  </Table.Action>
                </Table.Col>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </PageContent>
    </>
  );
};
