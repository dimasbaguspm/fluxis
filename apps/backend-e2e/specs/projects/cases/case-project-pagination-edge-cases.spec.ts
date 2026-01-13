import { test, expect } from "../../../fixtures";

/**
 * Project Pagination Edge Cases
 * Tests edge cases and boundary conditions for project pagination
 */
test.describe("Project Pagination Edge Cases", () => {
  const createdProjectIds: string[] = [];

  test.afterEach(async ({ projectAPI }) => {
    for (const id of createdProjectIds) {
      await projectAPI.remove(id).catch(() => {});
    }
    createdProjectIds.length = 0;
  });

  test("should handle page beyond total pages", async ({ projectAPI }) => {
    // Create a few projects
    for (let i = 0; i < 3; i++) {
      const project = await projectAPI.create({
        name: `Project ${i}`,
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);
    }

    // Request page 100 (way beyond available data)
    const response = await projectAPI.getPaginated({
      pageNumber: 100,
      pageSize: 10,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toEqual([]);
    expect(response.data?.pageNumber).toBe(100);
    // Total count reflects actual database state (may be 0 or more)
    expect(response.data?.totalCount).toBeGreaterThanOrEqual(0);
  });

  test("should handle very large page size", async ({ projectAPI }) => {
    // Create a few projects
    for (let i = 0; i < 5; i++) {
      const project = await projectAPI.create({
        name: `Project ${i}`,
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);
    }

    // Request very large page size
    const response = await projectAPI.getPaginated({
      pageNumber: 1,
      pageSize: 10000,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBeGreaterThanOrEqual(5);
    expect(response.data?.pageSize).toBe(10000);
  });

  test("should handle page size of 1", async ({ projectAPI }) => {
    // Create projects
    for (let i = 0; i < 3; i++) {
      const project = await projectAPI.create({
        name: `Project ${i}`,
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);
    }

    // Request page size of 1
    const response = await projectAPI.getPaginated({
      pageNumber: 1,
      pageSize: 1,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBe(1);
    expect(response.data?.pageSize).toBe(1);
    expect(response.data?.totalPages).toBeGreaterThanOrEqual(3);
  });

  test("should handle zero page number (should default to 1)", async ({
    projectAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const response = await projectAPI.getPaginated({
      pageNumber: 0,
      pageSize: 10,
    });

    // Should default to page 1 or return error
    expect(response.status).toBeLessThan(500);
  });

  test("should handle negative page number", async ({ projectAPI }) => {
    const response = await projectAPI.getPaginated({
      pageNumber: -1,
      pageSize: 10,
    });

    // Should return error or default to page 1
    expect(response.status).toBeLessThan(500);
  });

  test("should handle zero page size", async ({ projectAPI }) => {
    const response = await projectAPI.getPaginated({
      pageNumber: 1,
      pageSize: 0,
    });

    // Should return error or use default page size
    expect(response.status).toBeLessThan(500);
  });

  test("should handle negative page size", async ({ projectAPI }) => {
    const response = await projectAPI.getPaginated({
      pageNumber: 1,
      pageSize: -10,
    });

    // Should return error or use default page size
    expect(response.status).toBeLessThan(500);
  });

  test("should handle missing pagination parameters (should use defaults)", async ({
    projectAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Don't specify pagination params
    const response = await projectAPI.getPaginated({});

    expect(response.status).toBe(200);
    expect(response.data?.pageNumber).toBeDefined();
    expect(response.data?.pageSize).toBeDefined();
    expect(response.data?.items).toBeDefined();
  });

  test("should calculate total pages correctly with exact division", async ({
    projectAPI,
  }) => {
    // Create exactly 10 projects
    for (let i = 0; i < 10; i++) {
      const project = await projectAPI.create({
        name: `Project ${i}`,
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);
    }

    // Request with page size 5 (should give exactly 2 pages)
    const response = await projectAPI.getPaginated({
      pageNumber: 1,
      pageSize: 5,
    });

    expect(response.status).toBe(200);
    expect(response.data?.totalCount).toBeGreaterThanOrEqual(10);
    const expectedPages = Math.ceil(response.data!.totalCount / 5);
    expect(response.data?.totalPages).toBe(expectedPages);
  });

  test("should calculate total pages correctly with remainder", async ({
    projectAPI,
  }) => {
    // Create 7 projects
    for (let i = 0; i < 7; i++) {
      const project = await projectAPI.create({
        name: `Project ${i}`,
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);
    }

    // Request with page size 3 (should give 3 pages: 3, 3, 1)
    const response = await projectAPI.getPaginated({
      pageNumber: 1,
      pageSize: 3,
    });

    expect(response.status).toBe(200);
    expect(response.data?.totalCount).toBeGreaterThanOrEqual(7);
    const expectedPages = Math.ceil(response.data!.totalCount / 3);
    expect(response.data?.totalPages).toBe(expectedPages);
  });

  test("should handle pagination with filters returning no results", async ({
    projectAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Filter by non-existent ID
    const response = await projectAPI.getPaginated({
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
  }) => {
    // Create 15 projects
    for (let i = 0; i < 15; i++) {
      const project = await projectAPI.create({
        name: `Project ${i}`,
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);
    }

    // Get first page
    const page1 = await projectAPI.getPaginated({
      pageNumber: 1,
      pageSize: 5,
      sortBy: "createdAt",
      sortOrder: "asc",
    });

    // Get second page
    const page2 = await projectAPI.getPaginated({
      pageNumber: 2,
      pageSize: 5,
      sortBy: "createdAt",
      sortOrder: "asc",
    });

    // Get third page
    const page3 = await projectAPI.getPaginated({
      pageNumber: 3,
      pageSize: 5,
      sortBy: "createdAt",
      sortOrder: "asc",
    });

    expect(page1.status).toBe(200);
    expect(page2.status).toBe(200);
    expect(page3.status).toBe(200);

    // All pages should have consistent total count
    expect(page1.data?.totalCount).toBe(page2.data?.totalCount);
    expect(page2.data?.totalCount).toBe(page3.data?.totalCount);

    // IDs should not overlap
    const ids1 = page1.data!.items?.map((p) => p.id) || [];
    const ids2 = page2.data!.items?.map((p) => p.id) || [];
    const ids3 = page3.data!.items?.map((p) => p.id) || [];

    const allIds = [...ids1, ...ids2, ...ids3];
    const uniqueIds = new Set(allIds);
    expect(uniqueIds.size).toBe(allIds.length); // No duplicates
  });

  test("should handle empty database", async ({ projectAPI }) => {
    // Assuming test starts with empty state or we've cleaned up
    const response = await projectAPI.getPaginated({
      pageNumber: 1,
      pageSize: 10,
      // Filter to ensure we get no results
      id: ["00000000-0000-0000-0000-000000000000"],
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toEqual([]);
    expect(response.data?.totalCount).toBe(0);
    expect(response.data?.totalPages).toBe(0);
  });
});
