import { SearchIcon } from "@versaur/icons";
import { PageHeader } from "@versaur/react/blocks";
import { Badge, ButtonIcon } from "@versaur/react/primitive";

export const DashboardPage = () => {
  return (
    <>
      <PageHeader
        title={
          <PageHeader.Title
            action={<ButtonIcon aria-label="Search" as={SearchIcon} variant="ghost" />}
          >
            Dashboard
          </PageHeader.Title>
        }
        subtitle={
          <PageHeader.Subtitle additionalInfo={<Badge variant="info">8 total</Badge>}>
            Manage your projects and settings
          </PageHeader.Subtitle>
        }
      />
    </>
  );
};
