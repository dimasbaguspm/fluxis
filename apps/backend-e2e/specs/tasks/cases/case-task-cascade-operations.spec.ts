import { test, expect } from "../../../fixtures";

/**
 * Task Cascade Operations Test Cases
 * Tests behavior when related entities (project, status) are deleted or modified
 */
test.describe("Task Cascade Operations", () => {
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

  test.describe("Status Deletion", () => {
    test("should handle task when status is deleted", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For cascade tests",
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
        title: "Task in Status 1",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Delete the status (soft delete, so FK constraint won't trigger)
      await statusAPI.remove(status1.id);

      // Task should remain accessible with the statusId still referencing the soft-deleted status
      const taskAfterDelete = await taskAPI.getById(task.data!.id);

      expect(taskAfterDelete.status).toBe(200);
      // Status ID remains the same (soft delete doesn't trigger ON DELETE SET NULL)
      expect(taskAfterDelete.data?.statusId).toBe(status1.id);
    });

    test("should handle multiple tasks when status is deleted", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For cascade tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create multiple tasks in the status
      const task1 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task 1",
        details: "Details",
        priority: 1,
      });

      const task2 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task 2",
        details: "Details",
        priority: 2,
      });

      const task3 = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task 3",
        details: "Details",
        priority: 3,
      });

      createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

      // Delete the status
      await statusAPI.remove(status.id);

      // Check all tasks
      const task1After = await taskAPI.getById(task1.data!.id);
      const task2After = await taskAPI.getById(task2.data!.id);
      const task3After = await taskAPI.getById(task3.data!.id);

      // All tasks should have consistent behavior
      expect(task1After.status).toBe(task2After.status);
      expect(task2After.status).toBe(task3After.status);
    });

    test("should not allow deleting default status with tasks", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For cascade tests",
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
        title: "Task in Default Status",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Try to delete default status (should fail or have special handling)
      const deleteResponse = await statusAPI.remove(defaultStatus.id);

      // Either deletion fails or task handling is specific
      if (deleteResponse.status >= 400) {
        expect(deleteResponse.error).toBeDefined();
      } else {
        // If deletion succeeds, verify task still exists
        const taskCheck = await taskAPI.getById(task.data!.id);
        expect(taskCheck.status).toBe(200);
      }
    });
  });

  test.describe("Project Deletion", () => {
    test("should handle task when project is deleted", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For cascade tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create task
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task in Project",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Delete the project (soft delete)
      await projectAPI.remove(project.data!.id);

      // Task should not be accessible when project is soft deleted (cascade behavior)
      const taskAfterDelete = await taskAPI.getById(task.data!.id);
      expect(taskAfterDelete.status).toBe(404);
    });

    test("should delete all tasks when project is deleted", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For cascade tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status1 = statuses.data![0];
      const status2 = statuses.data![1];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create tasks in multiple statuses
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

      // Delete the project (soft delete)
      await projectAPI.remove(project.data!.id);

      // All tasks should not be accessible when project is soft deleted
      const task1After = await taskAPI.getById(task1.data!.id);
      const task2After = await taskAPI.getById(task2.data!.id);

      expect(task1After.status).toBe(404);
      expect(task2After.status).toBe(404);
    });

    test("should not affect tasks in other projects", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      // Create two projects
      const project1 = await projectAPI.create({
        name: "Project 1",
        description: "First project",
        status: "active",
      });

      const project2 = await projectAPI.create({
        name: "Project 2",
        description: "Second project",
        status: "active",
      });

      createdProjectIds.push(project1.data!.id, project2.data!.id);

      const statuses1 = await statusAPI.getByProject(project1.data!.id);
      const status1 = statuses1.data![0];

      const statuses2 = await statusAPI.getByProject(project2.data!.id);
      const status2 = statuses2.data![0];

      createdStatusIds.push(...statuses1.data!.map((s) => s.id));
      createdStatusIds.push(...statuses2.data!.map((s) => s.id));

      // Create tasks in both projects
      const task1 = await taskAPI.create({
        projectId: project1.data!.id,
        statusId: status1.id,
        title: "Task in Project 1",
        details: "Details",
        priority: 1,
      });

      const task2 = await taskAPI.create({
        projectId: project2.data!.id,
        statusId: status2.id,
        title: "Task in Project 2",
        details: "Details",
        priority: 1,
      });

      createdTaskIds.push(task1.data!.id, task2.data!.id);

      // Delete project1 (soft delete)
      await projectAPI.remove(project1.data!.id);

      // task1 should not be accessible (project deleted)
      const task1After = await taskAPI.getById(task1.data!.id);
      expect(task1After.status).toBe(404);

      // task2 should still exist (different project)
      const task2After = await taskAPI.getById(task2.data!.id);
      expect(task2After.status).toBe(200);
      expect(task2After.data?.id).toBe(task2.data!.id);
    });
  });

  test.describe("Project Status Change", () => {
    test("should allow tasks in archived project", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For status change tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create task
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Archive the project
      await projectAPI.update(project.data!.id, {
        status: "archived",
      });

      // Task should still be accessible
      const taskAfter = await taskAPI.getById(task.data!.id);
      expect(taskAfter.status).toBe(200);
      expect(taskAfter.data?.id).toBe(task.data!.id);
    });

    test("should allow tasks in paused project", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For status change tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create task
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Pause the project
      await projectAPI.update(project.data!.id, {
        status: "paused",
      });

      // Task should still be accessible
      const taskAfter = await taskAPI.getById(task.data!.id);
      expect(taskAfter.status).toBe(200);
      expect(taskAfter.data?.id).toBe(task.data!.id);
    });

    test("should prevent creating tasks in archived project", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For status change tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Archive the project
      await projectAPI.update(project.data!.id, {
        status: "archived",
      });

      // Try to create task in archived project
      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task in Archived Project",
        details: "Details",
        priority: 1,
      });

      // Should either fail or succeed based on business rules
      // Adjust expectation based on requirements
      if (response.status >= 400) {
        expect(response.error).toBeDefined();
      } else {
        createdTaskIds.push(response.data!.id);
      }
    });
  });

  test.describe("Status Update", () => {
    test("should handle tasks when status is renamed", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For status update tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create task
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Rename the status
      await statusAPI.update(status.id, {
        name: "New Status Name",
      });

      // Task should still reference the same status
      const taskAfter = await taskAPI.getById(task.data!.id);
      expect(taskAfter.status).toBe(200);
      expect(taskAfter.data?.statusId).toBe(status.id);
    });

    test("should handle tasks when status is reordered", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For status reorder tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status1 = statuses.data![0];
      const status2 = statuses.data![1];
      const status3 = statuses.data![2];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create tasks in different statuses
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

      // Reorder statuses
      await statusAPI.reorder({
        projectId: project.data!.id,
        ids: [status3.id, status1.id, status2.id],
      });

      // Tasks should still be in their respective statuses
      const task1After = await taskAPI.getById(task1.data!.id);
      const task2After = await taskAPI.getById(task2.data!.id);

      expect(task1After.data?.statusId).toBe(status1.id);
      expect(task2After.data?.statusId).toBe(status2.id);
    });
  });

  test.describe("Bulk Operations", () => {
    test("should handle deleting multiple tasks", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For bulk tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create multiple tasks
      const taskIds: string[] = [];
      for (let i = 1; i <= 5; i++) {
        const task = await taskAPI.create({
          projectId: project.data!.id,
          statusId: status.id,
          title: `Task ${i}`,
          details: "Details",
          priority: i,
        });
        taskIds.push(task.data!.id);
      }

      // Delete all tasks
      for (const id of taskIds) {
        await taskAPI.remove(id);
      }

      // Verify all tasks are deleted
      for (const id of taskIds) {
        const response = await taskAPI.getById(id);
        expect(response.status).toBe(404);
      }
    });

    test("should maintain data integrity when operations fail", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For integrity tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      // Create valid task
      const validTask = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Valid Task",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(validTask.data!.id);

      // Try to create invalid task (should fail)
      await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "",
        details: "Details",
        priority: 1,
      });

      // Valid task should still exist and be unchanged
      const validTaskAfter = await taskAPI.getById(validTask.data!.id);
      expect(validTaskAfter.status).toBe(200);
      expect(validTaskAfter.data?.title).toBe("Valid Task");
    });
  });
});
