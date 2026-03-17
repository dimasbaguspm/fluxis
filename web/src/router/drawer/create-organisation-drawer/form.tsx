import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { CREATE_ORGANISATION_FORM_ID } from "./constants";
import type { CreateOrganisationFormInputs } from "./types";

interface CreateOrganisationFormProps {
  onSubmit: (data: CreateOrganisationFormInputs) => void;
}

export const CreateOrganisationForm = ({
  onSubmit,
}: CreateOrganisationFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm<CreateOrganisationFormInputs>({
    mode: "all",
    defaultValues: { name: "" },
  });

  return (
    <FormGroup id={CREATE_ORGANISATION_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <TextInput
          placeholder="Organisation Name"
          label="Name"
          required
          error={errors.name?.message}
          {...register("name", {
            required: "This field required",
            minLength: {
              value: 1,
              message: "Organisation name must not be empty",
            },
          })}
        />
      </FormGroup.Field>
    </FormGroup>
  );
};
