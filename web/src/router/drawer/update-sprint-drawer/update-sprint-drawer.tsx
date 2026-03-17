import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetSprint, useUpdateSprint } from "@/hooks/use-api";
import { dateFormat, FormatDate } from "@/lib";
import { cx } from "@/lib/cx";
import { useDrawer } from "@/providers/drawer";
import { vGap4 } from "@versaur/core/utilities";
import { ButtonGroup, Drawer, FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { Banner, Button } from "@versaur/react/primitive";
import dayjs from "dayjs";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

interface UpdateSprintFormInputs {
  name: string;
  goal?: string;
  plannedStartedAt?: string;
  plannedCompletedAt?: string;
}

export const UpdateSprintDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_SPRINT>();
  const sprintId = params?.sprintId ?? "";

  const [sprint, , { isLoading: isLoadingSprint }] = useGetSprint(sprintId);
  const [updateSprint, err, { isPending }] = useUpdateSprint();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
  } = useForm<UpdateSprintFormInputs>({
    defaultValues: {
      name: "",
      goal: "",
      plannedStartedAt: "",
      plannedCompletedAt: "",
    },
  });

  useEffect(() => {
    if (sprint?.name) {
      setValue("name", sprint.name);
      setValue("goal", sprint.goal ?? "");
      setValue("plannedStartedAt", dateFormat(sprint.plannedStartedAt, FormatDate.ShortDate) ?? "");
      setValue(
        "plannedCompletedAt",
        dateFormat(sprint.plannedCompletedAt, FormatDate.ShortDate) ?? "",
      );
    }
  }, [sprint?.name, sprint?.goal, sprint?.plannedStartedAt, sprint?.plannedCompletedAt, setValue]);

  const onSubmit = async (data: UpdateSprintFormInputs) => {
    if (!sprintId) return;
    await updateSprint({
      sprintId: sprintId,
      data: {
        name: data.name,
        goal: data.goal,
        plannedStartedAt: data.plannedStartedAt
          ? dayjs(data.plannedStartedAt).toISOString()
          : undefined,
        plannedCompletedAt: data.plannedCompletedAt
          ? dayjs(data.plannedCompletedAt).toISOString()
          : undefined,
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
    return "Failed to update sprint";
  })();

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Update Sprint</Drawer.Title>
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
              placeholder="Sprint Name"
              label="Name"
              required
              disabled={isLoadingSprint}
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
              disabled={isLoadingSprint}
              error={errors.goal?.message}
              {...register("goal")}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Planned Start (optional)"
              label="Planned Start"
              type="date"
              disabled={isLoadingSprint}
              error={errors.plannedStartedAt?.message}
              {...register("plannedStartedAt")}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Planned End (optional)"
              label="Planned End"
              type="date"
              disabled={isLoadingSprint}
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
          <Button
            onClick={handleSubmit(onSubmit)}
            loading={isPending || isLoadingSprint}
            disabled={isPending || isLoadingSprint}
          >
            {isPending ? "Updating..." : isLoadingSprint ? "Loading..." : "Update"}
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
