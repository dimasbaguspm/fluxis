import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { SelectProjectsInput, SelectSprintsInput } from "@/components/ui";
import { CREATE_TICKET_FORM_ID } from "./constants";
import type { CreateTicketFormInputs } from "./types";

interface CreateTicketFormProps {
  onSubmit: (data: CreateTicketFormInputs) => void;
  projectId?: string;
}

export const CreateTicketForm = ({
  onSubmit,
  projectId,
}: CreateTicketFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
    control,
  } = useForm<CreateTicketFormInputs>({
    mode: "all",
    defaultValues: {
      projectId: projectId ?? "",
      title: "",
      type: "",
      priority: "",
      description: "",
      sprintId: "",
      storyPoints: undefined,
      dueDate: "",
    },
  });

  return (
    <FormGroup id={CREATE_TICKET_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <SelectProjectsInput control={control} name="projectId" required />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="Ticket Title"
          label="Title"
          required
          error={errors.title?.message}
          {...register("title", {
            required: "Title is required",
          })}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="bug / story / task / epic"
          label="Type"
          required
          error={errors.type?.message}
          {...register("type", {
            required: "Type is required",
          })}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="low / medium / high / critical"
          label="Priority"
          required
          error={errors.priority?.message}
          {...register("priority", {
            required: "Priority is required",
          })}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="Description"
          label="Description"
          {...register("description")}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <SelectSprintsInput control={control} name="sprintId" />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="Story Points"
          label="Story Points"
          type="number"
          {...register("storyPoints", {
            valueAsNumber: true,
          })}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="Due Date"
          label="Due Date"
          type="date"
          {...register("dueDate")}
        />
      </FormGroup.Field>
    </FormGroup>
  );
};
