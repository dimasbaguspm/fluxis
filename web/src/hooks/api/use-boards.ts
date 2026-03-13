import type {
  DomainBoardColumnCreateModel,
  DomainBoardColumnModel,
  DomainBoardColumnsPagedModel,
  DomainBoardColumnUpdateModel,
  DomainBoardCreateModel,
  DomainBoardModel,
  DomainBoardsPagedModel,
  DomainBoardUpdateModel,
} from "@interfaces/openapi.generated";
import type { UseQueryOptions } from "@tanstack/react-query";
import { useApiMutation } from "../use-api-mutation";
import { useApiQuery } from "../use-api-query";

interface ListBoardsParams {
  id?: string[];
  name?: string;
  pageNumber?: number;
  pageSize?: number;
  sprintId?: string[];
}

/**
 * List boards with pagination and filtering
 */
export function useListBoards(
  params?: ListBoardsParams,
  options?: Omit<UseQueryOptions<DomainBoardsPagedModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["boards", params],
    {
      method: "GET",
      path: "/boards",
      params,
    },
    options,
  );
}

/**
 * Get a single board by ID
 */
export function useGetBoard(
  boardId: string,
  options?: Omit<UseQueryOptions<DomainBoardModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["boards", boardId],
    {
      method: "GET",
      path: `/boards/${boardId}`,
    },
    { ...options, enabled: !!boardId },
  );
}

/**
 * Create a new board
 */
export function useCreateBoard() {
  return useApiMutation<DomainBoardModel, unknown, DomainBoardCreateModel>((variables) => ({
    method: "POST",
    path: "/boards",
    body: variables,
  }));
}

/**
 * Update a board
 */
export function useUpdateBoard() {
  return useApiMutation<
    DomainBoardModel,
    unknown,
    { boardId: string; data: DomainBoardUpdateModel }
  >((variables) => ({
    method: "PATCH",
    path: `/boards/${variables.boardId}`,
    body: variables.data,
  }));
}

/**
 * Delete a board
 */
export function useDeleteBoard() {
  return useApiMutation<void, unknown, string>((boardId) => ({
    method: "DELETE",
    path: `/boards/${boardId}`,
  }));
}

/**
 * Reorder boards
 */
export function useReorderBoards() {
  return useApiMutation<DomainBoardModel[], unknown, { sprintId: string; boardIds: string[] }>(
    (variables) => ({
      method: "PATCH",
      path: `/boards/reorder`,
      body: { sprintId: variables.sprintId, boardIds: variables.boardIds },
    }),
  );
}

interface ListBoardColumnsParams {
  boardId?: string[];
  id?: string[];
  name?: string;
  pageNumber?: number;
  pageSize?: number;
}

/**
 * List board columns
 */
export function useListBoardColumns(
  boardId: string,
  params?: ListBoardColumnsParams,
  options?: Omit<UseQueryOptions<DomainBoardColumnsPagedModel>, "queryKey" | "queryFn">,
) {
  return useApiQuery(
    ["boards", boardId, "columns", params],
    {
      method: "GET",
      path: `/boards/${boardId}/columns`,
      params,
    },
    { ...options, enabled: !!boardId },
  );
}

/**
 * Create a board column
 */
export function useCreateBoardColumn() {
  return useApiMutation<
    DomainBoardColumnModel,
    unknown,
    { boardId: string; data: DomainBoardColumnCreateModel }
  >((variables) => ({
    method: "POST",
    path: `/boards/${variables.boardId}/columns`,
    body: variables.data,
  }));
}

/**
 * Update a board column
 */
export function useUpdateBoardColumn() {
  return useApiMutation<
    DomainBoardColumnModel,
    unknown,
    { boardId: string; columnId: string; data: DomainBoardColumnUpdateModel }
  >((variables) => ({
    method: "PATCH",
    path: `/boards/${variables.boardId}/columns/${variables.columnId}`,
    body: variables.data,
  }));
}

/**
 * Delete a board column
 */
export function useDeleteBoardColumn() {
  return useApiMutation<void, unknown, { boardId: string; columnId: string }>((variables) => ({
    method: "DELETE",
    path: `/boards/${variables.boardId}/columns/${variables.columnId}`,
  }));
}

/**
 * Reorder board columns
 */
export function useReorderBoardColumns() {
  return useApiMutation<
    DomainBoardColumnModel[],
    unknown,
    { boardId: string; columnIds: string[] }
  >((variables) => ({
    method: "PATCH",
    path: `/boards/${variables.boardId}/columns/reorder`,
    body: { columnIds: variables.columnIds },
  }));
}
