import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { dateFormat, FormatDate } from "@/lib";
import { useListTickets } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { SelectProjectsInput } from "@/components/ui";
import { MenuIcon, PlusIcon } from "@versaur/icons";
import { PageContent, PageHeader, Table } from "@versaur/react/blocks";
import { ButtonIcon, Text } from "@versaur/react/primitive";
import { useForm, useWatch } from "react-hook-form";

export const TicketsPage = () => {
  const { openDrawer } = useDrawer();

  const { control } = useForm({
    defaultValues: {
      projectId: "",
    },
  });

  const selectedProjectId = useWatch({ control, name: "projectId" });

  const [tickets, error, { isLoading }] = useListTickets(
    selectedProjectId ? { projectId: [selectedProjectId] } : undefined,
    { enabled: !!selectedProjectId },
  );

  const handleOnAddButtonClick = () => {
    openDrawer(DRAWER_ROUTES.CREATE_TICKET, { projectId: selectedProjectId });
  };

  const handleOnEditTicket = (ticketId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_TICKET, { ticketId: ticketId });
  };

  const ticketsList = tickets?.items || [];

  return (
    <>
      <PageHeader
        title={
          <PageHeader.Title
            action={
              selectedProjectId && (
                <ButtonIcon
                  aria-label="Add ticket"
                  as={PlusIcon}
                  variant="ghost"
                  onClick={handleOnAddButtonClick}
                />
              )
            }
          >
            Tickets
          </PageHeader.Title>
        }
        subtitle={<PageHeader.Subtitle>Manage your tickets</PageHeader.Subtitle>}
      />

      <PageContent>
        <div style={{ marginBottom: "2rem", maxWidth: "400px" }}>
          <SelectProjectsInput control={control} name="projectId" required />
        </div>

        {!selectedProjectId ? (
          <div style={{ padding: "2rem", textAlign: "center", color: "#666" }}>
            Select a project to view tickets
          </div>
        ) : (
          <>
            {isLoading && <div>Loading tickets...</div>}

            {error && <div>Error loading tickets</div>}

            {!isLoading && !error && (
              <Table columns="40px 1fr 100px 100px 1fr 100px">
                <Table.Toolbar
                  leftContent={(selectedIds) => (
                    <Text size="xs">
                      {selectedIds.size > 0 ? `${selectedIds.size} selected` : "No selection"}
                    </Text>
                  )}
                  rightContent={(selectedIds) => (
                    <Text size="xs">
                      {selectedIds.size > 0 ? `${selectedIds.size} selected` : "No selection"}
                    </Text>
                  )}
                />
                <Table.Header>
                  <Table.Col as="th" variant="checkbox">
                    <Table.Checkbox isMain />
                  </Table.Col>
                  <Table.Col as="th">Title</Table.Col>
                  <Table.Col as="th">Type</Table.Col>
                  <Table.Col as="th">Priority</Table.Col>
                  <Table.Col as="th">Updated</Table.Col>
                  <Table.Col as="th">Actions</Table.Col>
                </Table.Header>
                <Table.Body>
                  {ticketsList.map((ticket) => (
                    <Table.Row key={ticket.id}>
                      <Table.Col as="td" variant="checkbox">
                        <Table.Checkbox rowId={ticket.id} />
                      </Table.Col>
                      <Table.Col as="td">{ticket.title}</Table.Col>
                      <Table.Col as="td">{ticket.type}</Table.Col>
                      <Table.Col as="td">{ticket.priority}</Table.Col>
                      <Table.Col as="td">
                        {dateFormat(ticket.updatedAt, FormatDate.ShortDate)}
                      </Table.Col>
                      <Table.Col as="td" variant="action">
                        <Table.Action icon={MenuIcon}>
                          <Table.ActionItem onClick={() => handleOnEditTicket(ticket.id)}>
                            Edit
                          </Table.ActionItem>
                        </Table.Action>
                      </Table.Col>
                    </Table.Row>
                  ))}
                </Table.Body>
              </Table>
            )}
          </>
        )}
      </PageContent>
    </>
  );
};
