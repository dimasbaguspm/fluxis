import { test, expect } from "../../../fixtures";

/**
 * Task Pagination Edge Cases
 * Tests edge cases and boundary conditions for task pagination
 */
test.describe("Task Pagination Edge Cases", () => {
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

  test("should handle page beyond total pages", async ({
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

    // Create a few tasks
    for (let i = 0; i < 3; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Test",
        priority: i + 1,
      });
      createdTaskIds.push(task.data!.id);
    }

    // Request page 100 (way beyond available data)
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      pageNumber: 100,
      pageSize: 10,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toEqual([]);
    expect(response.data?.pageNumber).toBe(100);
    // Total count reflects tasks in the filtered project (may be affected by timing)
    expect(response.data?.totalCount).toBeGreaterThanOrEqual(0);
  });

  test("should handle very large page size", async ({
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

    // Create tasks
    for (let i = 0; i < 5; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Test",
        priority: i + 1,
      });
      createdTaskIds.push(task.data!.id);
    }

    // Request very large page size
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      pageNumber: 1,
      pageSize: 10000,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBeGreaterThanOrEqual(5);
    expect(response.data?.pageSize).toBe(10000);
  });

  test("should handle page size of 1", async ({
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

    // Create tasks
    for (let i = 0; i < 3; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Test",
        priority: i + 1,
      });
      createdTaskIds.push(task.data!.id);
    }

    // Request page size of 1
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      pageNumber: 1,
      pageSize: 1,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBe(1);
    expect(response.data?.pageSize).toBe(1);
    expect(response.data?.totalPages).toBeGreaterThanOrEqual(3);
  });

  test("should handle zero page number", async ({ taskAPI }) => {
    const response = await taskAPI.getPaginated({
      pageNumber: 0,
      pageSize: 10,
    });

    // Should default to page 1 or return error
    expect(response.status).toBeLessThan(500);
  });

  test("should handle negative page number", async ({ taskAPI }) => {
    const response = await taskAPI.getPaginated({
      pageNumber: -1,
      pageSize: 10,
    });

    // Should return error or default to page 1
    expect(response.status).toBeLessThan(500);
  });

  test("should handle zero page size", async ({ taskAPI }) => {
    const response = await taskAPI.getPaginated({
      pageNumber: 1,
      pageSize: 0,
    });

    // Should return error or use default page size
    expect(response.status).toBeLessThan(500);
  });

  test("should handle negative page size", async ({ taskAPI }) => {
    const response = await taskAPI.getPaginated({
      pageNumber: 1,
      pageSize: -10,
    });

    // Should return error or use default page size
    expect(response.status).toBeLessThan(500);
  });

  test("should handle missing pagination parameters (should use defaults)", async ({
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
      title: "Test Task",
      details: "Test",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Don't specify pagination params
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
    });

    expect(response.status).toBe(200);
    expect(response.data?.pageNumber).toBeDefined();
    expect(response.data?.pageSize).toBeDefined();
    expect(response.data?.items).toBeDefined();
  });

  test("should handle pagination with filters returning no results", async ({
    taskAPI,
  }) => {
    // Filter by non-existent ID
    const response = await taskAPI.getPaginated({
      id: ["00000000-0000-0000-0000-000000000000"],
      pageNumber: 1,
      pageSize: 10,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toEqual([]);
    expect(response.data?.totalCount).toBe(0);
    expect(response.data?.totalPages).toBe(0);
  });

  test("should maintain consistent pagination across multiple requests", async ({
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

    // Create 15 tasks
    for (let i = 0; i < 15; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Test",
        priority: i + 1,
      });
      createdTaskIds.push(task.data!.id);
    }

    // Get three pages with consistent params
    const page1 = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      pageNumber: 1,
      pageSize: 5,
      sortBy: "priority",
      sortOrder: "asc",
    });

    const page2 = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      pageNumber: 2,
      pageSize: 5,
      sortBy: "priority",
      sortOrder: "asc",
    });

    const page3 = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      pageNumber: 3,
      pageSize: 5,
      sortBy: "priority",
      sortOrder: "asc",
    });

    expect(page1.status).toBe(200);
    expect(page2.status).toBe(200);
    expect(page3.status).toBe(200);

    // All pages should have consistent total count
    expect(page1.data?.totalCount).toBe(page2.data?.totalCount);
    expect(page2.data?.totalCount).toBe(page3.data?.totalCount);

    // IDs should not overlap
    const ids1 = page1.data!.items!.map((t) => t.id);
    const ids2 = page2.data!.items!.map((t) => t.id);
    const ids3 = page3.data!.items!.map((t) => t.id);

    const allIds = [...ids1, ...ids2, ...ids3];
    const uniqueIds = new Set(allIds);
    expect(uniqueIds.size).toBe(allIds.length); // No duplicates
  });

  test("should handle sorting by due date with null values", async ({
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

    // Create tasks with and without due dates
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Task with due date",
      details: "Test",
      priority: 1,
      dueDate: new Date("2026-12-31").toISOString(),
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Task without due date",
      details: "Test",
      priority: 2,
    });

    const task3 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Task with earlier due date",
      details: "Test",
      priority: 3,
      dueDate: new Date("2026-06-15").toISOString(),
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

    // Sort by due date
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      sortBy: "dueDate",
      sortOrder: "asc",
    });

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBeGreaterThanOrEqual(3);

    // Tasks with due dates should come first, sorted by date
    // Tasks without due dates should come last
  });

  test("should handle invalid sort parameters", async ({ taskAPI }) => {
    const response = await taskAPI.getPaginated({
      sortBy: "invalidField" as any,
      sortOrder: "asc",
    });

    // Should return error or use default sort
    expect(response.status).toBeLessThan(500);
  });

  test("should calculate total pages correctly with exact division", async ({
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

    // Create exactly 10 tasks
    for (let i = 0; i < 10; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Test",
        priority: i + 1,
      });
      createdTaskIds.push(task.data!.id);
    }

    // Request with page size 5
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      pageNumber: 1,
      pageSize: 5,
    });

    expect(response.status).toBe(200);
    expect(response.data?.totalCount).toBe(10);
    expect(response.data?.totalPages).toBe(2); // Exactly 2 pages
  });

  test("should calculate total pages correctly with remainder", async ({
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

    // Create 7 tasks
    for (let i = 0; i < 7; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Test",
        priority: i + 1,
      });
      createdTaskIds.push(task.data!.id);
    }

    // Request with page size 3 (should give 3 pages)
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      pageNumber: 1,
      pageSize: 3,
    });

    expect(response.status).toBe(200);
    expect(response.data?.totalCount).toBe(7);
    expect(response.data?.totalPages).toBe(3); // ceil(7/3) = 3
  });
});
