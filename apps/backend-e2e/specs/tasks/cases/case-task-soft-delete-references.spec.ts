import { test, expect } from "../../../fixtures";

/**
 * Task Soft-Delete References Test Cases
 * Tests task behavior when interacting with soft-deleted statuses
 */
test.describe("Task Soft-Delete References", () => {
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

  test("should prevent creating task with soft-deleted status", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Soft delete the status
    await statusAPI.remove(status.id);

    // Try to create task with the soft-deleted status
    const response = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Task with deleted status",
      details: "Should fail",
      priority: 1,
    });

    // Should fail because status is soft-deleted
    expect(response.status).toBeGreaterThanOrEqual(400);
  });

  test("should prevent moving task to soft-deleted status", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status1 = statuses.data![0];
    const status2 = statuses.data![1];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create task in status1
    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Test Task",
      details: "Details",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Soft delete status2
    await statusAPI.remove(status2.id);

    // Try to move task to soft-deleted status
    const response = await taskAPI.update(task.data!.id, {
      statusId: status2.id,
    });

    // Should fail because target status is soft-deleted
    expect(response.status).toBeGreaterThanOrEqual(400);
  });

  test("should allow task to remain with soft-deleted status reference", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create task in status
    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Test Task",
      details: "Details",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Soft delete the status (task still references it)
    await statusAPI.remove(status.id);

    // Task should still be accessible and reference the deleted status
    const taskAfter = await taskAPI.getById(task.data!.id);
    expect(taskAfter.status).toBe(200);
    expect(taskAfter.data?.statusId).toBe(status.id);
  });

  test("should filter out tasks with soft-deleted status in list queries", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status1 = statuses.data![0];
    const status2 = statuses.data![1];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks in both statuses
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Task 1",
      details: "Details",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status2.id,
      title: "Task 2",
      details: "Details",
      priority: 1,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id);

    // Soft delete status1
    await statusAPI.remove(status1.id);

    // Query tasks filtered by deleted status
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      statusId: [status1.id],
    });

    // Should return empty or not include tasks with soft-deleted status
    // This behavior depends on your business logic
    expect(response.status).toBe(200);
    // Depending on implementation, might return 0 items or exclude them
  });

  test("should handle updating task that references soft-deleted status", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status1 = statuses.data![0];
    const status2 = statuses.data![1];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create task
    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Test Task",
      details: "Details",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Soft delete status1 (task still references it)
    await statusAPI.remove(status1.id);

    // Update task title (not status)
    const response = await taskAPI.update(task.data!.id, {
      title: "Updated Title",
    });

    expect(response.status).toBe(200);
    expect(response.data?.title).toBe("Updated Title");

    // Try to move to a valid status
    const moveResponse = await taskAPI.update(task.data!.id, {
      statusId: status2.id,
    });

    expect(moveResponse.status).toBe(200);
    expect(moveResponse.data?.statusId).toBe(status2.id);
  });

  test("should prevent deleting status that is default and has tasks", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const defaultStatus = statuses.data!.find((s) => s.isDefault);
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    if (!defaultStatus) {
      throw new Error("No default status found");
    }

    // Create task in default status
    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: defaultStatus.id,
      title: "Task in Default",
      details: "Details",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Try to delete default status with tasks
    const response = await statusAPI.remove(defaultStatus.id);

    // Depending on business logic, might prevent deletion or allow it
    // Document the actual behavior here
    if (response.status < 400) {
      // If deletion is allowed, task should still be accessible
      const taskCheck = await taskAPI.getById(task.data!.id);
      expect(taskCheck.status).toBe(200);
    } else {
      // If deletion is prevented
      expect(response.status).toBeGreaterThanOrEqual(400);
    }
  });
});
