import { test, expect } from "../../../fixtures";

/**
 * Case: Project Bulk Operations
 * Tests scenarios involving multiple projects and bulk filtering
 */
test.describe("Case: Project Bulk Operations", () => {
  const createdProjectIds: string[] = [];

  test.afterEach(async ({ projectAPI }) => {
    for (const id of createdProjectIds) {
      await projectAPI.remove(id).catch(() => {});
    }
    createdProjectIds.length = 0;
  });

  test("should filter projects by multiple IDs", async ({ projectAPI }) => {
    // Create multiple projects
    const project1 = await projectAPI.create({
      name: "Bulk Project 1",
      description: "First",
      status: "active",
    });
    const project2 = await projectAPI.create({
      name: "Bulk Project 2",
      description: "Second",
      status: "active",
    });
    const project3 = await projectAPI.create({
      name: "Bulk Project 3",
      description: "Third",
      status: "active",
    });

    createdProjectIds.push(
      project1.data!.id,
      project2.data!.id,
      project3.data!.id
    );

    // Filter by specific IDs
    const response = await projectAPI.getPaginated({
      id: [project1.data!.id, project2.data!.id],
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toBeDefined();

    // Should only return the 2 requested projects
    const returnedIds = response.data!.items!.map((p) => p.id);
    expect(returnedIds).toContain(project1.data!.id);
    expect(returnedIds).toContain(project2.data!.id);
  });

  test("should handle large number of projects", async ({ projectAPI }) => {
    // Create 15 projects
    for (let i = 1; i <= 15; i++) {
      const created = await projectAPI.create({
        name: `Bulk Test Project ${i}`,
        description: `Description ${i}`,
        status: i % 3 === 0 ? "archived" : i % 2 === 0 ? "paused" : "active",
      });
      createdProjectIds.push(created.data!.id);
    }

    // Get first page
    const page1 = await projectAPI.getPaginated({
      pageNumber: 1,
      pageSize: 10,
    });

    expect(page1.status).toBe(200);
    expect(page1.data?.items).toBeDefined();
    expect(page1.data?.totalCount).toBeGreaterThanOrEqual(15);

    // Get second page
    const page2 = await projectAPI.getPaginated({
      pageNumber: 2,
      pageSize: 10,
    });

    expect(page2.status).toBe(200);
    expect(page2.data?.items).toBeDefined();
  });

  test("should combine multiple filters", async ({ projectAPI }) => {
    const timestamp = Date.now();

    // Create projects with different combinations
    const activeSearch = await projectAPI.create({
      name: `Searchable Active ${timestamp}`,
      description: "Active and searchable",
      status: "active",
    });
    const pausedSearch = await projectAPI.create({
      name: `Searchable Paused ${timestamp}`,
      description: "Paused and searchable",
      status: "paused",
    });
    const activeNonSearch = await projectAPI.create({
      name: `NonSearch Active ${timestamp}`,
      description: "Active but not searchable",
      status: "active",
    });

    createdProjectIds.push(
      activeSearch.data!.id,
      pausedSearch.data!.id,
      activeNonSearch.data!.id
    );

    // Search for "Searchable" with status filter
    const response = await projectAPI.getPaginated({
      query: "Searchable",
      status: ["active"],
    });

    expect(response.status).toBe(200);

    // Should find the active searchable project
    const found = response.data!.items!.find(
      (p) => p.id === activeSearch.data!.id
    );
    expect(found).toBeDefined();

    // Should not find paused searchable project (status filter)
    const pausedFound = response.data!.items!.find(
      (p) => p.id === pausedSearch.data!.id
    );
    expect(pausedFound).toBeUndefined();
  });

  test("should sort and paginate correctly with filters", async ({
    projectAPI,
  }) => {
    // Create multiple projects
    for (let i = 1; i <= 5; i++) {
      const created = await projectAPI.create({
        name: `Sorted Project ${i}`,
        description: `Description ${i}`,
        status: "active",
      });
      createdProjectIds.push(created.data!.id);
    }

    // Get sorted list
    const response = await projectAPI.getPaginated({
      status: ["active"],
      sortBy: "createdAt",
      sortOrder: "desc",
      pageSize: 3,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toBeDefined();

    // Verify sorting
    if (response.data!.items!.length > 1) {
      for (let i = 0; i < response.data!.items!.length - 1; i++) {
        const current = new Date(response.data!.items![i].createdAt);
        const next = new Date(response.data!.items![i + 1].createdAt);
        expect(current.getTime()).toBeGreaterThanOrEqual(next.getTime());
      }
    }
  });

  test("should handle empty results with filters", async ({ projectAPI }) => {
    const response = await projectAPI.getPaginated({
      query: `nonexistent-project-${Date.now()}-${Math.random()}`,
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toBeDefined();
    expect(response.data?.totalCount).toBe(0);
    expect(response.data?.items).toHaveLength(0);
  });
});
