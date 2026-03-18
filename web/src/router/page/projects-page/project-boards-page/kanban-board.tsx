import { Text } from "@versaur/react/primitive";
import type { DomainBoardColumnModel } from "@/interfaces/openapi.generated";

interface KanbanBoardProps {
  boardId: string;
}

// Mock data matching DomainBoardColumnModel structure
const MOCK_COLUMNS: DomainBoardColumnModel[] = [
  {
    id: "col-1",
    boardId: "board-1",
    name: "Backlog",
    position: 0,
    createdAt: "2026-03-18T09:00:00Z",
    updatedAt: "2026-03-18T09:00:00Z",
  },
  {
    id: "col-2",
    boardId: "board-1",
    name: "To Do",
    position: 1,
    createdAt: "2026-03-18T09:00:00Z",
    updatedAt: "2026-03-18T09:00:00Z",
  },
  {
    id: "col-3",
    boardId: "board-1",
    name: "In Progress",
    position: 2,
    createdAt: "2026-03-18T09:00:00Z",
    updatedAt: "2026-03-18T09:00:00Z",
  },
  {
    id: "col-4",
    boardId: "board-1",
    name: "In Review",
    position: 3,
    createdAt: "2026-03-18T09:00:00Z",
    updatedAt: "2026-03-18T09:00:00Z",
  },
  {
    id: "col-5",
    boardId: "board-1",
    name: "Done",
    position: 4,
    createdAt: "2026-03-18T09:00:00Z",
    updatedAt: "2026-03-18T09:00:00Z",
  },
];

export const KanbanBoard = ({ boardId: _boardId }: KanbanBoardProps) => {
  // TODO: Replace mock data with actual API call
  // const [columns, error, { isLoading }] = useListBoardColumns(_boardId);
  // const columnsList = columns?.items || [];

  const columnsList = MOCK_COLUMNS;

  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: `repeat(${columnsList.length}, minmax(300px, 1fr))`,
        gap: "1rem",
        padding: "1.5rem",
        overflowX: "auto",
        overflowY: "hidden",
        width: "100%",
        height: "100%",
      }}
    >
      {columnsList.map((column) => (
        <BoardColumn key={column.id} column={column} />
      ))}
    </div>
  );
};

interface BoardColumnProps {
  column: DomainBoardColumnModel;
}

function BoardColumn({ column }: BoardColumnProps) {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        backgroundColor: "#f3f4f6",
        borderRadius: "0.5rem",
        border: "1px solid #e5e7eb",
        height: "100%",
        minWidth: "300px",
      }}
    >
      {/* Column Header */}
      <div
        style={{
          padding: "1rem",
          borderBottom: "1px solid #e5e7eb",
          backgroundColor: "#fff",
          borderRadius: "0.5rem 0.5rem 0 0",
          flexShrink: 0,
        }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <Text size="sm" style={{ fontWeight: "600", color: "#1f2937" }}>
            {column.name}
          </Text>
          <div
            style={{
              padding: "0.25rem 0.5rem",
              backgroundColor: "#e5e7eb",
              borderRadius: "0.25rem",
              fontSize: "0.75rem",
              color: "#4b5563",
              fontWeight: "500",
            }}
          >
            0
          </div>
        </div>
      </div>

      {/* Column Content */}
      <div
        style={{
          flex: 1,
          padding: "1rem",
          overflowY: "auto",
          display: "flex",
          flexDirection: "column",
          gap: "0.75rem",
        }}
      >
        {/* Tickets will go here */}
        <div
          style={{
            padding: "2rem 1rem",
            textAlign: "center",
            color: "#9ca3af",
          }}
        >
          <Text size="xs">No tickets yet</Text>
        </div>
      </div>
    </div>
  );
}
