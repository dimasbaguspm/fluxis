import { useCreateBoard } from "@/hooks/use-api";
import { cx } from "@/lib/cx";
import { useDrawer } from "@/providers/drawer";
import { vGap4 } from "@versaur/core/utilities";
import { ButtonGroup, Drawer, FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { Banner, Button } from "@versaur/react/primitive";
import { useForm } from "react-hook-form";

interface CreateBoardFormInputs {
  sprintId: string;
  name: string;
}

export const CreateBoardDrawer = () => {
  const { closeDrawer, params } = useDrawer();
  const [createBoard, err, { isPending }] = useCreateBoard();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateBoardFormInputs>({
    defaultValues: {
      sprintId: params?.sprintId ?? "",
      name: "",
    },
  });

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
        <FormGroup onSubmit={handleSubmit(onSubmit)} className={cx(vGap4)}>
          {err && (
            <FormGroup.Field>
              <Banner variant="warning">{(err as any)?.message || "Failed to create board"}</Banner>
            </FormGroup.Field>
          )}
          <FormGroup.Field>
            <TextInput
              placeholder="Sprint ID"
              label="Sprint ID"
              required
              error={errors.sprintId?.message}
              {...register("sprintId", {
                required: "Sprint ID is required",
              })}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Board Name"
              label="Name"
              required
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
          <Button onClick={handleSubmit(onSubmit)} loading={isPending} disabled={isPending}>
            {isPending ? "Creating..." : "Create"}
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
