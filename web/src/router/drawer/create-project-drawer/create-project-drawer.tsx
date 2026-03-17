import type { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useCreateProject } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { ButtonGroup, Drawer } from "@versaur/react/blocks";
import { Banner, Button } from "@versaur/react/primitive";
import { CREATE_PROJECT_FORM_ID } from "./constants";
import { CreateProjectForm } from "./form";
import type { CreateProjectFormInputs } from "./types";
import { cx } from "@/lib";
import { vMb4 } from "@versaur/core/utilities";

export const CreateProjectDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.CREATE_PROJECT>();
  const [createProject, err, { isPending }] = useCreateProject();

  const onSubmit = async (data: CreateProjectFormInputs) => {
    await createProject({
      data: {
        orgId: data.orgId,
        name: data.name,
        key: data.key,
        visibility: data.visibility,
        description: data.description,
      },
    });
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Create Project</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        {err && (
          <Banner variant="warning" className={cx(vMb4)}>
            {err?.message}
          </Banner>
        )}
        <CreateProjectForm orgId={params?.orgId} onSubmit={onSubmit} />
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button
            type="submit"
            form={CREATE_PROJECT_FORM_ID}
            loading={isPending}
            disabled={isPending}
          >
            Create
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
