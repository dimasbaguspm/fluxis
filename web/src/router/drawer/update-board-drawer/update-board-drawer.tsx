import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetBoard, useUpdateBoard } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { ButtonGroup, Drawer, NoResults } from "@versaur/react/blocks";
import { Banner, Button, Loader } from "@versaur/react/primitive";
import { When } from "@/lib/when";
import { SearchXIcon } from "@versaur/icons";
import { UPDATE_BOARD_FORM_ID } from "./constants";
import { UpdateBoardForm } from "./form";
import type { UpdateBoardFormInputs } from "./types";

export const UpdateBoardDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_BOARD>();
  const boardId = params?.boardId ?? "";

  const [board, , { isLoading: isLoadingBoard }] = useGetBoard(boardId);
  const [updateBoard, err, { isPending }] = useUpdateBoard();

  const onSubmit = async (data: UpdateBoardFormInputs) => {
    if (!boardId) return;
    await updateBoard({
      boardId: boardId,
      data: {
        sprintId: data.sprintId,
        name: data.name,
      },
    });
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Update Board</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <When condition={isLoadingBoard}>
          <Loader type="bar" />
        </When>
        <When condition={!isLoadingBoard}>
          <When condition={!board}>
            <NoResults
              icon={SearchXIcon}
              title="Board not found"
              subtitle="We couldn't find the board you're looking for"
            />
          </When>
          <When condition={!!board}>
            <When condition={err?.message}>
              <Banner variant="warning">{err?.message}</Banner>
            </When>
            <UpdateBoardForm board={board!} onSubmit={onSubmit} />
          </When>
        </When>
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button
            type="submit"
            form={UPDATE_BOARD_FORM_ID}
            loading={isPending || isLoadingBoard}
            disabled={isPending || isLoadingBoard}
          >
            Update
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
