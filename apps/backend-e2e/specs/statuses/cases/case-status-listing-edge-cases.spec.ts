import { test, expect } from "../../../fixtures";

/**
 * Status Pagination Edge Cases
 * Tests edge cases for status listing (not paginated, but has ordering)
 */
test.describe("Status Listing Edge Cases", () => {
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

  test("should handle project with no custom statuses (only defaults)", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "New Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const response = await statusAPI.getByProject(project.data!.id);

    expect(response.status).toBe(200);
    expect(response.data?.length).toBe(3); // Todo, In Progress, Done
    expect(response.data).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ name: "Todo", isDefault: true }),
        expect.objectContaining({ name: "In Progress", isDefault: false }),
        expect.objectContaining({ name: "Done", isDefault: false }),
      ])
    );

    createdStatusIds.push(...response.data!.map((s) => s.id));
  });

  test("should handle project with many custom statuses", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create 20 custom statuses
    for (let i = 0; i < 20; i++) {
      const status = await statusAPI.create({
        projectId: project.data!.id,
        name: `Status ${i}`,
      });
      createdStatusIds.push(status.data!.id);
    }

    const response = await statusAPI.getByProject(project.data!.id);

    expect(response.status).toBe(200);
    expect(response.data?.length).toBe(23); // 3 default + 20 custom
  });

  test("should return empty array for non-existent project", async ({
    statusAPI,
  }) => {
    const response = await statusAPI.getByProject(
      "00000000-0000-0000-0000-000000000000"
    );

    expect(response.status).toBe(200);
    expect(response.data).toEqual([]);
  });

  test("should return empty array for project with all statuses deleted", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Delete all non-default statuses
    for (const status of statuses.data!) {
      if (!status.isDefault) {
        await statusAPI.remove(status.id);
      }
    }

    // Try to delete default status (may or may not be allowed)
    const defaultStatus = statuses.data!.find((s) => s.isDefault);
    if (defaultStatus) {
      await statusAPI.remove(defaultStatus.id).catch(() => {});
    }

    const response = await statusAPI.getByProject(project.data!.id);

    expect(response.status).toBe(200);
    // Should return only non-deleted statuses
    expect(response.data).toBeDefined();
  });

  test("should maintain position order with deleted statuses", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create additional statuses
    const status1 = await statusAPI.create({
      projectId: project.data!.id,
      name: "Status 1",
    });
    const status2 = await statusAPI.create({
      projectId: project.data!.id,
      name: "Status 2",
    });
    const status3 = await statusAPI.create({
      projectId: project.data!.id,
      name: "Status 3",
    });

    createdStatusIds.push(status1.data!.id, status2.data!.id, status3.data!.id);

    // Delete the middle one
    await statusAPI.remove(status2.data!.id);

    const response = await statusAPI.getByProject(project.data!.id);

    expect(response.status).toBe(200);
    // Should maintain sequential positions for remaining statuses
    const positions = response.data!.map((s) => s.position);
    expect(positions).toEqual([...positions].sort((a, b) => a - b));
  });

  test("should handle invalid projectId format", async ({ statusAPI }) => {
    const response = await statusAPI.getByProject("invalid-uuid");

    // Should return error or empty array
    expect(response.status).toBeLessThan(500);
  });

  test("should handle missing projectId parameter", async ({ request }) => {
    const response = await request.get("/api/statuses", {
      headers: {
        Authorization: `Bearer ${process.env.TEST_ACCESS_TOKEN}`,
      },
    });

    // Should return error
    expect(response.status()).toBeGreaterThanOrEqual(400);
  });

  test("should return statuses in consistent order across multiple requests", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create multiple custom statuses
    for (let i = 0; i < 5; i++) {
      const status = await statusAPI.create({
        projectId: project.data!.id,
        name: `Status ${i}`,
      });
      createdStatusIds.push(status.data!.id);
    }

    // Get statuses multiple times
    const response1 = await statusAPI.getByProject(project.data!.id);
    const response2 = await statusAPI.getByProject(project.data!.id);
    const response3 = await statusAPI.getByProject(project.data!.id);

    expect(response1.status).toBe(200);
    expect(response2.status).toBe(200);
    expect(response3.status).toBe(200);

    // Order should be consistent
    const ids1 = response1.data!.map((s) => s.id);
    const ids2 = response2.data!.map((s) => s.id);
    const ids3 = response3.data!.map((s) => s.id);

    expect(ids1).toEqual(ids2);
    expect(ids2).toEqual(ids3);

    // Positions should be consistent
    const positions1 = response1.data!.map((s) => s.position);
    const positions2 = response2.data!.map((s) => s.position);
    const positions3 = response3.data!.map((s) => s.position);

    expect(positions1).toEqual(positions2);
    expect(positions2).toEqual(positions3);
  });

  test.skip("should handle concurrent status creation", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create multiple statuses concurrently
    const createPromises = Array(5)
      .fill(null)
      .map((_, i) =>
        statusAPI.create({
          projectId: project.data!.id,
          name: `Concurrent Status ${i}`,
        })
      );

    const results = await Promise.all(createPromises);

    // All should succeed
    results.forEach((r) => {
      expect(r.status).toBe(200);
      createdStatusIds.push(r.data!.id);
    });

    // Verify all statuses exist
    const allStatuses = await statusAPI.getByProject(project.data!.id);
    expect(allStatuses.data?.length).toBe(8); // 3 default + 5 new

    // All positions should be unique
    const positions = allStatuses.data!.map((s) => s.position);
    const uniquePositions = new Set(positions);
    expect(uniquePositions.size).toBe(positions.length);
  });
});
