import type { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useCreateBoard } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { ButtonGroup, Drawer } from "@versaur/react/blocks";
import { Banner, Button } from "@versaur/react/primitive";
import { CREATE_BOARD_FORM_ID } from "./constants";
import { CreateBoardForm } from "./form";
import type { CreateBoardFormInputs } from "./types";

export const CreateBoardDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.CREATE_BOARD>();
  const [createBoard, err, { isPending }] = useCreateBoard();

  const onSubmit = async (data: CreateBoardFormInputs) => {
    await createBoard({
      sprintId: data.sprintId,
      name: data.name,
    });
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Create Board</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        {err && (
          <Banner variant="warning" style={{ marginBottom: "1rem" }}>
            {err?.message}
          </Banner>
        )}
        <CreateBoardForm sprintId={params?.sprintId ?? ""} onSubmit={onSubmit} />
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button
            type="submit"
            form={CREATE_BOARD_FORM_ID}
            loading={isPending}
            disabled={isPending}
          >
            Create
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
