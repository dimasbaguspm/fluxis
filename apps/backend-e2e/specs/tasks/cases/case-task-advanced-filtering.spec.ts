import { test, expect } from "../../../fixtures";

/**
 * Task Advanced Filtering Test Cases
 * Tests complex filtering combinations for tasks
 */
test.describe("Task Advanced Filtering", () => {
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

  test("should filter by multiple task IDs", async ({
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

    // Create multiple tasks
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Task 1",
      details: "Test",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Task 2",
      details: "Test",
      priority: 2,
    });

    const task3 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Task 3",
      details: "Test",
      priority: 3,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

    // Filter by specific IDs
    const response = await taskAPI.getPaginated({
      id: [task1.data!.id, task3.data!.id],
    });

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBe(2);

    const ids = response.data!.items!.map((t) => t.id);
    expect(ids).toContain(task1.data!.id);
    expect(ids).toContain(task3.data!.id);
    expect(ids).not.toContain(task2.data!.id);
  });

  test("should filter by multiple project IDs", async ({
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
    const statuses2 = await statusAPI.getByProject(project2.data!.id);

    createdStatusIds.push(...statuses1.data!.map((s) => s.id));
    createdStatusIds.push(...statuses2.data!.map((s) => s.id));

    // Create tasks in both projects
    const task1 = await taskAPI.create({
      projectId: project1.data!.id,
      statusId: statuses1.data![0].id,
      title: "Task in Project 1",
      details: "Test",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project2.data!.id,
      statusId: statuses2.data![0].id,
      title: "Task in Project 2",
      details: "Test",
      priority: 1,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id);

    // Filter by both project IDs
    const response = await taskAPI.getPaginated({
      projectId: [project1.data!.id, project2.data!.id],
    });

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBeGreaterThanOrEqual(2);

    const ids = response.data!.items!.map((t) => t.id);
    expect(ids).toContain(task1.data!.id);
    expect(ids).toContain(task2.data!.id);
  });

  test("should filter by multiple status IDs", async ({
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
    const status1 = statuses.data![0];
    const status2 = statuses.data![1];
    const status3 = statuses.data![2];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks in different statuses
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Task 1",
      details: "Test",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status2.id,
      title: "Task 2",
      details: "Test",
      priority: 1,
    });

    const task3 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status3.id,
      title: "Task 3",
      details: "Test",
      priority: 1,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

    // Filter by two specific statuses
    const response = await taskAPI.getPaginated({
      statusId: [status1.id, status3.id],
    });

    expect(response.status).toBe(200);

    const ids = response.data!.items!.map((t) => t.id);
    expect(ids).toContain(task1.data!.id);
    expect(ids).toContain(task3.data!.id);
    expect(ids).not.toContain(task2.data!.id);
  });

  test("should combine projectId and statusId filters", async ({
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
    const statuses2 = await statusAPI.getByProject(project2.data!.id);

    createdStatusIds.push(...statuses1.data!.map((s) => s.id));
    createdStatusIds.push(...statuses2.data!.map((s) => s.id));

    // Create tasks
    const task1 = await taskAPI.create({
      projectId: project1.data!.id,
      statusId: statuses1.data![0].id,
      title: "P1 Task 1",
      details: "Test",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project1.data!.id,
      statusId: statuses1.data![1].id,
      title: "P1 Task 2",
      details: "Test",
      priority: 1,
    });

    const task3 = await taskAPI.create({
      projectId: project2.data!.id,
      statusId: statuses2.data![0].id,
      title: "P2 Task 1",
      details: "Test",
      priority: 1,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

    // Filter by project1 AND specific status
    const response = await taskAPI.getPaginated({
      projectId: [project1.data!.id],
      statusId: [statuses1.data![0].id],
    });

    expect(response.status).toBe(200);

    const ids = response.data!.items!.map((t) => t.id);
    expect(ids).toContain(task1.data!.id);
    expect(ids).not.toContain(task2.data!.id); // Different status
    expect(ids).not.toContain(task3.data!.id); // Different project
  });

  test("should combine filters with search query", async ({
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
    const status1 = statuses.data![0];
    const status2 = statuses.data![1];
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create tasks with different titles
    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Implement login feature",
      details: "Backend work",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status1.id,
      title: "Design login page",
      details: "Frontend design",
      priority: 2,
    });

    const task3 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status2.id,
      title: "Implement dashboard",
      details: "Backend work",
      priority: 1,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

    // Search for "login" in specific status
    const response = await taskAPI.getPaginated({
      statusId: [status1.id],
      query: "login",
    });

    expect(response.status).toBe(200);

    const ids = response.data!.items!.map((t) => t.id);
    expect(ids).toContain(task1.data!.id); // Has "login" in status1
    expect(ids).toContain(task2.data!.id); // Has "login" in status1
    expect(ids).not.toContain(task3.data!.id); // Different status
  });

  test("should search in both title and details", async ({
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

    const task1 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Backend API development",
      details: "Node.js service",
      priority: 1,
    });

    const task2 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Frontend application",
      details: "Backend integration required",
      priority: 2,
    });

    const task3 = await taskAPI.create({
      projectId: project.data!.id,
      statusId: status.id,
      title: "Mobile app",
      details: "iOS development",
      priority: 3,
    });

    createdTaskIds.push(task1.data!.id, task2.data!.id, task3.data!.id);

    // Search for "Backend" (appears in title of task1, details of task2)
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      query: "Backend",
    });

    expect(response.status).toBe(200);

    const ids = response.data!.items!.map((t) => t.id);
    expect(ids).toContain(task1.data!.id); // Title matches
    expect(ids).toContain(task2.data!.id); // Details match
    expect(ids).not.toContain(task3.data!.id); // No match
  });

  test("should handle case-insensitive search", async ({
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
      title: "React Development",
      details: "Frontend work",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    // Search with different cases
    const lower = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      query: "react",
    });
    const upper = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      query: "REACT",
    });
    const mixed = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      query: "ReAcT",
    });

    // All should find the task
    expect(lower.data!.items!.some((t) => t.id === task.data!.id)).toBe(true);
    expect(upper.data!.items!.some((t) => t.id === task.data!.id)).toBe(true);
    expect(mixed.data!.items!.some((t) => t.id === task.data!.id)).toBe(true);
  });

  test("should combine all filters with pagination and sorting", async ({
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

    // Create multiple tasks
    for (let i = 0; i < 10; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i} Feature`,
        details: "Implementation work",
        priority: i + 1,
      });
      createdTaskIds.push(task.data!.id);
    }

    // Complex query: filter + search + sort + paginate
    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      statusId: [status.id],
      query: "Feature",
      sortBy: "priority",
      sortOrder: "desc",
      pageNumber: 1,
      pageSize: 5,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBe(5);

    // Should be sorted by priority descending
    const priorities = response.data!.items!.map((t) => t.priority);
    for (let i = 1; i < priorities.length; i++) {
      expect(priorities[i]).toBeLessThanOrEqual(priorities[i - 1]);
    }
  });

  test("should handle unicode characters in search", async ({
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
      title: "実装タスク (Implementation Task)",
      details: "日本語の説明",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    const response = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      query: "実装",
    });

    expect(response.status).toBe(200);
    expect(response.data!.items!.some((t) => t.id === task.data!.id)).toBe(
      true
    );
  });

  test("should handle special characters in search", async ({
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
      title: "Task @2024 #important",
      details: "Cost: $1000",
      priority: 1,
    });
    createdTaskIds.push(task.data!.id);

    const response1 = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      query: "@2024",
    });

    const response2 = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      query: "#important",
    });

    expect(response1.status).toBe(200);
    expect(response2.status).toBe(200);
  });

  test("should handle filters returning no results", async ({ taskAPI }) => {
    const response = await taskAPI.getPaginated({
      id: ["00000000-0000-0000-0000-000000000000"],
      projectId: ["00000000-0000-0000-0000-000000000000"],
      statusId: ["00000000-0000-0000-0000-000000000000"],
      query: "nonexistent-xyz-abc",
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toEqual([]);
    expect(response.data?.totalCount).toBe(0);
  });

  test("should maintain filter consistency across multiple pages", async ({
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

    // Create 20 tasks
    for (let i = 0; i < 20; i++) {
      const task = await taskAPI.create({
        projectId: project.data!.id,
        statusId: status.id,
        title: `Task ${i}`,
        details: "Test",
        priority: i + 1,
      });
      createdTaskIds.push(task.data!.id);
    }

    // Get all pages with same filter
    const page1 = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      statusId: [status.id],
      pageNumber: 1,
      pageSize: 7,
    });

    const page2 = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      statusId: [status.id],
      pageNumber: 2,
      pageSize: 7,
    });

    const page3 = await taskAPI.getPaginated({
      projectId: [project.data!.id],
      statusId: [status.id],
      pageNumber: 3,
      pageSize: 7,
    });

    // All should return consistent total counts
    expect(page1.data?.totalCount).toBe(page2.data?.totalCount);
    expect(page2.data?.totalCount).toBe(page3.data?.totalCount);

    // No overlapping IDs
    const ids1 = page1.data!.items!.map((t) => t.id);
    const ids2 = page2.data!.items!.map((t) => t.id);
    const ids3 = page3.data!.items!.map((t) => t.id);

    const overlap12 = ids1.filter((id) => ids2.includes(id));
    const overlap23 = ids2.filter((id) => ids3.includes(id));
    const overlap13 = ids1.filter((id) => ids3.includes(id));

    expect(overlap12.length).toBe(0);
    expect(overlap23.length).toBe(0);
    expect(overlap13.length).toBe(0);
  });
});
