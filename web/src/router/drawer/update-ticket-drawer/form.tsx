import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { SelectSprintsInput } from "@/components/ui";
import { UPDATE_TICKET_FORM_ID } from "./constants";
import type { UpdateTicketFormInputs } from "./types";
import type { DomainTicketModel } from "@/interfaces/openapi.generated";

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
      dueDate: ticket.dueDate ?? "",
    },
  });

  return (
    <FormGroup id={UPDATE_TICKET_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <TextInput
          placeholder="Ticket Title"
          label="Title"
          error={errors.title?.message}
          {...register("title")}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="bug / story / task / epic"
          label="Type"
          error={errors.type?.message}
          {...register("type")}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="low / medium / high / critical"
          label="Priority"
          error={errors.priority?.message}
          {...register("priority")}
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
