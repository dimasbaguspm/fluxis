import { useEffect, useRef } from "react";
import { MenuIcon, HamburgerMenuIcon } from "@versaur/icons";
import { Table } from "@versaur/react/blocks";
import type { DomainTicketModel } from "@/interfaces/openapi.generated";
import { draggable } from "@atlaskit/pragmatic-drag-and-drop/element/adapter";
import { combine } from "@atlaskit/pragmatic-drag-and-drop/combine";

interface TicketsTableProps {
  tickets: DomainTicketModel[];
  onEditTicket: (ticketId: string) => void;
  isLoading?: boolean;
  groupId?: string | null;
}

export const TicketsTable = ({ tickets, onEditTicket, isLoading, groupId }: TicketsTableProps) => {
  const columns = "40px 70px 1fr 65px 70px 90px 80px";

  if (isLoading) {
    return (
      <Table columns={columns}>
        <Table.Header>
          <Table.Col as="th"></Table.Col>
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
          <Table.Col as="th"></Table.Col>
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
        <Table.Col as="th"></Table.Col>
        <Table.Col as="th">Key</Table.Col>
        <Table.Col as="th">Title</Table.Col>
        <Table.Col as="th">Type</Table.Col>
        <Table.Col as="th">Priority</Table.Col>
        <Table.Col as="th">Points</Table.Col>
        <Table.Col as="th">Actions</Table.Col>
      </Table.Header>
      <Table.Body>
        {tickets.map((ticket) => (
          <TicketRow
            key={ticket.id}
            ticket={ticket}
            onEditTicket={onEditTicket}
            groupId={groupId}
          />
        ))}
      </Table.Body>
    </Table>
  );
};

interface TicketRowProps {
  ticket: DomainTicketModel;
  onEditTicket: (ticketId: string) => void;
  groupId?: string | null;
}

const TicketRow = ({ ticket, onEditTicket, groupId }: TicketRowProps) => {
  const dragRef = useRef<HTMLDivElement>(null);
  const rowRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!dragRef.current || !rowRef.current) return;

    return combine(
      draggable({
        element: rowRef.current,
        onDragStart: () => {
          // Optional: Add visual feedback during drag
        },
        getInitialData: () => ({
          ticketId: ticket.id,
          sourceGroupId: groupId,
          type: "ticket",
        }),
      })
    );
  }, [ticket.id, groupId]);

  return (
    <Table.Row ref={rowRef}>
      <Table.Col as="td" style={{ cursor: "grab", display: "flex", alignItems: "center", justifyContent: "center" }} ref={dragRef}>
        <HamburgerMenuIcon style={{ width: "16px", height: "16px", color: "#999" }} />
      </Table.Col>
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
  );
};
