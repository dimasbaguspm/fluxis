import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { SelectSprintsInput } from "@/components/ui";
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
    control,
  } = useForm<UpdateBoardFormInputs>({
    mode: "all",
    defaultValues: {
      sprintId: board.sprintId,
      name: board.name,
    },
  });

  return (
    <FormGroup id={UPDATE_BOARD_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <SelectSprintsInput control={control} name="sprintId" required />
      </FormGroup.Field>
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
