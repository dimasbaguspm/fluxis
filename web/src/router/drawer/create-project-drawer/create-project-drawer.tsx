import { useCreateProject } from "@/hooks/use-api";
import { cx } from "@/lib/cx";
import { useDrawer } from "@/providers/drawer";
import { vGap4 } from "@versaur/core/utilities";
import { ButtonGroup, Drawer, FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { Banner, Button } from "@versaur/react/primitive";
import { useForm } from "react-hook-form";

interface CreateProjectFormInputs {
  orgId: string;
  name: string;
  key: string;
  visibility: "public" | "private";
  description?: string;
}

export const CreateProjectDrawer = () => {
  const { closeDrawer, params } = useDrawer();
  const [createProject, err, { isPending }] = useCreateProject();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateProjectFormInputs>({
    defaultValues: {
      orgId: params?.orgId ?? "",
      name: "",
      key: "",
      visibility: "private",
      description: "",
    },
  });

  const onSubmit = async (data: CreateProjectFormInputs) => {
    await createProject({
      orgId: data.orgId,
      data: {
        name: data.name,
        key: data.key,
        visibility: data.visibility,
        description: data.description,
      },
    });
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Create Project</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <FormGroup onSubmit={handleSubmit(onSubmit)} className={cx(vGap4)}>
          {err && (
            <FormGroup.Field>
              <Banner variant="warning">{(err as any)?.message || "Failed to create project"}</Banner>
            </FormGroup.Field>
          )}
          <FormGroup.Field>
            <TextInput
              placeholder="Organisation ID"
              label="Organisation ID"
              required
              error={errors.orgId?.message}
              {...register("orgId", {
                required: "Organisation ID is required",
              })}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Project Name"
              label="Name"
              required
              error={errors.name?.message}
              {...register("name", {
                required: "Project name is required",
              })}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Project Key"
              label="Key"
              required
              error={errors.key?.message}
              {...register("key", {
                required: "Project key is required",
                maxLength: {
                  value: 10,
                  message: "Key must not exceed 10 characters",
                },
              })}
            />
          </FormGroup.Field>
          <FormGroup.Field>
            <label htmlFor="visibility">Visibility</label>
            <select
              id="visibility"
              {...register("visibility")}
              style={{
                padding: "8px",
                borderRadius: "4px",
                border: "1px solid #e5e5e5",
                fontSize: "14px",
              }}
            >
              <option value="public">Public</option>
              <option value="private">Private</option>
            </select>
          </FormGroup.Field>
          <FormGroup.Field>
            <TextInput
              placeholder="Description (optional)"
              label="Description"
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
          <Button onClick={handleSubmit(onSubmit)} loading={isPending} disabled={isPending}>
            {isPending ? "Creating..." : "Create"}
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
