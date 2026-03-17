import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetBoard, useUpdateBoard } from "@/hooks/use-api";
import { cx } from "@/lib/cx";
import { useDrawer } from "@/providers/drawer";
import { vGap4 } from "@versaur/core/utilities";
import { ButtonGroup, Drawer, FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { Banner, Button } from "@versaur/react/primitive";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

interface UpdateBoardFormInputs {
  name: string;
}

export const UpdateBoardDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_BOARD>();
  const boardId = params?.boardId ?? "";

  const [board, , { isLoading: isLoadingBoard }] = useGetBoard(boardId);
  const [updateBoard, err, { isPending }] = useUpdateBoard();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
  } = useForm<UpdateBoardFormInputs>({
    defaultValues: { name: "" },
  });

  useEffect(() => {
    if (board?.name) {
      setValue("name", board.name);
    }
  }, [board?.name, setValue]);

  const onSubmit = async (data: UpdateBoardFormInputs) => {
    if (!boardId) return;
    await updateBoard({
      boardId: boardId,
      data: {
        name: data.name,
      },
    });
    reset();
    closeDrawer();
  };

  const errorMessage = (() => {
    if (!err) return null;
    if (typeof err === "object" && err !== null && "message" in err) {
      return String((err as any).message);
    }
    return "Failed to update board";
  })();

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Update Board</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <FormGroup onSubmit={handleSubmit(onSubmit)} className={cx(vGap4)}>
          {errorMessage ? (
            <FormGroup.Field>
              <Banner variant="warning">{errorMessage}</Banner>
            </FormGroup.Field>
          ) : null}
          <FormGroup.Field>
            <TextInput
              placeholder="Board Name"
              label="Name"
              required
              disabled={isLoadingBoard}
              error={errors.name?.message}
              {...register("name", {
                required: "Board name is required",
              })}
            />
          </FormGroup.Field>
        </FormGroup>
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button
            onClick={handleSubmit(onSubmit)}
            loading={isPending || isLoadingBoard}
            disabled={isPending || isLoadingBoard}
          >
            {isPending ? "Updating..." : isLoadingBoard ? "Loading..." : "Update"}
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
