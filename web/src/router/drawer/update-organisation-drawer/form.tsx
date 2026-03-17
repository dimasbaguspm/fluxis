import { FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { useForm } from "react-hook-form";
import { UPDATE_ORGANISATION_FORM_ID } from "./constants";
import type { UpdateOrganisationFormInputs } from "./types";
import type { DomainOrganisationModel } from "@/interfaces/openapi.generated";

interface UpdateOrganisationFormProps {
  onSubmit: (data: UpdateOrganisationFormInputs) => void;
  organisation: DomainOrganisationModel;
}

export const UpdateOrganisationForm = ({ onSubmit, organisation }: UpdateOrganisationFormProps) => {
  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm<UpdateOrganisationFormInputs>({
    mode: "all",
    defaultValues: {
      name: organisation.name,
    },
  });

  return (
    <FormGroup id={UPDATE_ORGANISATION_FORM_ID} onSubmit={handleSubmit(onSubmit)}>
      <FormGroup.Field>
        <TextInput
          placeholder="Organisation Name"
          label="Name"
          required
          error={errors.name?.message}
          {...register("name", {
            required: "This field is required",
            minLength: {
              value: 1,
              message: "This field must not be empty",
            },
          })}
        />
      </FormGroup.Field>
    </FormGroup>
  );
};
