import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useDrawer } from "@/providers/drawer";
import { MenuIcon, PlusIcon } from "@versaur/icons";
import { PageHeader, Table } from "@versaur/react/blocks";
import { ButtonIcon, Text } from "@versaur/react/primitive";

export const OrganizationsPage = () => {
  const { openDrawer } = useDrawer();

  const handleOnAddButtonClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_ORGANISATION);
  };
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
          <Table.Col as="th">Product</Table.Col>
          <Table.Col as="th">Category</Table.Col>
          <Table.Col as="th">Price</Table.Col>
          <Table.Col as="th">Stock</Table.Col>
          <Table.Col as="th">Actions</Table.Col>
        </Table.Header>
        <Table.Body>
          {["1", "2", "3"].map((rowId) => (
            <Table.Row key={rowId}>
              <Table.Col as="td" variant="checkbox">
                <Table.Checkbox rowId={rowId} />
              </Table.Col>
              <Table.DoubleLine
                as="td"
                title={`Premium Widget ${rowId}`}
                subtitle={`SKU: PWG-2024-00${rowId}`}
                size="md"
              />
              <Table.Col as="td">Electronics</Table.Col>
              <Table.Col as="td" variant="numeric">
                $199.99
              </Table.Col>
              <Table.Col as="td" variant="numeric">
                {12 + Number(rowId)}
              </Table.Col>
              <Table.Col as="td" variant="action">
                <Table.Action icon={MenuIcon}>
                  <Table.ActionItem
                    onClick={() => {
                      // eslint-disable-next-line no-console
                      console.log(`Edit ${rowId}`);
                    }}
                  >
                    Edit
                  </Table.ActionItem>
                  <Table.ActionItem
                    onClick={() => {
                      // eslint-disable-next-line no-console
                      console.log(`Delete ${rowId}`);
                    }}
                  >
                    Delete
                  </Table.ActionItem>
                </Table.Action>
              </Table.Col>
            </Table.Row>
          ))}
        </Table.Body>
        <Table.Footer>
          <Table.Col as="td" area="span 2">
            Total
          </Table.Col>
          <Table.Col as="td">3 items</Table.Col>
          <Table.Col as="td" variant="numeric">
            $599.97
          </Table.Col>
          <Table.Col as="td" variant="numeric">
            39 units
          </Table.Col>
          <Table.Col as="td" />
        </Table.Footer>
      </Table>
    </>
  );
};
