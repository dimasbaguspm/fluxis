import { test, expect } from "../../../fixtures";

/**
 * Case: Project Logs
 * Tests project activity log retrieval and filtering
 */
test.describe("Case: Project Logs", () => {
  const createdProjectIds: string[] = [];

  test.afterEach(async ({ projectAPI }) => {
    for (const id of createdProjectIds) {
      await projectAPI.remove(id).catch(() => {});
    }
    createdProjectIds.length = 0;
  });

  test("should retrieve logs for a project", async ({ projectAPI }) => {
    const created = await projectAPI.create({
      name: "Project with Logs",
      description: "Testing log retrieval",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    const logsResponse = await projectAPI.getLogs(created.data!.id);

    expect(logsResponse.status).toBe(200);
    expect(logsResponse.data).toBeDefined();
    expect(logsResponse.data?.items).toBeDefined();
    expect(Array.isArray(logsResponse.data?.items)).toBe(true);
    expect(logsResponse.data?.pageNumber).toBeDefined();
    expect(logsResponse.data?.pageSize).toBeDefined();
    expect(logsResponse.data?.totalCount).toBeDefined();
    expect(logsResponse.data?.totalPages).toBeDefined();
  });

  test("should paginate project logs", async ({ projectAPI }) => {
    const created = await projectAPI.create({
      name: "Project for Log Pagination",
      description: "Testing pagination",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    // Get first page
    const page1 = await projectAPI.getLogs(created.data!.id, {
      pageNumber: 1,
      pageSize: 10,
    });

    expect(page1.status).toBe(200);
    expect(page1.data?.pageNumber).toBe(1);
    expect(page1.data?.pageSize).toBe(10);
  });

  test("should handle logs query for non-existent project", async ({
    projectAPI,
  }) => {
    const fakeId = "00000000-0000-0000-0000-000000000000";
    const response = await projectAPI.getLogs(fakeId);

    // API may return 200 with empty results or 404
    expect([200, 404]).toContain(response.status);
    if (response.status === 200) {
      expect(response.data?.items).toBeDefined();
    }
  });

  test("should filter logs by query parameter", async ({ projectAPI }) => {
    const created = await projectAPI.create({
      name: "Project for Log Search",
      description: "Testing search",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    const response = await projectAPI.getLogs(created.data!.id, {
      query: "search term",
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toBeDefined();
  });

  test("should handle empty log list for new project", async ({
    projectAPI,
  }) => {
    const created = await projectAPI.create({
      name: "New Project No Logs",
      description: "Brand new",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    const logsResponse = await projectAPI.getLogs(created.data!.id);

    expect(logsResponse.status).toBe(200);
    // New project might have creation log or be empty
    expect(logsResponse.data?.items).toBeDefined();
    expect(Array.isArray(logsResponse.data?.items)).toBe(true);
  });
});
