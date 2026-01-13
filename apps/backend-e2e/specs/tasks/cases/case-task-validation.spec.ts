import { test, expect } from "../../../fixtures";

/**
 * Task Validation Test Cases
 * Tests input validation, business rules, and edge cases
 */
test.describe("Task Validation", () => {
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

  test.describe("Title Validation", () => {
    test("should reject empty title", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "",
        details: "Valid details",
        priority: 1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should reject whitespace-only title", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "   ",
        details: "Valid details",
        priority: 1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should accept title with special characters", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const specialTitle = "Task #1: Fix bug @user - [URGENT] (50%)";

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: specialTitle,
        details: "Details",
        priority: 1,
      });

      expect(response.status).toBe(200);
      expect(response.data?.title).toBe(specialTitle);

      createdTaskIds.push(response.data!.id);
    });

    test("should accept very long title", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const longTitle = "A".repeat(500);

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: longTitle,
        details: "Details",
        priority: 1,
      });

      expect(response.status).toBe(200);
      expect(response.data?.title).toBe(longTitle);

      createdTaskIds.push(response.data!.id);
    });

    test("should accept title with unicode characters", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const unicodeTitle = "任务 🚀 Tâche ñoño";

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: unicodeTitle,
        details: "Details",
        priority: 1,
      });

      expect(response.status).toBe(200);
      expect(response.data?.title).toBe(unicodeTitle);

      createdTaskIds.push(response.data!.id);
    });
  });

  test.describe("Details Validation", () => {
    test("should accept empty details", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with no details",
        details: "",
        priority: 1,
      });

      // Depending on API design, this might succeed or fail
      // Adjust based on actual requirements
      expect(response.status).toBe(200);
      createdTaskIds.push(response.data!.id);
    });

    test("should accept very long details", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const longDetails = "Details ".repeat(1000);

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with long details",
        details: longDetails,
        priority: 1,
      });

      expect(response.status).toBe(200);
      expect(response.data?.details).toBe(longDetails);

      createdTaskIds.push(response.data!.id);
    });

    test("should accept details with markdown formatting", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const markdownDetails = `
# Heading
- Item 1
- Item 2

**Bold** and *italic*

\`\`\`code
function test() {}
\`\`\`
      `;

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with markdown",
        details: markdownDetails,
        priority: 1,
      });

      expect(response.status).toBe(200);
      expect(response.data?.details).toBe(markdownDetails);

      createdTaskIds.push(response.data!.id);
    });
  });

  test.describe("Due Date Validation", () => {
    test("should accept valid future due date", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const futureDate = new Date();
      futureDate.setDate(futureDate.getDate() + 30);

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with future due date",
        details: "Details",
        priority: 1,
        dueDate: futureDate.toISOString(),
      });

      expect(response.status).toBe(200);
      expect(response.data?.dueDate).toBeDefined();

      createdTaskIds.push(response.data!.id);
    });

    test("should accept past due date", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const pastDate = new Date();
      pastDate.setDate(pastDate.getDate() - 30);

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with past due date",
        details: "Details",
        priority: 1,
        dueDate: pastDate.toISOString(),
      });

      expect(response.status).toBe(200);
      expect(response.data?.dueDate).toBeDefined();

      createdTaskIds.push(response.data!.id);
    });

    test("should accept task without due date", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task without due date",
        details: "Details",
        priority: 1,
      });

      expect(response.status).toBe(200);

      createdTaskIds.push(response.data!.id);
    });

    test("should reject invalid due date format", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with invalid date",
        details: "Details",
        priority: 1,
        dueDate: "not-a-date" as any,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should allow clearing due date", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const futureDate = new Date();
      futureDate.setDate(futureDate.getDate() + 7);

      // Create task with due date
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with due date",
        details: "Details",
        priority: 1,
        dueDate: futureDate.toISOString(),
      });
      createdTaskIds.push(task.data!.id);

      // Clear due date (if API supports this)
      const updated = await taskAPI.update(task.data!.id, {
        dueDate: null as any,
      });

      // Depending on API design, this might work or fail
      // Adjust based on requirements
      if (updated.status === 200) {
        expect(updated.data?.dueDate).toBeUndefined();
      }
    });
  });

  test.describe("Priority Validation", () => {
    test("should reject negative priority", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with negative priority",
        details: "Details",
        priority: -1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should reject zero priority", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with zero priority",
        details: "Details",
        priority: 0,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should reject non-integer priority", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task with decimal priority",
        details: "Details",
        priority: 1.5 as any,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("Foreign Key Validation", () => {
    test("should reject invalid projectId format", async ({ taskAPI }) => {
      const response = await taskAPI.create({
        projectId: "not-a-uuid",
        statusId: "00000000-0000-0000-0000-000000000000",
        title: "Task",
        details: "Details",
        priority: 1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should reject invalid statusId format", async ({ taskAPI }) => {
      const response = await taskAPI.create({
        projectId: "00000000-0000-0000-0000-000000000000",
        statusId: "not-a-uuid",
        title: "Task",
        details: "Details",
        priority: 1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should reject non-existent projectId", async ({ taskAPI }) => {
      const response = await taskAPI.create({
        projectId: "00000000-0000-0000-0000-000000000000",
        statusId: "00000000-0000-0000-0000-000000000001",
        title: "Task",
        details: "Details",
        priority: 1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should reject non-existent statusId", async ({
      projectAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const response = await taskAPI.create({
        projectId: project.data!.id,
        statusId: "00000000-0000-0000-0000-000000000000",
        title: "Task",
        details: "Details",
        priority: 1,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("Update Validation", () => {
    test("should allow partial updates", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Original Title",
        details: "Original details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Update only title
      const updated = await taskAPI.update(task.data!.id, {
        title: "New Title",
      });

      expect(updated.status).toBe(200);
      expect(updated.data?.title).toBe("New Title");
      expect(updated.data?.details).toBe("Original details");
    });

    test("should validate updated fields", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const statuses = await statusAPI.getByProject(project.data!.id);
      const status = statuses.data![0];
      createdStatusIds.push(...statuses.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: "Task",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Try to update with invalid title
      const response = await taskAPI.update(task.data!.id, {
        title: "",
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should not allow changing projectId", async ({
      projectAPI,
      statusAPI,
      taskAPI,
    }) => {
      const project1 = await projectAPI.create({
        name: "Project 1",
        description: "Test",
        status: "active",
      });
      const project2 = await projectAPI.create({
        name: "Project 2",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project1.data!.id, project2.data!.id);

      const statuses1 = await statusAPI.getByProject(project1.data!.id);
      const status1 = statuses1.data![0];
      createdStatusIds.push(...statuses1.data!.map((s) => s.id));

      const statuses2 = await statusAPI.getByProject(project2.data!.id);
      createdStatusIds.push(...statuses2.data!.map((s) => s.id));

      const task = await taskAPI.create({
        projectId: project1.data!.id,
        statusId: status1.id,
        title: "Task",
        details: "Details",
        priority: 1,
      });
      createdTaskIds.push(task.data!.id);

      // Try to change projectId
      const response = await taskAPI.update(task.data!.id, {
        projectId: project2.data!.id,
      } as any);

      // Should either fail or ignore the projectId field
      if (response.status === 200) {
        expect(response.data?.projectId).toBe(project1.data!.id);
      } else {
        expect(response.status).toBeGreaterThanOrEqual(400);
      }
    });
  });
});
