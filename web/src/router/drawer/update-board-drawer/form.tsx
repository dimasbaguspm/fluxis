import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { UPDATE_BOARD_FORM_ID } from "./constants";
import type { UpdateBoardFormInputs } from "./types";
import type { DomainBoardModel } from "@/interfaces/openapi.generated";

interface UpdateBoardFormProps {
  onSubmit: (data: UpdateBoardFormInputs) => void;
  board: DomainBoardModel;
}

export const UpdateBoardForm = ({ onSubmit, board }: UpdateBoardFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm<UpdateBoardFormInputs>({
    mode: "all",
    defaultValues: {
      name: board.name,
    },
  });

  return (
    <FormGroup id={UPDATE_BOARD_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <TextInput
          placeholder="Name"
          label="Name"
          required
          error={errors.name?.message}
          {...register("name", {
            required: "This field is required",
          })}
        />
      </FormGroup.Field>
    </FormGroup>
  );
};
