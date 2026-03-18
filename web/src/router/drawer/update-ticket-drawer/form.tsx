import { SelectSprintsInput } from "@/components/ui";
import type { DomainTicketModel } from "@/interfaces/openapi.generated";
import { dateFormat, FormatDate } from "@/lib";
import { FormGroup } from "@versaur/react/blocks";
import { Select, TextArea, TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { TICKET_PRIORITIES, TICKET_TYPES, UPDATE_TICKET_FORM_ID } from "./constants";
import type { UpdateTicketFormInputs } from "./types";

interface UpdateTicketFormProps {
  onSubmit: (data: UpdateTicketFormInputs) => void;
  ticket: DomainTicketModel;
}

export const UpdateTicketForm = ({ onSubmit, ticket }: UpdateTicketFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
    control,
  } = useForm<UpdateTicketFormInputs>({
    mode: "all",
    defaultValues: {
      title: ticket.title ?? "",
      type: ticket.type ?? "",
      priority: ticket.priority ?? "",
      description: ticket.description ?? "",
      sprintId: ticket.sprintId ?? "",
      storyPoints: ticket.storyPoints,
      dueDate: ticket.dueDate ? dateFormat(ticket.dueDate, FormatDate.DateInput) : "",
    },
  });

  return (
    <FormGroup id={UPDATE_TICKET_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <SelectSprintsInput control={control} name="sprintId" projectId={ticket.projectId} />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="Ticket Title"
          label="Title"
          error={errors.title?.message}
          {...register("title")}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <Select label="Type" error={errors.type?.message} {...register("type")}>
          <Select.Option value="">Select a type</Select.Option>
          {TICKET_TYPES.map((type) => (
            <Select.Option key={type} value={type}>
              {type.charAt(0).toUpperCase() + type.slice(1)}
            </Select.Option>
          ))}
        </Select>
      </FormGroup.Field>
      <FormGroup.Field>
        <Select label="Priority" error={errors.priority?.message} {...register("priority")}>
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
