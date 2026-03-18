import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { dateFormat, FormatDate } from "@/lib";
import { useListBoards } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { SelectSprintsInput } from "@/components/ui";
import { MenuIcon, PlusIcon } from "@versaur/icons";
import { PageContent, PageHeader, Table } from "@versaur/react/blocks";
import { ButtonIcon, Text } from "@versaur/react/primitive";
import { useForm, useWatch } from "react-hook-form";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";

export const ProjectBoardsPage = () => {
  const { projectId } = useOutletContext<ProjectDetailContextType>();
  const { openDrawer } = useDrawer();

  const { control } = useForm({
    defaultValues: {
      sprintId: "",
    },
  });

  const selectedSprintId = useWatch({ control, name: "sprintId" });

  const [boards, error, { isLoading }] = useListBoards(
    selectedSprintId ? { sprintId: [selectedSprintId] } : undefined,
    { enabled: !!selectedSprintId },
  );

  const handleOnAddButtonClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_BOARD, { sprintId: selectedSprintId });
  };

  const handleOnEditBoard = (boardId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_BOARD, { boardId });
  };

  const boardsList = boards?.items || [];

  return (
    <>
      <PageHeader
        title={
          <PageHeader.Title
            action={
              selectedSprintId && (
                <ButtonIcon
                  aria-label="Add board"
                  as={PlusIcon}
                  variant="ghost"
                  onClick={handleOnAddButtonClick}
                />
              )
            }
          >
            Boards
          </PageHeader.Title>
        }
      />

      <PageContent>
        <div style={{ marginBottom: "2rem", maxWidth: "400px" }}>
          <SelectSprintsInput control={control} name="sprintId" projectId={projectId} />
        </div>

        {!selectedSprintId ? (
          <div style={{ padding: "2rem", textAlign: "center", color: "#666" }}>
            Select a sprint to view boards
          </div>
        ) : (
          <>
            {isLoading && <div>Loading boards...</div>}

            {error && <div>Error loading boards</div>}

            {!isLoading && !error && (
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
                  <Table.Col as="th">Position</Table.Col>
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
                      <Table.Col as="td" variant="numeric">
                        {board.position}
                      </Table.Col>
                      <Table.Col as="td">
                        {dateFormat(board.updatedAt, FormatDate.ShortDate)}
                      </Table.Col>
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
            )}
          </>
        )}
      </PageContent>
    </>
  );
};
