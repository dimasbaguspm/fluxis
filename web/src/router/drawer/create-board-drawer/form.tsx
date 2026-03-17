import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { SelectSprintsInput } from "@/components/ui";
import { CREATE_BOARD_FORM_ID } from "./constants";
import type { CreateBoardFormInputs } from "./types";

interface CreateBoardFormProps {
  onSubmit: (data: CreateBoardFormInputs) => void;
  sprintId?: string;
}

export const CreateBoardForm = ({
  onSubmit,
  sprintId,
}: CreateBoardFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
    control,
  } = useForm<CreateBoardFormInputs>({
    mode: "all",
    defaultValues: {
      sprintId: sprintId ?? "",
      name: "",
    },
  });

  return (
    <FormGroup id={CREATE_BOARD_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <SelectSprintsInput control={control} name="sprintId" required />
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
