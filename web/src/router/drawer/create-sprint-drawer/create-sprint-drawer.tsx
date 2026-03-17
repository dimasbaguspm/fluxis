import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useCreateSprint } from "@/hooks/use-api";
import { cx } from "@/lib/cx";
import { useDrawer } from "@/providers/drawer";
import { vGap4 } from "@versaur/core/utilities";
import { ButtonGroup, Drawer, FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { Banner, Button } from "@versaur/react/primitive";
import dayjs from "dayjs";
import { useForm } from "react-hook-form";

interface CreateSprintFormInputs {
  projectId: string;
  name: string;
  goal?: string;
  plannedStartedAt?: string;
  plannedCompletedAt?: string;
}

export const CreateSprintDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.CREATE_SPRINT>();
  const [createSprint, err, { isPending }] = useCreateSprint();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateSprintFormInputs>({
    defaultValues: {
      projectId: params?.projectId ?? "",
      name: "",
      goal: "",
      plannedStartedAt: "",
      plannedCompletedAt: "",
    },
  });

  const onSubmit = async (data: CreateSprintFormInputs) => {
    await createSprint({
      projectId: data.projectId,
      name: data.name,
      goal: data.goal,
      plannedStartedAt: data.plannedStartedAt
        ? dayjs(data.plannedStartedAt).toISOString()
        : undefined,
      plannedCompletedAt: data.plannedCompletedAt
        ? dayjs(data.plannedCompletedAt).toISOString()
        : undefined,
    });
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Create Sprint</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <FormGroup onSubmit={handleSubmit(onSubmit)} className={cx(vGap4)}>
          {err && (
            <FormGroup.Field>
              <Banner variant="warning">
                {(err as any)?.message || "Failed to create sprint"}
              </Banner>
            </FormGroup.Field>
          )}
          <FormGroup.Field>
            <TextInput
              placeholder="Project ID"
              label="Project ID"
              required
              error={errors.projectId?.message}
              {...register("projectId", {
                required: "Project ID is required",
              })}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Sprint Name"
              label="Name"
              required
              error={errors.name?.message}
              {...register("name", {
                required: "Sprint name is required",
              })}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Sprint Goal (optional)"
              label="Goal"
              error={errors.goal?.message}
              {...register("goal")}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Planned Start (optional)"
              label="Planned Start"
              type="date"
              error={errors.plannedStartedAt?.message}
              {...register("plannedStartedAt")}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Planned End (optional)"
              label="Planned End"
              type="date"
              error={errors.plannedCompletedAt?.message}
              {...register("plannedCompletedAt")}
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
