import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { dateFormat, FormatDate } from "@/lib";
import { useListBoards } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { MenuIcon, PlusIcon } from "@versaur/icons";
import { PageContent, PageHeader, Table } from "@versaur/react/blocks";
import { ButtonIcon, Text } from "@versaur/react/primitive";

export const BoardsPage = () => {
  const { openDrawer } = useDrawer();
  const [boards, error, { isLoading }] = useListBoards();

  const handleOnAddButtonClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_BOARD);
  };

  const handleOnEditBoard = (boardId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_BOARD, { boardId: boardId });
  };

  if (isLoading) {
    return <div>Loading boards...</div>;
  }

  if (error) {
    return <div>Error loading boards</div>;
  }

  const boardsList = boards?.items || [];

  return (
    <>
      <PageHeader
        title={
          <PageHeader.Title
            action={
              <ButtonIcon
                aria-label="Add board"
                as={PlusIcon}
                variant="ghost"
                onClick={handleOnAddButtonClick}
              />
            }
          >
            Boards
          </PageHeader.Title>
        }
        subtitle={<PageHeader.Subtitle>Manage your boards</PageHeader.Subtitle>}
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
            <Table.Col as="th">Updated</Table.Col>
            <Table.Col as="th">Actions</Table.Col>
          </Table.Header>
          <Table.Body>
            {boardsList.map((board) => (
              <Table.Row key={board.id}>
                <Table.Col as="td" variant="checkbox">
                  <Table.Checkbox rowId={board.id} />
                </Table.Col>
                <Table.Col as="td">{board.name}</Table.Col>
                <Table.Col as="td">{dateFormat(board.updatedAt, FormatDate.ShortDate)}</Table.Col>
                <Table.Col as="td" variant="action">
                  <Table.Action icon={MenuIcon}>
                    <Table.ActionItem onClick={() => handleOnEditBoard(board.id)}>
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
