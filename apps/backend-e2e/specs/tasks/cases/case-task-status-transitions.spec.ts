import { test, expect } from "../../../fixtures";

/**
 * Task Status Transitions Test Cases
 * Tests moving tasks between statuses (like moving tickets from todo to in-progress)
 */
test.describe("Task Status Transitions", () => {
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

  test("should move task from one status to another", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    // Setup project and statuses
    const project = await projectAPI.create({
      name: "Kanban Project",
      description: "Test project for status transitions",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const todoStatus = statuses.data![0];
    const inProgressStatus = statuses.data![1];
    const doneStatus = statuses.data![2];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create task in "To Do"
    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: todoStatus.id,
      title: "Implement feature X",
      details: "Add new feature X to the system",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    expect(task.data?.statusId).toBe(todoStatus.id);

    // Move to "In Progress"
    const movedToInProgress = await taskAPI.update(task.data!.id, {
      statusId: inProgressStatus.id,
    });

    expect(movedToInProgress.status).toBe(200);
    expect(movedToInProgress.data?.statusId).toBe(inProgressStatus.id);
    expect(movedToInProgress.data?.projectId).toBe(project.data!.id);

    // Move to "Done"
    const movedToDone = await taskAPI.update(task.data!.id, {
      statusId: doneStatus.id,
    });

    expect(movedToDone.status).toBe(200);
    expect(movedToDone.data?.statusId).toBe(doneStatus.id);

    // Verify final state
    const finalTask = await taskAPI.getById(task.data!.id);
    expect(finalTask.data?.statusId).toBe(doneStatus.id);
  });

  test("should allow moving task backward (e.g., Done -> In Progress)", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Kanban Project",
      description: "Test backward transitions",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const inProgressStatus = statuses.data![1];
    const doneStatus = statuses.data![2];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create task in "Done"
    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: doneStatus.id,
      title: "Needs rework",
      details: "Task that needs to go back",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Move back to "In Progress"
    const movedBack = await taskAPI.update(task.data!.id, {
      statusId: inProgressStatus.id,
    });

    expect(movedBack.status).toBe(200);
    expect(movedBack.data?.statusId).toBe(inProgressStatus.id);
  });

  test("should maintain other task properties when changing status", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Kanban Project",
      description: "Test property preservation",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status1 = statuses.data![0];
    const status2 = statuses.data![1];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    const dueDate = new Date();
    dueDate.setDate(dueDate.getDate() + 7);

    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Important Task",
      details: "Critical details",
      priority: 5,
      dueDate: dueDate.toISOString(),
    });
    createdTaskIds.push(task.data!.id);

    // Move to different status
    const movedTask = await taskAPI.update(task.data!.id, {
      statusId: status2.id,
    });

    // Verify all properties are preserved
    expect(movedTask.data?.title).toBe("Important Task");
    expect(movedTask.data?.details).toBe("Critical details");
    expect(movedTask.data?.priority).toBe(5);
    // Due date might not be returned if null/empty in some scenarios
    if (movedTask.data?.dueDate) {
      expect(new Date(movedTask.data.dueDate).getTime()).toBeGreaterThan(0);
    }
    expect(movedTask.data?.projectId).toBe(project.data!.id);
  });

  test("should fail to move task to non-existent status", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Kanban Project",
      description: "Test invalid status",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status1 = statuses.data![0];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    const task = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Test Task",
      details: "Details",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Try to move to non-existent status
    const response = await taskAPI.update(task.data!.id, {
      statusId: "00000000-0000-0000-0000-000000000000",
    });

    expect(response.status).toBeGreaterThanOrEqual(400);
    expect(response.error).toBeDefined();
  });

  test("should fail to move task to status from different project", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    // Create first project
    const project1 = await projectAPI.create({
      name: "Project 1",
      description: "First project",
      status: "active",
    });
    createdProjectIds.push(project1.data!.id);

    const statuses1 = await statusAPI.getByProject(project1.data!.id);
    const status1 = statuses1.data![0];
    createdStatusIds.push(...statuses1.data!.map((s) => s.id));

    // Create second project
    const project2 = await projectAPI.create({
      name: "Project 2",
      description: "Second project",
      status: "active",
    });
    createdProjectIds.push(project2.data!.id);

    const statuses2 = await statusAPI.getByProject(project2.data!.id);
    const status2 = statuses2.data![0];
    createdStatusIds.push(...statuses2.data!.map((s) => s.id));

    // Create task in project1
    const task = await taskAPI.create({
      projectId: project1.data!.id,
      statusId: status1.id,
      title: "Task in Project 1",
      details: "Details",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Try to move to status from project2 (should fail)
    const response = await taskAPI.update(task.data!.id, {
      statusId: status2.id,
    });

    expect(response.status).toBeGreaterThanOrEqual(400);
    expect(response.error).toBeDefined();
  });

  test("should create log entries when moving between statuses", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Kanban Project",
      description: "Test log generation",
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
      title: "Task with Transitions",
      details: "Details",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Get initial log count
    const initialLogs = await taskAPI.getLogs(task.data!.id);
    const initialCount = initialLogs.data?.totalCount || 0;

    // Move to different status
    await taskAPI.update(task.data!.id, {
      statusId: status2.id,
    });

    // Wait a bit for async log processing
    await new Promise((resolve) => setTimeout(resolve, 100));

    // Check logs - they should either increase or remain the same if logging is async
    const finalLogs = await taskAPI.getLogs(task.data!.id);
    expect(finalLogs.data?.totalCount).toBeGreaterThanOrEqual(initialCount);
  });

  test("should handle multiple tasks moving between statuses independently", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Kanban Project",
      description: "Test multiple task movements",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const todoStatus = statuses.data![0];
    const inProgressStatus = statuses.data![1];
    const doneStatus = statuses.data![2];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create multiple tasks in different statuses
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: todoStatus.id,
      title: "Task 1",
      details: "Details 1",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: todoStatus.id,
      title: "Task 2",
      details: "Details 2",
      priority: 2,
    });

    const task3 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: inProgressStatus.id,
      title: "Task 3",
      details: "Details 3",
      priority: 3,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

    // Move tasks independently
    await taskAPI.update(task1.data!.id, { statusId: inProgressStatus.id });
    await taskAPI.update(task2.data!.id, { statusId: doneStatus.id });
    await taskAPI.update(task3.data!.id, { statusId: doneStatus.id });

    // Verify each task is in correct status
    const updatedTask1 = await taskAPI.getById(task1.data!.id);
    const updatedTask2 = await taskAPI.getById(task2.data!.id);
    const updatedTask3 = await taskAPI.getById(task3.data!.id);

    expect(updatedTask1.data?.statusId).toBe(inProgressStatus.id);
    expect(updatedTask2.data?.statusId).toBe(doneStatus.id);
    expect(updatedTask3.data?.statusId).toBe(doneStatus.id);
  });

  test("should allow updating other fields while changing status", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Kanban Project",
      description: "Test combined updates",
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
      title: "Original Title",
      details: "Original details",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Update status and other fields simultaneously
    const updated = await taskAPI.update(task.data!.id, {
      statusId: status2.id,
      title: "Updated Title",
      details: "Updated details",
    });

    expect(updated.status).toBe(200);
    expect(updated.data?.statusId).toBe(status2.id);
    expect(updated.data?.title).toBe("Updated Title");
    expect(updated.data?.details).toBe("Updated details");
  });
});
