import { MenuIcon } from "@versaur/icons";
import { Table } from "@versaur/react/blocks";
import type { DomainTicketModel, DomainSprintModel } from "@/interfaces/openapi.generated";

interface TicketsTableProps {
  tickets: DomainTicketModel[];
  onEditTicket: (ticketId: string) => void;
  isLoading?: boolean;
  sprint?: DomainSprintModel;
}

export const TicketsTable = ({ tickets, onEditTicket, isLoading }: TicketsTableProps) => {
  const columns = "70px 1fr 65px 70px 90px 80px";

  if (isLoading) {
    return (
      <Table columns={columns}>
        <Table.Header>
          <Table.Col as="th">Key</Table.Col>
          <Table.Col as="th">Title</Table.Col>
          <Table.Col as="th">Type</Table.Col>
          <Table.Col as="th">Priority</Table.Col>
          <Table.Col as="th">Points</Table.Col>
          <Table.Col as="th">Actions</Table.Col>
        </Table.Header>
        <Table.Body>
          <Table.Row>
            <Table.Col as="td" style={{ gridColumn: "1 / -1", textAlign: "center", padding: "2rem" }}>
              Loading tickets...
            </Table.Col>
          </Table.Row>
        </Table.Body>
      </Table>
    );
  }

  if (tickets.length === 0) {
    return (
      <Table columns={columns}>
        <Table.Header>
          <Table.Col as="th">Key</Table.Col>
          <Table.Col as="th">Title</Table.Col>
          <Table.Col as="th">Type</Table.Col>
          <Table.Col as="th">Priority</Table.Col>
          <Table.Col as="th">Points</Table.Col>
          <Table.Col as="th">Actions</Table.Col>
        </Table.Header>
        <Table.Body>
          <Table.Row>
            <Table.Col as="td" style={{ gridColumn: "1 / -1", textAlign: "center", padding: "2rem", color: "#999" }}>
              No tickets
            </Table.Col>
          </Table.Row>
        </Table.Body>
      </Table>
    );
  }

  return (
    <Table columns={columns}>
      <Table.Header>
        <Table.Col as="th">Key</Table.Col>
        <Table.Col as="th">Title</Table.Col>
        <Table.Col as="th">Type</Table.Col>
        <Table.Col as="th">Priority</Table.Col>
        <Table.Col as="th">Points</Table.Col>
        <Table.Col as="th">Actions</Table.Col>
      </Table.Header>
      <Table.Body>
        {tickets.map((ticket) => (
          <Table.Row key={ticket.id}>
            <Table.Col as="td">{ticket.key}</Table.Col>
            <Table.Col as="td">{ticket.title}</Table.Col>
            <Table.Col as="td">{ticket.type}</Table.Col>
            <Table.Col as="td">{ticket.priority}</Table.Col>
            <Table.Col as="td">{ticket.storyPoints ?? "-"}</Table.Col>
            <Table.Col as="td" variant="action">
              <Table.Action icon={MenuIcon}>
                <Table.ActionItem onClick={() => onEditTicket(ticket.id)}>
                  Edit
                </Table.ActionItem>
              </Table.Action>
            </Table.Col>
          </Table.Row>
        ))}
      </Table.Body>
    </Table>
  );
};
