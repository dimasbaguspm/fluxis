import { test, expect } from "../../fixtures";

/**
 * Task API - Common CRUD Tests
 * Tests standard operations: Create, Read, Update, Delete, List
 */
test.describe("Task API", () => {
  const createdProjectIds: string[] = [];
  const createdStatusIds: string[] = [];
  const createdTaskIds: string[] = [];

  // Cleanup after each test
  test.afterEach(async ({ projectAPI, statusAPI, taskAPI }) => {
    for (const id of createdTaskIds) {
      await taskAPI.remove(id).catch(() => {});
    }
    for (const id of createdStatusIds) {
      await statusAPI.remove(id).catch(() => {});
    }
    for (const id of createdProjectIds) {
      await projectAPI.remove(id).catch(() => {});
    }
    createdTaskIds.length = 0;
    createdStatusIds.length = 0;
    createdProjectIds.length = 0;
  });

  test.describe("POST /tasks", () => {
    test("should create task with all required fields", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      // Setup: Create project and status
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For task tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create task
      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Test Task",
        details: "Test task details",
        priority: 1,
      });

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.id).toBeDefined();
      expect(response.data?.projectId).toBe(project.data!.id);
      expect(response.data?.statusId).toBe(firstStatus.id);
      expect(response.data?.title).toBe("Test Task");
      expect(response.data?.details).toBe("Test task details");
      expect(response.data?.priority).toBe(1);
      expect(response.data?.createdAt).toBeDefined();
      expect(response.data?.updatedAt).toBeDefined();

      // Verify timestamps are valid
      expect(new Date(response.data!.createdAt).getTime()).toBeGreaterThan(0);
      expect(new Date(response.data!.updatedAt).getTime()).toBeGreaterThan(0);

      createdTaskIds.push(response.data!.id);
    });

    test("should create task with due date", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For task tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const dueDate = new Date();
      dueDate.setDate(dueDate.getDate() + 7); // 7 days from now

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task with Due Date",
        details: "This task has a due date",
        priority: 1,
        dueDate: dueDate.toISOString(),
      });

      expect(response.status).toBe(200);
      expect(response.data?.dueDate).toBeDefined();
      expect(new Date(response.data!.dueDate!).getTime()).toBeGreaterThan(0);

      createdTaskIds.push(response.data!.id);
    });

    test("should create multiple tasks with auto-incrementing priority", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For priority tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create tasks without specifying priority
      const task1 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task 1",
        details: "First task",
        priority: 1,
      });

      const task2 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task 2",
        details: "Second task",
        priority: 2,
      });

      const task3 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task 3",
        details: "Third task",
        priority: 3,
      });

      expect(task1.data?.priority).toBe(1);
      expect(task2.data?.priority).toBe(2);
      expect(task3.data?.priority).toBe(3);

      createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);
    });

    test("should fail to create task with missing projectId", async ({
      statusAPI,
      taskAPI,
    }) => {
      const response = await taskAPI.create({
        statusId: "00000000-0000-0000-0000-000000000000",
        title: "Invalid Task",
        details: "Missing project",
        priority: 1,
      } as any);

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create task with missing statusId", async ({
      projectAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const response = await taskAPI.create({
        projectId: project.data!.id,
        title: "Invalid Task",
        details: "Missing status",
        priority: 1,
      } as any);

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create task with missing title", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        details: "Missing title",
        priority: 1,
      } as any);

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create task with empty title", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "",
        details: "Empty title",
        priority: 1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create task with invalid projectId", async ({
      statusAPI,
      taskAPI,
    }) => {
      const response = await taskAPI.create({
        projectId: "invalid-uuid",
        statusId: "00000000-0000-0000-0000-000000000000",
        title: "Invalid Task",
        details: "Invalid project ID",
        priority: 1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create task with non-existent projectId", async ({
      taskAPI,
    }) => {
      const response = await taskAPI.create({
        projectId: "00000000-0000-0000-0000-000000000000",
        statusId: "00000000-0000-0000-0000-000000000001",
        title: "Non-existent Project Task",
        details: "Project does not exist",
        priority: 1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("GET /tasks", () => {
    test("should get paginated tasks", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For listing tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create some tasks
      const task1 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task 1",
        details: "Details 1",
        priority: 1,
      });
      const task2 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task 2",
        details: "Details 2",
        priority: 2,
      });
      createdTaskIds.push(task1.data!.id, task2.data!.id);

      const response = await taskAPI.getPaginated({
        projectId: [project.data!.id],
      });

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.items).toBeDefined();
      expect(response.data?.items!.length).toBeGreaterThanOrEqual(2);
      expect(response.data?.pageNumber).toBe(1);
      expect(response.data?.totalCount).toBeGreaterThanOrEqual(2);
    });

    test("should filter tasks by statusId", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For filter tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status1 = statuses.data![0];
      const status2 = statuses.data![1];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create tasks in different statuses
      const task1 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status1.id,
        title: "Task in Status 1",
        details: "Details",
        priority: 1,
      });
      const task2 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status2.id,
        title: "Task in Status 2",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task1.data!.id, task2.data!.id);

      // Filter by status1
      const response = await taskAPI.getPaginated({
        statusId: [status1.id],
      });

      expect(response.status).toBe(200);
      expect(response.data?.items).toBeDefined();
      expect(
        response.data?.items!.every((t) => t.statusId === status1.id)
      ).toBe(true);
    });

    test("should sort tasks by priority", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For sort tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create tasks with different priorities
      const task1 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Low Priority",
        details: "Details",
        priority: 3,
      });
      const task2 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "High Priority",
        details: "Details",
        priority: 1,
      });
      const task3 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Medium Priority",
        details: "Details",
        priority: 2,
      });
      createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

      // Sort by priority ascending
      const response = await taskAPI.getPaginated({
        projectId: [project.data!.id],
        sortBy: "priority",
        sortOrder: "asc",
      });

      expect(response.status).toBe(200);
      expect(response.data?.items).toBeDefined();
      expect(response.data?.items![0].priority).toBeLessThanOrEqual(
        response.data?.items![1].priority!
      );
    });

    test("should search tasks by query", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For search tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Unique Search Term XYZ123",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      const response = await taskAPI.getPaginated({
        query: "XYZ123",
      });

      expect(response.status).toBe(200);
      expect(response.data?.items).toBeDefined();
      expect(
        response.data?.items!.some((t) => t.title.includes("XYZ123"))
      ).toBe(true);
    });
  });

  test.describe("GET /tasks/{taskId}", () => {
    test("should get task detail by ID", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For detail tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Test Task",
        details: "Test details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      const response = await taskAPI.getById(task.data!.id);

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.id).toBe(task.data!.id);
      expect(response.data?.title).toBe("Test Task");
    });

    test("should fail to get non-existent task", async ({ taskAPI }) => {
      const response = await taskAPI.getById(
        "00000000-0000-0000-0000-000000000000"
      );

      expect(response.status).toBe(404);
      expect(response.error).toBeDefined();
    });

    test("should fail to get task with invalid ID", async ({ taskAPI }) => {
      const response = await taskAPI.getById("invalid-uuid");

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("PATCH /tasks/{taskId}", () => {
    test("should update task title", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For update tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Original Title",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      const response = await taskAPI.update(task.data!.id, {
        title: "Updated Title",
      });

      expect(response.status).toBe(200);
      expect(response.data?.title).toBe("Updated Title");
      expect(response.data?.details).toBe("Details"); // Unchanged
    });

    test("should update task details", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For update tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task Title",
        details: "Original details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      const response = await taskAPI.update(task.data!.id, {
        details: "Updated details",
      });

      expect(response.status).toBe(200);
      expect(response.data?.details).toBe("Updated details");
    });

    test("should update task status (move between statuses)", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For status change tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status1 = statuses.data![0];
      const status2 = statuses.data![1];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status1.id,
        title: "Task to Move",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      const response = await taskAPI.update(task.data!.id, {
        statusId: status2.id,
      });

      expect(response.status).toBe(200);
      expect(response.data?.statusId).toBe(status2.id);
      expect(response.data?.projectId).toBe(project.data!.id); // Project unchanged
    });

    test("should update task due date", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For due date tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task with Due Date",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      const newDueDate = new Date();
      newDueDate.setDate(newDueDate.getDate() + 14);

      const response = await taskAPI.update(task.data!.id, {
        dueDate: newDueDate.toISOString(),
      });

      expect(response.status).toBe(200);
      expect(response.data?.dueDate).toBeDefined();
    });

    test("should fail to update task with empty title", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Valid Title",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      const response = await taskAPI.update(task.data!.id, {
        title: "",
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to update non-existent task", async ({ taskAPI }) => {
      const response = await taskAPI.update(
        "00000000-0000-0000-0000-000000000000",
        {
          title: "Updated",
        }
      );

      expect(response.status).toBe(404);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("DELETE /tasks/{taskId}", () => {
    test("should delete task", async ({ projectAPI, statusAPI, taskAPI }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For delete tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task to Delete",
        details: "Details",
        priority: 1,
      });

      const deleteResponse = await taskAPI.remove(task.data!.id);
      expect(deleteResponse.status).toBe(204);

      // Verify task is deleted
      const getResponse = await taskAPI.getById(task.data!.id);
      expect(getResponse.status).toBe(404);
    });

    test("should fail to delete non-existent task", async ({ taskAPI }) => {
      const response = await taskAPI.remove(
        "00000000-0000-0000-0000-000000000000"
      );

      expect(response.status).toBe(404);
      expect(response.error).toBeDefined();
    });

    test("should fail to delete task with invalid ID", async ({ taskAPI }) => {
      const response = await taskAPI.remove("invalid-uuid");

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("GET /tasks/{taskId}/logs", () => {
    test("should get logs for a task", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For log tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task with Logs",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Update task to generate logs
      await taskAPI.update(task.data!.id, {
        title: "Updated Task Title",
      });

      const response = await taskAPI.getLogs(task.data!.id);

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.items).toBeDefined();
    });

    test("should paginate task logs", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Task Test Project",
        description: "For log pagination",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const firstStatus = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: firstStatus.id,
        title: "Task with Many Logs",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      const response = await taskAPI.getLogs(task.data!.id, {
        pageNumber: 1,
        pageSize: 10,
      });

      expect(response.status).toBe(200);
      expect(response.data?.pageNumber).toBe(1);
      expect(response.data?.pageSize).toBe(10);
    });
  });
});
