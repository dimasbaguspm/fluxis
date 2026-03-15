import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetOrg, useUpdateOrg } from "@/hooks/use-api";
import { cx } from "@/lib/cx";
import { useDrawer } from "@/providers/drawer";
import { vGap4 } from "@versaur/core/utilities";
import { ButtonGroup, Drawer, FormGroup } from "@versaur/react/blocks";
import { TextInput } from "@versaur/react/forms";
import { Banner, Button } from "@versaur/react/primitive";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

interface UpdateOrganisationFormInputs {
  name: string;
}

export const UpdateOrganisationDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_ORGANISATION>();
  const orgId = params?.orgId ?? "";

  const [org, , { isLoading: isLoadingOrg }] = useGetOrg(orgId);
  const [updateOrg, err, { isPending }] = useUpdateOrg();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
  } = useForm<UpdateOrganisationFormInputs>({
    defaultValues: { name: "" },
  });

  useEffect(() => {
    if (org?.name) {
      setValue("name", org.name);
    }
  }, [org?.name, setValue]);

  const onSubmit = async (data: UpdateOrganisationFormInputs) => {
    if (!orgId) return;
    await updateOrg({
      id: orgId,
      data: { name: data.name },
    });
    reset();
    closeDrawer();
  };

  const errorMessage = (() => {
    if (!err) return null;
    if (typeof err === "object" && err !== null && "message" in err) {
      return String((err as any).message);
    }
    return "Failed to update organisation";
  })();

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Update Organisation</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <FormGroup onSubmit={handleSubmit(onSubmit)} className={cx(vGap4)}>
          {errorMessage ? (
            <FormGroup.Field>
              <Banner variant="warning">{errorMessage}</Banner>
            </FormGroup.Field>
          ) : null}
          <FormGroup.Field>
            <TextInput
              placeholder="Organisation Name"
              label="Organisation Name"
              required
              disabled={isLoadingOrg}
              error={errors.name?.message}
              {...register("name", {
                required: "Organisation name is required",
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
          <Button
            onClick={handleSubmit(onSubmit)}
            loading={isPending || isLoadingOrg}
            disabled={isPending || isLoadingOrg}
          >
            {isPending ? "Updating..." : isLoadingOrg ? "Loading..." : "Update"}
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
