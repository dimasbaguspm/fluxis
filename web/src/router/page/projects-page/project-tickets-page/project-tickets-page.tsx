import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { dateFormat, FormatDate } from "@/lib";
import { useListTickets } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { MenuIcon } from "@versaur/icons";
import { PageContent, Table } from "@versaur/react/blocks";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";

export const ProjectTicketsPage = () => {
  const { projectId } = useOutletContext<ProjectDetailContextType>();
  const { openDrawer } = useDrawer();
  const [tickets, error, { isLoading }] = useListTickets({ projectId: [projectId], pageSize: 20 });

  const handleOnEditTicket = (ticketId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_TICKET, { ticketId });
  };

  if (isLoading) {
    return <PageContent>Loading tickets...</PageContent>;
  }

  if (error) {
    return <PageContent>Error loading tickets</PageContent>;
  }

  const ticketsList = tickets?.items || [];

  return (
    <PageContent>
      <Table columns="1fr 100px 100px 1fr 100px">
        <Table.Header>
          <Table.Col as="th">Title</Table.Col>
          <Table.Col as="th">Type</Table.Col>
          <Table.Col as="th">Priority</Table.Col>
          <Table.Col as="th">Updated</Table.Col>
          <Table.Col as="th">Actions</Table.Col>
        </Table.Header>
        <Table.Body>
          {ticketsList.map((ticket) => (
            <Table.Row key={ticket.id}>
              <Table.Col as="td">{ticket.title}</Table.Col>
              <Table.Col as="td">{ticket.type}</Table.Col>
              <Table.Col as="td">{ticket.priority}</Table.Col>
              <Table.Col as="td">{dateFormat(ticket.updatedAt, FormatDate.ShortDate)}</Table.Col>
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
    </PageContent>
  );
};
