import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useCreateOrg } from "@/hooks/use-api";
import { cx } from "@/lib/cx";
import { useDrawer } from "@/providers/drawer";
import { vGap4 } from "@versaur/core/utilities";
import { ButtonGroup, Drawer, FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { Banner, Button } from "@versaur/react/primitive";
import { useForm } from "react-hook-form";

interface CreateOrganisationFormInputs {
  name: string;
}

export const CreateOrganisationDrawer = () => {
  const { closeDrawer } = useDrawer<typeof DRAWER_ROUTES.CREATE_ORGANISATION>();
  const [createOrg, err, { isPending }] = useCreateOrg();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateOrganisationFormInputs>({
    defaultValues: { name: "" },
  });

  const onSubmit = async (data: CreateOrganisationFormInputs) => {
    await createOrg(data);
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Create Organisation</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <FormGroup onSubmit={handleSubmit(onSubmit)} className={cx(vGap4)}>
          {err && (
            <FormGroup.Field>
              <Banner variant="warning">{err.message}</Banner>
            </FormGroup.Field>
          )}
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
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button onClick={handleSubmit(onSubmit)} loading={isPending} disabled={isPending}>
            {isPending ? "Creating..." : "Create"}
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
