import { SelectProjectsInput } from "@/components/ui";
import { When } from "@/lib/when";
import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { CREATE_SPRINT_FORM_ID } from "./constants";
import type { CreateSprintFormInputs } from "./types";

interface CreateSprintFormProps {
  onSubmit: (data: CreateSprintFormInputs) => void;
  projectId?: string;
}

export const CreateSprintForm = ({ onSubmit, projectId }: CreateSprintFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
    control,
  } = useForm<CreateSprintFormInputs>({
    mode: "all",
    defaultValues: {
      projectId: projectId ?? "",
      name: "",
      goal: "",
      plannedStartedAt: "",
      plannedCompletedAt: "",
    },
  });

  return (
    <FormGroup id={CREATE_SPRINT_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <When condition={!projectId}>
        <FormGroup.Field>
          <SelectProjectsInput control={control} name="projectId" required />
        </FormGroup.Field>
      </When>
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
  );
};
