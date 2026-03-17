import { cx } from "@/lib/cx";
import { vGap4 } from "@versaur/core/utilities";
import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { UPDATE_PROJECT_FORM_ID } from "./constants";
import type { UpdateProjectFormInputs } from "./types";
import type { DomainProjectModel } from "@/interfaces/openapi.generated";

interface UpdateProjectFormProps {
  onSubmit: (data: UpdateProjectFormInputs) => void;
  project: DomainProjectModel;
}

export const UpdateProjectForm = ({ onSubmit, project }: UpdateProjectFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm<UpdateProjectFormInputs>({
    mode: "all",
    defaultValues: {
      orgId: project.orgId,
      name: project.name,
      description: project.description ?? "",
    },
  });

  return (
    <FormGroup id={UPDATE_PROJECT_FORM_ID} onSubmit={handleSubmit(onSubmit)} className={cx(vGap4)}>
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
          placeholder="Description"
          label="Description"
          error={errors.description?.message}
          {...register("description")}
        />
      </FormGroup.Field>
    </FormGroup>
  );
};
