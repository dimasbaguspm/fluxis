import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetProject, useUpdateProject } from "@/hooks/use-api";
import { cx } from "@/lib/cx";
import { useDrawer } from "@/providers/drawer";
import { SelectOrganisationsInput } from "@/components/ui/select-organisations-input";
import { vGap4 } from "@versaur/core/utilities";
import { ButtonGroup, Drawer, FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { Banner, Button } from "@versaur/react/primitive";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

interface UpdateProjectFormInputs {
  orgId: string;
  name: string;
  description?: string;
}

export const UpdateProjectDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_PROJECT>();
  const projectId = params?.projectId ?? "";

  const [project, , { isLoading: isLoadingProject }] = useGetProject(projectId);
  const [updateProject, err, { isPending }] = useUpdateProject();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
    control,
  } = useForm<UpdateProjectFormInputs>({
    mode: "onBlur",
    defaultValues: { orgId: "", name: "", description: "" },
  });

  useEffect(() => {
    if (project?.name) {
      setValue("orgId", project.orgId);
      setValue("name", project.name);
      setValue("description", project.description ?? "");
    }
  }, [project?.name, project?.description, project?.orgId, setValue]);

  const onSubmit = async (data: UpdateProjectFormInputs) => {
    if (!projectId) return;
    await updateProject({
      projectId: projectId,
      data: {
        name: data.name,
        description: data.description,
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
    return "Failed to update project";
  })();

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Update Project</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <FormGroup onSubmit={handleSubmit(onSubmit)} className={cx(vGap4)}>
          {errorMessage ? (
            <FormGroup.Field>
              <Banner variant="warning">{errorMessage}</Banner>
            </FormGroup.Field>
          ) : null}
          <FormGroup.Field>
            <SelectOrganisationsInput
              control={control}
              name="orgId"
              disabled={true}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Project Name"
              label="Name"
              required
              disabled={isLoadingProject}
              error={errors.name?.message}
              {...register("name", {
                required: "Project name is required",
              })}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Description (optional)"
              label="Description"
              disabled={isLoadingProject}
              error={errors.description?.message}
              {...register("description")}
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
            loading={isPending || isLoadingProject}
            disabled={isPending || isLoadingProject}
          >
            {isPending ? "Updating..." : isLoadingProject ? "Loading..." : "Update"}
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
