import { useState } from "react";
import { DRAWER_ROUTES } from "@/constants/drawer-routes";
import { dateFormat, FormatDate } from "@/lib";
import { useListSprints, useListBoards } from "@/hooks/use-api";
import { useDrawer } from "@/providers/drawer";
import { MenuIcon, PlusIcon, ChevronDownIcon, ChevronRightIcon } from "@versaur/icons";
import { PageContent, Table } from "@versaur/react/blocks";
import { ButtonIcon, Text } from "@versaur/react/primitive";
import { useOutletContext } from "react-router";
import type { ProjectDetailContextType } from "../project-detail-layout";

export const ProjectSprintsPage = () => {
  const { projectId } = useOutletContext<ProjectDetailContextType>();
  const { openDrawer } = useDrawer();
  const [expandedSprintId, setExpandedSprintId] = useState<string | null>(null);
  const [sprints, error, { isLoading }] = useListSprints({ projectId: [projectId] });

  const handleOnEditSprint = (sprintId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_SPRINT, { sprintId });
  };

  const handleOnAddBoardClick = (sprintId: string) => {
    openDrawer(DRAWER_ROUTES.CREATE_BOARD, { sprintId });
  };

  const handleOnEditBoard = (boardId: string) => {
    openDrawer(DRAWER_ROUTES.UPDATE_BOARD, { boardId });
  };

  if (isLoading) {
    return <PageContent>Loading sprints...</PageContent>;
  }

  if (error) {
    return <PageContent>Error loading sprints</PageContent>;
  }

  const sprintsList = sprints?.items || [];

  return (
    <PageContent>
      <Table columns="2fr 1fr 1fr 1fr 100px">
        <Table.Header>
          <Table.Col as="th">Name</Table.Col>
          <Table.Col as="th">Status</Table.Col>
          <Table.Col as="th">Planned Start</Table.Col>
          <Table.Col as="th">Planned End</Table.Col>
          <Table.Col as="th">Actions</Table.Col>
        </Table.Header>
        <Table.Body>
          {sprintsList.map((sprint) => (
            <SprintRow
              key={sprint.id}
              sprint={sprint}
              isExpanded={expandedSprintId === sprint.id}
              onToggleExpand={() =>
                setExpandedSprintId(expandedSprintId === sprint.id ? null : sprint.id)
              }
              onEditSprint={handleOnEditSprint}
              onAddBoard={handleOnAddBoardClick}
              onEditBoard={handleOnEditBoard}
            />
          ))}
        </Table.Body>
      </Table>
    </PageContent>
  );
};

interface SprintRowProps {
  sprint: any;
  isExpanded: boolean;
  onToggleExpand: () => void;
  onEditSprint: (sprintId: string) => void;
  onAddBoard: (sprintId: string) => void;
  onEditBoard: (boardId: string) => void;
}

function SprintRow({
  sprint,
  isExpanded,
  onToggleExpand,
  onEditSprint,
  onAddBoard,
  onEditBoard,
}: SprintRowProps) {
  const [boards] = useListBoards({ sprintId: [sprint.id] });

  const boardsList = boards?.items || [];

  return (
    <>
      <Table.Row key={sprint.id}>
        <Table.Col as="td">
          <div style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
            <ButtonIcon
              as={isExpanded ? ChevronDownIcon : ChevronRightIcon}
              variant="ghost"
              size="small"
              onClick={onToggleExpand}
              aria-label={isExpanded ? "Collapse" : "Expand"}
            />
            {sprint.name}
          </div>
        </Table.Col>
        <Table.Col as="td">{sprint.status}</Table.Col>
        <Table.Col as="td">
          {sprint.plannedStartedAt
            ? dateFormat(sprint.plannedStartedAt, FormatDate.ShortDate)
            : "-"}
        </Table.Col>
        <Table.Col as="td">
          {sprint.plannedCompletedAt
            ? dateFormat(sprint.plannedCompletedAt, FormatDate.ShortDate)
            : "-"}
        </Table.Col>
        <Table.Col as="td" variant="action">
          <Table.Action icon={MenuIcon}>
            <Table.ActionItem onClick={() => onEditSprint(sprint.id)}>Edit Sprint</Table.ActionItem>
          </Table.Action>
        </Table.Col>
      </Table.Row>

      {isExpanded && (
        <>
          <Table.Row>
            <Table.Col
              as="td"
              style={{ paddingLeft: "3rem", "--_bg-row": "var(--color-neutral-light)" } as any}
            >
              <Text size="xs" weight="medium">
                Name
              </Text>
            </Table.Col>
            <Table.Col as="td" style={{ "--_bg-row": "var(--color-neutral-light)" } as any} />
            <Table.Col
              as="td"
              style={{ "--_bg-row": "var(--color-neutral-light)" } as any}
            ></Table.Col>
            <Table.Col
              as="td"
              style={{ "--_bg-row": "var(--color-neutral-light)" } as any}
            ></Table.Col>
            <Table.Col
              as="td"
              variant="action"
              style={{ "--_bg-row": "var(--color-neutral-light)" } as any}
            >
              <ButtonIcon
                aria-label="Add board"
                as={PlusIcon}
                variant="ghost"
                size="small"
                onClick={() => onAddBoard(sprint.id)}
              />
            </Table.Col>
          </Table.Row>

          {boardsList.length === 0 ? (
            <Table.Row>
              <Table.Col
                as="td"
                style={{ paddingLeft: "3rem", "--_bg-row": "var(--color-neutral-light)" } as any}
              >
                <Text size="xs">No boards yet</Text>
              </Table.Col>
              <Table.Col as="td" style={{ "--_bg-row": "var(--color-neutral-light)" } as any} />
              <Table.Col
                as="td"
                style={{ "--_bg-row": "var(--color-neutral-light)" } as any}
              ></Table.Col>
              <Table.Col
                as="td"
                style={{ "--_bg-row": "var(--color-neutral-light)" } as any}
              ></Table.Col>
              <Table.Col
                as="td"
                style={{ "--_bg-row": "var(--color-neutral-light)" } as any}
              ></Table.Col>
            </Table.Row>
          ) : (
            boardsList.map((board) => (
              <Table.Row key={board.id}>
                <Table.Col
                  as="td"
                  style={{ paddingLeft: "3rem", "--_bg-row": "var(--color-neutral-light)" } as any}
                >
                  <Text size="xs">{board.name}</Text>
                </Table.Col>
                <Table.Col as="td" style={{ "--_bg-row": "var(--color-neutral-light)" } as any} />
                <Table.Col
                  as="td"
                  style={{ "--_bg-row": "var(--color-neutral-light)" } as any}
                ></Table.Col>
                <Table.Col
                  as="td"
                  style={{ "--_bg-row": "var(--color-neutral-light)" } as any}
                ></Table.Col>
                <Table.Col
                  as="td"
                  variant="action"
                  style={{ "--_bg-row": "var(--color-neutral-light)" } as any}
                >
                  <Table.Action icon={MenuIcon}>
                    <Table.ActionItem onClick={() => onEditBoard(board.id)}>Edit</Table.ActionItem>
                  </Table.Action>
                </Table.Col>
              </Table.Row>
            ))
          )}
        </>
      )}
    </>
  );
}
