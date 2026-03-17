import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { CREATE_BOARD_FORM_ID } from "./constants";
import type { CreateBoardFormInputs } from "./types";

interface CreateBoardFormProps {
  onSubmit: (data: CreateBoardFormInputs) => void;
  sprintId: string;
}

export const CreateBoardForm = ({
  onSubmit,
  sprintId,
}: CreateBoardFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm<CreateBoardFormInputs>({
    mode: "all",
    defaultValues: {
      sprintId,
      name: "",
    },
  });

  return (
    <FormGroup id={CREATE_BOARD_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <TextInput
          placeholder="Sprint ID"
          label="Sprint ID"
          required
          disabled
          error={errors.sprintId?.message}
          {...register("sprintId", {
            required: "Sprint ID is required",
          })}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="Board Name"
          label="Name"
          required
          error={errors.name?.message}
          {...register("name", {
            required: "Board name is required",
          })}
        />
      </FormGroup.Field>
    </FormGroup>
  );
};
