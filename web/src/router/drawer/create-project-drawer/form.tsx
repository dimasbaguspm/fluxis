import { FormGroup } from "@versaur/react/blocks";
import { Select, TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { SelectOrganisationsInput } from "@/components/ui/select-organisations-input";
import { CREATE_PROJECT_FORM_ID } from "./constants";
import type { CreateProjectFormInputs } from "./types";

interface CreateProjectFormProps {
  onSubmit: (data: CreateProjectFormInputs) => void;
  orgId?: string;
}

export const CreateProjectForm = ({
  onSubmit,
  orgId,
}: CreateProjectFormProps) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
    control,
  } = useForm<CreateProjectFormInputs>({
    mode: "onBlur",
    defaultValues: {
      orgId: orgId ?? "",
      name: "",
      key: "",
      visibility: "private",
      description: "",
    },
  });

  return (
    <FormGroup id={CREATE_PROJECT_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <SelectOrganisationsInput control={control} name="orgId" required />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="Project Name"
          label="Name"
          required
          error={errors.name?.message}
          {...register("name", {
            required: "Project name is required",
          })}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="Project Key"
          label="Key"
          required
          error={errors.key?.message}
          {...register("key", {
            required: "Project key is required",
            maxLength: {
              value: 10,
              message: "Key must not exceed 10 characters",
            },
          })}
        />
      </FormGroup.Field>
      <FormGroup.Field>
        <Select {...register("visibility")} label="Visibility" required>
          <Select.Option value="public">Public</Select.Option>
          <Select.Option value="private">Private</Select.Option>
        </Select>
      </FormGroup.Field>
      <FormGroup.Field>
        <TextInput
          placeholder="Description (optional)"
          label="Description"
          error={errors.description?.message}
          {...register("description")}
        />
      </FormGroup.Field>
    </FormGroup>
  );
};
