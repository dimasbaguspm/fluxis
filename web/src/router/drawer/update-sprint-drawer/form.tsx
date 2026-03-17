import { dateFormat, FormatDate } from "@/lib";
import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { UPDATE_SPRINT_FORM_ID } from "./constants";
import type { UpdateSprintFormInputs } from "./types";
import type { DomainSprintModel } from "@/interfaces/openapi.generated";

interface UpdateSprintFormProps {
  onSubmit: (data: UpdateSprintFormInputs) => void;
  sprint: DomainSprintModel;
}

export const UpdateSprintForm = ({ onSubmit, sprint }: UpdateSprintFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm<UpdateSprintFormInputs>({
    mode: "all",
    defaultValues: {
      name: sprint.name,
      goal: sprint.goal ?? "",
      plannedStartedAt: dateFormat(sprint.plannedStartedAt, FormatDate.ShortDate) ?? "",
      plannedCompletedAt: dateFormat(sprint.plannedCompletedAt, FormatDate.ShortDate) ?? "",
    },
  });

  return (
    <FormGroup id={UPDATE_SPRINT_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
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
