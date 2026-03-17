import type { QueryClient } from "@tanstack/react-query";
import type { NotifierEvent } from "@/providers/notifier";
import { NotifierEventType } from "@/providers/notifier";
import { queryKeys } from "./query-keys";

/**
 * Invalidates query cache based on notifier event type
 * Automatically invalidates all queries for the affected entity regardless of ID
 *
 * @example
 * const queryClient = useQueryClient();
 * invalidateQueriesByEvent(queryClient, notifierEvent);
 */
export async function invalidateQueriesByEvent(
  queryClient: QueryClient,
  event: NotifierEvent,
): Promise<void> {
  switch (event.type) {
    case NotifierEventType.Org:
      await queryClient.invalidateQueries({ queryKey: queryKeys.orgs.all });
      break;

    case NotifierEventType.OrgMember:
      await queryClient.invalidateQueries({ queryKey: queryKeys.orgs.all });
      break;

    case NotifierEventType.Project:
      await queryClient.invalidateQueries({ queryKey: queryKeys.projects.all });
      break;

    case NotifierEventType.Sprint:
      await queryClient.invalidateQueries({ queryKey: queryKeys.sprints.all });
      break;

    case NotifierEventType.Board:
      await queryClient.invalidateQueries({ queryKey: queryKeys.boards.all });
      break;

    case NotifierEventType.BoardColumn:
      await queryClient.invalidateQueries({ queryKey: queryKeys.boards.all });
      break;

    case NotifierEventType.Ticket:
      await queryClient.invalidateQueries({ queryKey: queryKeys.tickets.all });
      break;

    case NotifierEventType.User:
      await queryClient.invalidateQueries({ queryKey: queryKeys.user.all });
      break;

    case NotifierEventType.Auth:
      await queryClient.invalidateQueries({ queryKey: queryKeys.user.all });
      break;

    case NotifierEventType.Unknown:
      break;
  }
}
