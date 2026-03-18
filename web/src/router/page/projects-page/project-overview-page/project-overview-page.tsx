import { dateFormat, FormatDate } from "@/lib";
import { PageContent, AttributeList } from "@versaur/react/blocks";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";

export const ProjectOverviewPage = () => {
  const { project } = useOutletContext<ProjectDetailContextType>();

  return (
    <PageContent>
      <AttributeList columns="repeat(2, 1fr)">
        <AttributeList.Item title="Name" area="span 2">
          {project?.name || "-"}
        </AttributeList.Item>
        <AttributeList.Item title="Key">
          {project?.key || "-"}
        </AttributeList.Item>
        <AttributeList.Item title="Visibility">
          {project?.visibility || "-"}
        </AttributeList.Item>
        <AttributeList.Item title="Organization ID" area="span 2">
          {project?.orgId || "-"}
        </AttributeList.Item>
        <AttributeList.Item title="Description" area="span 2">
          {project?.description || "-"}
        </AttributeList.Item>
        <AttributeList.Item title="Created">
          {project?.createdAt ? dateFormat(project.createdAt, FormatDate.ShortDate) : "-"}
        </AttributeList.Item>
      </AttributeList>
    </PageContent>
  );
};
