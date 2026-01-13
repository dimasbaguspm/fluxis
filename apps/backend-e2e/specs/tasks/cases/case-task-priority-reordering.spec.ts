import { test, expect } from "../../../fixtures";

/**
 * Task Priority and Reordering Test Cases
 * Tests priority management and task reordering within statuses
 * Note: Priority can only be changed via reordering, not direct updates
 */
test.describe("Task Priority and Reordering", () => {
  const createdProjectIds: string[] = [];
  const createdStatusIds: string[] = [];
  const createdTaskIds: string[] = [];

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

  test("should assign sequential priorities when creating tasks", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test priority assignment",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks sequentially
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "First Task",
      details: "Details",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Second Task",
      details: "Details",
      priority: 2,
    });

    const task3 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Third Task",
      details: "Details",
      priority: 3,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

    // Verify priorities are sequential
    expect(task1.data?.priority).toBe(1);
    expect(task2.data?.priority).toBe(2);
    expect(task3.data?.priority).toBe(3);
  });

  test("should maintain priority uniqueness within a status", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test priority uniqueness",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create multiple tasks
    for (let i = 1; i <= 5; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Details",
        priority: i,
      });
      createdTaskIds.push(task.data!.id);
    }

    // Get all tasks for this status
    const tasks = await taskAPI.getPaginated({
      statusId: [status.id],
      sortBy: "priority",
      sortOrder: "asc",
    });

    // Verify all priorities are unique
    const priorities = tasks.data?.items?.map((t) => t.priority) || [];
    const uniquePriorities = new Set(priorities);
    expect(priorities.length).toBe(uniquePriorities.size);
  });

  test("should allow different priorities across different statuses", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test cross-status priorities",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status1 = statuses.data![0];
    const status2 = statuses.data![1];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks with same priority numbers in different statuses
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Task 1 in Status 1",
      details: "Details",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status2.id,
      title: "Task 1 in Status 2",
      details: "Details",
      priority: 1,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id);

    // Both should have priority 1 in their respective statuses
    expect(task1.data?.priority).toBe(1);
    expect(task2.data?.priority).toBe(1);
  });

  test("should sort tasks by priority correctly", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test priority sorting",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks with specific priorities
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Low Priority",
      details: "Details",
      priority: 10,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "High Priority",
      details: "Details",
      priority: 1,
    });

    const task3 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Medium Priority",
      details: "Details",
      priority: 5,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

    // Get tasks sorted by priority ascending
    const ascTasks = await taskAPI.getPaginated({
      statusId: [status.id],
      sortBy: "priority",
      sortOrder: "asc",
    });

    expect(ascTasks.data?.items![0].priority).toBe(1);
    expect(ascTasks.data?.items![1].priority).toBe(5);
    expect(ascTasks.data?.items![2].priority).toBe(10);

    // Get tasks sorted by priority descending
    const descTasks = await taskAPI.getPaginated({
      statusId: [status.id],
      sortBy: "priority",
      sortOrder: "desc",
    });

    expect(descTasks.data?.items![0].priority).toBe(10);
    expect(descTasks.data?.items![1].priority).toBe(5);
    expect(descTasks.data?.items![2].priority).toBe(1);
  });

  test("should allow updating priority directly via PATCH", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test priority update",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Test Task",
      details: "Details",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Update priority directly
    const response = await taskAPI.update(task.data!.id, {
      priority: 99,
    });

    expect(response.status).toBe(200);
    expect(response.data?.priority).toBe(99);
  });

  test("should handle priority when moving task between statuses", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test priority during status change",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status1 = statuses.data![0];
    const status2 = statuses.data![1];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks in status1
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Task 1",
      details: "Details",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Task 2",
      details: "Details",
      priority: 2,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id);

    // Get count of tasks in status2
    const status2TasksBefore = await taskAPI.getPaginated({
      statusId: [status2.id],
    });
    const countBefore = status2TasksBefore.data?.totalCount || 0;

    // Move task1 to status2
    const movedTask = await taskAPI.update(task1.data!.id, {
      statusId: status2.id,
    });

    // Verify task has a priority in new status
    expect(movedTask.data?.priority).toBeDefined();
    expect(movedTask.data?.priority).toBeGreaterThan(0);

    // Get tasks in status2 after move
    const status2TasksAfter = await taskAPI.getPaginated({
      statusId: [status2.id],
    });

    expect(status2TasksAfter.data?.totalCount).toBe(countBefore + 1);
  });

  test("should maintain priority order when deleting tasks", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test priority after deletion",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks with priorities 1, 2, 3, 4, 5
    const tasks = [];
    for (let i = 1; i <= 5; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Details",
        priority: i,
      });
      tasks.push(task.data!);
      createdTaskIds.push(task.data!.id);
    }

    // Delete middle task (priority 3)
    await taskAPI.remove(tasks[2].id);

    // Get remaining tasks
    const remainingTasks = await taskAPI.getPaginated({
      statusId: [status.id],
      sortBy: "priority",
      sortOrder: "asc",
    });

    // Should have 4 tasks left with priorities 1, 2, 4, 5
    expect(remainingTasks.data?.items?.length).toBe(4);
    const priorities = remainingTasks.data?.items?.map((t) => t.priority);
    expect(priorities).toContain(1);
    expect(priorities).toContain(2);
    expect(priorities).not.toContain(3);
    expect(priorities).toContain(4);
    expect(priorities).toContain(5);
  });

  test("should handle high priority values", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test high priority values",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create task with high priority value
    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "High Priority Task",
      details: "Details",
      priority: 999999,
    });

    createdTaskIds.push(task.data!.id);

    expect(task.status).toBe(200);
    expect(task.data?.priority).toBe(999999);
  });

  test("should list tasks with correct priority metadata", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test priority metadata",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks
    for (let i = 1; i <= 3; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Details",
        priority: i,
      });
      createdTaskIds.push(task.data!.id);
    }

    const response = await taskAPI.getPaginated({
      statusId: [status.id],
      sortBy: "priority",
      sortOrder: "asc",
    });

    // Verify each task has priority field
    response.data?.items?.forEach((task) => {
      expect(task.priority).toBeDefined();
      expect(typeof task.priority).toBe("number");
      expect(task.priority).toBeGreaterThan(0);
    });
  });

  test("should handle concurrent task creation with priorities", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Priority Test Project",
      description: "Test concurrent creation",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks concurrently
    const taskPromises = [];
    for (let i = 1; i <= 5; i++) {
      taskPromises.push(
        taskAPI.create({
          projectId: project.data!.id,
          statusId: status.id,
          title: `Concurrent Task ${i}`,
          details: "Details",
          priority: i,
        })
      );
    }

    const tasks = await Promise.all(taskPromises);

    // All should succeed
    tasks.forEach((task) => {
      expect(task.status).toBe(200);
      expect(task.data?.priority).toBeDefined();
      createdTaskIds.push(task.data!.id);
    });

    // Verify all tasks have unique priorities
    const priorities = tasks.map((t) => t.data!.priority);
    const uniquePriorities = new Set(priorities);
    expect(priorities.length).toBe(uniquePriorities.size);
  });
});
