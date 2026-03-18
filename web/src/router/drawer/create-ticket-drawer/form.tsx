import { SelectProjectsInput, SelectSprintsInput } from "@/components/ui";
import { When } from "@/lib/when";
import { FormGroup } from "@versaur/react/blocks";
import { Select, TextArea, TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { CREATE_TICKET_FORM_ID, TICKET_PRIORITIES, TICKET_TYPES } from "./constants";
import type { CreateTicketFormInputs } from "./types";

interface CreateTicketFormProps {
  onSubmit: (data: CreateTicketFormInputs) => void;
  projectId?: string;
  sprintId?: string;
}

export const CreateTicketForm = ({ onSubmit, projectId, sprintId }: CreateTicketFormProps) => {
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
      sprintId: sprintId ?? "",
      storyPoints: undefined,
      dueDate: "",
    },
  });

  return (
    <FormGroup id={CREATE_TICKET_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <When condition={!projectId}>
        <FormGroup.Field>
          <SelectProjectsInput control={control} name="projectId" required />
        </FormGroup.Field>
      </When>
      <When condition={!sprintId}>
        <FormGroup.Field>
          <SelectSprintsInput control={control} name="sprintId" projectId={projectId} />
        </FormGroup.Field>
      </When>
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
        <Select
          label="Type"
          required
          error={errors.type?.message}
          {...register("type", {
            required: "Type is required",
          })}
        >
          <Select.Option value="">Select a type</Select.Option>
          {TICKET_TYPES.map((type) => (
            <Select.Option key={type} value={type}>
              {type.charAt(0).toUpperCase() + type.slice(1)}
            </Select.Option>
          ))}
        </Select>
      </FormGroup.Field>
      <FormGroup.Field>
        <Select
          label="Priority"
          required
          error={errors.priority?.message}
          {...register("priority", {
            required: "Priority is required",
          })}
        >
          <Select.Option value="">Select a priority</Select.Option>
          {TICKET_PRIORITIES.map((priority) => (
            <Select.Option key={priority} value={priority}>
              {priority.charAt(0).toUpperCase() + priority.slice(1)}
            </Select.Option>
          ))}
        </Select>
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
        <TextInput placeholder="Due Date" label="Due Date" type="date" {...register("dueDate")} />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextArea
          placeholder="Enter ticket description..."
          label="Description"
          minRows={5}
          resizable={false}
          {...register("description")}
        />
      </FormGroup.Field>
    </FormGroup>
  );
};
