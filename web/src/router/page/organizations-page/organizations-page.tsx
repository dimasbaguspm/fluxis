import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { dateFormat, FormatDate } from "@/lib";
import { useListOrgs } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { MenuIcon, PlusIcon } from "@versaur/icons";
import { PageContent, PageHeader, Table } from "@versaur/react/blocks";
import { ButtonIcon, Text } from "@versaur/react/primitive";

export const OrganizationsPage = () => {
  const { openDrawer } = useDrawer();
  const [orgs, error, { isLoading }] = useListOrgs();

  const handleOnAddButtonClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_ORGANISATION);
  };

  const handleOnEditOrganisation = (orgId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_ORGANISATION, { orgId: orgId });
  };

  if (isLoading) {
    return <div>Loading organizations...</div>;
  }

  if (error) {
    return <div>Error loading organizations</div>;
  }

  const organizations = orgs?.items || [];

  return (
    <>
      <PageHeader
        title={
          <PageHeader.Title
            action={
              <ButtonIcon
                aria-label="Search"
                as={PlusIcon}
                variant="ghost"
                onClick={handleOnAddButtonClick}
              />
            }
          >
            Organisation
          </PageHeader.Title>
        }
        subtitle={<PageHeader.Subtitle>Manage your projects and settings</PageHeader.Subtitle>}
      />
      <PageContent>
        <Table columns="40px 2fr 1fr 1fr 100px">
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
            <Table.Col as="th">Members</Table.Col>
            <Table.Col as="th">Updated</Table.Col>
            <Table.Col as="th">Actions</Table.Col>
          </Table.Header>
          <Table.Body>
            {organizations.map((org) => (
              <Table.Row key={org.id}>
                <Table.Col as="td" variant="checkbox">
                  <Table.Checkbox rowId={org.id} />
                </Table.Col>
                <Table.Col as="td">{org.name}</Table.Col>
                <Table.Col as="td" variant="numeric">
                  {org.totalMembers || 0}
                </Table.Col>
                <Table.Col as="td">{dateFormat(org.updatedAt, FormatDate.ShortDate)}</Table.Col>
                <Table.Col as="td" variant="action">
                  <Table.Action icon={MenuIcon}>
                    <Table.ActionItem onClick={() => handleOnEditOrganisation(org.id)}>
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
