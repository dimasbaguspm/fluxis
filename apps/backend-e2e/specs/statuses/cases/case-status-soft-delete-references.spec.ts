import { test, expect } from "../../../fixtures";

/**
 * Status Soft-Delete References Test Cases
 * Tests status behavior when interacting with soft-deleted projects
 */
test.describe("Status Soft-Delete References", () => {
  const createdProjectIds: string[] = [];
  const createdStatusIds: string[] = [];

  test.afterEach(async ({ projectAPI, statusAPI }) => {
    for (const id of createdStatusIds) {
      await statusAPI.remove(id).catch(() => {});
    }
    for (const id of createdProjectIds) {
      await projectAPI.remove(id).catch(() => {});
    }
    createdStatusIds.length = 0;
    createdProjectIds.length = 0;
  });

  test("should return 404 when getting status from soft-deleted project", async ({
    projectAPI,
    statusAPI,
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

    // Soft delete the project
    await projectAPI.remove(project.data!.id);

    // Try to get status detail
    const response = await statusAPI.getById(status.id);

    // Should return 404 because parent project is soft-deleted
    expect(response.status).toBe(404);
  });

  test("should return empty list when querying statuses of soft-deleted project", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    expect(statuses.data!.length).toBeGreaterThan(0);

    // Soft delete the project
    await projectAPI.remove(project.data!.id);

    // Try to get statuses by project
    const response = await statusAPI.getByProject(project.data!.id);

    // Should return empty list or error
    expect(response.status).toBe(200);
    expect(response.data?.length).toBe(0);
  });

  test("should prevent updating status from soft-deleted project", async ({
    projectAPI,
    statusAPI,
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

    // Soft delete the project
    await projectAPI.remove(project.data!.id);

    // Try to update status
    const response = await statusAPI.update(status.id, {
      name: "Updated Name",
    });

    // Should return 404 because parent project is soft-deleted
    expect(response.status).toBe(404);
  });

  test("should prevent deleting status from soft-deleted project", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const status = statuses.data![1]; // Don't use default
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Soft delete the project
    await projectAPI.remove(project.data!.id);

    // Try to delete status
    const response = await statusAPI.remove(status.id);

    // Should return 404 because parent project is soft-deleted
    expect(response.status).toBe(404);
  });

  test("should prevent creating status for soft-deleted project", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Soft delete the project
    await projectAPI.remove(project.data!.id);

    // Try to create status
    const response = await statusAPI.create({
      projectId: project.data!.id,
      name: "New Status",
    });

    // Should fail because parent project is soft-deleted
    expect(response.status).toBeGreaterThanOrEqual(400);
  });

  test("should prevent reordering statuses for soft-deleted project", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "For soft delete tests",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    const statusIds = statuses.data!.map((s) => s.id);

    // Soft delete the project
    await projectAPI.remove(project.data!.id);

    // Try to reorder statuses
    const response = await statusAPI.reorder({
      projectId: project.data!.id,
      ids: statusIds.reverse(),
    });

    // Should fail because parent project is soft-deleted
    expect(response.status).toBeGreaterThanOrEqual(400);
  });

  test("should maintain status independence across projects with soft deletes", async ({
    projectAPI,
    statusAPI,
  }) => {
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
    const statuses2 = await statusAPI.getByProject(project2.data!.id);

    createdStatusIds.push(...statuses1.data!.map((s) => s.id));
    createdStatusIds.push(...statuses2.data!.map((s) => s.id));

    const status1 = statuses1.data![0];
    const status2 = statuses2.data![0];

    // Soft delete project1
    await projectAPI.remove(project1.data!.id);

    // Status from project1 should not be accessible
    const response1 = await statusAPI.getById(status1.id);
    expect(response1.status).toBe(404);

    // Status from project2 should still be accessible
    const response2 = await statusAPI.getById(status2.id);
    expect(response2.status).toBe(200);
    expect(response2.data?.id).toBe(status2.id);
  });
});
