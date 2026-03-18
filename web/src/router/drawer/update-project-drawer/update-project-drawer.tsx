import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { useGetProject, useUpdateProject } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { ButtonGroup, Drawer, NoResults } from "@versaur/react/blocks";
import { Banner, Button, Loader } from "@versaur/react/primitive";
import { UPDATE_PROJECT_FORM_ID } from "./constants";
import { UpdateProjectForm } from "./form";
import type { UpdateProjectFormInputs } from "./types";
import { When } from "@/lib/when";
import { SearchXIcon } from "@versaur/icons";

export const UpdateProjectDrawer = () => {
  const { closeDrawer, params } = useDrawer<typeof DRAWER_ROUTES.UPDATE_PROJECT>();
  const projectId = params?.projectId ?? "";

  const [project, , { isLoading: isLoadingProject }] = useGetProject(projectId);
  const [updateProject, err, { isPending }] = useUpdateProject();

  const onSubmit = async (data: UpdateProjectFormInputs) => {
    if (!projectId) return;
    await updateProject({
      projectId,
      data: {
        name: data.name,
        description: data.description,
      },
    });
    closeDrawer();
  };

  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Update Project</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <When condition={isLoadingProject}>
          <Loader type="bar" />
        </When>
        <When condition={!isLoadingProject}>
          <When condition={!project}>
            <NoResults
              icon={SearchXIcon}
              title="Project not found"
              subtitle="We couldn't find the project you're looking for"
            />
          </When>
          <When condition={!!project}>
            <When condition={err?.message}>
              <Banner variant="warning">{err?.message}</Banner>
            </When>
            <UpdateProjectForm project={project!} onSubmit={onSubmit} />
          </When>
        </When>
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost" onClick={closeDrawer}>
            Cancel
          </Button>
          <Button
            type="submit"
            form={UPDATE_PROJECT_FORM_ID}
            loading={isPending || isLoadingProject}
            disabled={isPending || isLoadingProject}
          >
            Update
          </Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
