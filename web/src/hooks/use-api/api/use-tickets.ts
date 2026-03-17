import type {
  DomainTicketBoardMoveModel,
  DomainTicketCreateModel,
  DomainTicketModel,
  DomainTicketsPagedModel,
  DomainTicketUpdateModel,
  HttpxErrBlock,
} from "@interfaces/openapi.generated";
import type { UseQueryOptions } from "@tanstack/react-query";
import { useApiMutation } from "../use-api-mutation";
import { useApiQuery } from "../use-api-query";
import { queryKeys } from "../query-keys";

interface ListTicketsParams {
  boardId?: string[];
  id?: string[];
  pageNumber?: number;
  pageSize?: number;
  projectId?: string[];
  sprintId?: string[];
}

/**
 * List tickets with pagination and filtering
 */
export function useListTickets(
  params?: ListTicketsParams,
  options?: Omit<UseQueryOptions<DomainTicketsPagedModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    queryKeys.tickets.list(params),
    {
      method: "GET",
      path: "/tickets",
      params,
    },
    options,
  );
}

/**
 * Get a single ticket by ID
 */
export function useGetTicket(
  ticketId: string,
  options?: Omit<UseQueryOptions<DomainTicketModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    queryKeys.tickets.detail(ticketId),
    {
      method: "GET",
      path: `/tickets/${ticketId}`,
    },
    { ...options, enabled: !!ticketId },
  );
}

/**
 * Create a new ticket
 */
export function useCreateTicket() {
  return useApiMutation<
    DomainTicketModel,
    HttpxErrBlock,
    { projectId: string; data: DomainTicketCreateModel }
  >((variables) => ({
    method: "POST",
    path: "/tickets",
    body: { ...variables.data, projectId: variables.projectId },
  }));
}

/**
 * Update a ticket
 */
export function useUpdateTicket() {
  return useApiMutation<
    DomainTicketModel,
    HttpxErrBlock,
    { ticketId: string; data: DomainTicketUpdateModel }
  >((variables) => ({
    method: "PATCH",
    path: `/tickets/${variables.ticketId}`,
    body: variables.data,
  }));
}

/**
 * Delete a ticket
 */
export function useDeleteTicket() {
  return useApiMutation<void, HttpxErrBlock, string>((ticketId) => ({
    method: "DELETE",
    path: `/tickets/${ticketId}`,
  }));
}

/**
 * Move ticket to a board column
 */
export function useMoveBoardColumn() {
  return useApiMutation<
    DomainTicketModel,
    HttpxErrBlock,
    { ticketId: string; move: DomainTicketBoardMoveModel }
  >((variables) => ({
    method: "PATCH",
    path: `/tickets/${variables.ticketId}/move-board-column`,
    body: variables.move,
  }));
}

/**
 * Move ticket to a board
 */
export function useMoveToBoard() {
  return useApiMutation<
    DomainTicketModel,
    HttpxErrBlock,
    { ticketId: string; move: DomainTicketBoardMoveModel }
  >((variables) => ({
    method: "PATCH",
    path: `/tickets/${variables.ticketId}/move-to-board`,
    body: variables.move,
  }));
}

/**
 * Move ticket to a sprint
 */
export function useMoveToSprint() {
  return useApiMutation<DomainTicketModel, HttpxErrBlock, { ticketId: string; sprintId: string }>(
    (variables) => ({
      method: "PATCH",
      path: `/tickets/${variables.ticketId}/move-to-sprint`,
      body: { sprintId: variables.sprintId },
    }),
  );
}
