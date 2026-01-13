import { test, expect } from "../../fixtures";

/**
 * Project API - Common CRUD Tests
 * Tests standard operations: Create, Read, Update, Delete, List
 */
test.describe("Project API", () => {
  // Store created project IDs for cleanup
  const createdProjectIds: string[] = [];

  // Cleanup after each test
  test.afterEach(async ({ projectAPI }) => {
    for (const id of createdProjectIds) {
      await projectAPI.remove(id).catch(() => {
        // Ignore errors (project may already be deleted)
      });
    }
    createdProjectIds.length = 0;
  });

  test.describe("POST /projects", () => {
    test("should create project with all required fields", async ({
      projectAPI,
    }) => {
      const response = await projectAPI.create({
        name: "Test Project",
        description: "Test project description",
        status: "active",
      });

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.id).toBeDefined();
      expect(response.data?.name).toBe("Test Project");
      expect(response.data?.description).toBe("Test project description");
      expect(response.data?.status).toBe("active");
      expect(response.data?.createdAt).toBeDefined();
      expect(response.data?.updatedAt).toBeDefined();

      // Verify timestamps are valid dates
      expect(new Date(response.data!.createdAt).getTime()).toBeGreaterThan(0);
      expect(new Date(response.data!.updatedAt).getTime()).toBeGreaterThan(0);

      createdProjectIds.push(response.data!.id);
    });

    test("should create project with paused status", async ({ projectAPI }) => {
      const response = await projectAPI.create({
        name: "Paused Project",
        description: "This project is paused",
        status: "paused",
      });

      expect(response.status).toBe(200);
      expect(response.data?.status).toBe("paused");

      createdProjectIds.push(response.data!.id);
    });

    test("should create project with archived status", async ({
      projectAPI,
    }) => {
      const response = await projectAPI.create({
        name: "Archived Project",
        description: "This project is archived",
        status: "archived",
      });

      expect(response.status).toBe(200);
      expect(response.data?.status).toBe("archived");

      createdProjectIds.push(response.data!.id);
    });

    test("should fail to create project with missing name", async ({
      projectAPI,
    }) => {
      const response = await projectAPI.create({
        description: "Description without name",
        status: "active",
      } as any);

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.status).toBeLessThan(500);
      expect(response.error).toBeDefined();
    });

    test("should fail to create project with missing description", async ({
      projectAPI,
    }) => {
      const response = await projectAPI.create({
        name: "Project without description",
        status: "active",
      } as any);

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create project with missing status", async ({
      projectAPI,
    }) => {
      const response = await projectAPI.create({
        name: "Project without status",
        description: "Description",
      } as any);

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create project with invalid status", async ({
      projectAPI,
    }) => {
      const response = await projectAPI.create({
        name: "Project with invalid status",
        description: "Description",
        status: "invalid_status" as any,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create project with empty name", async ({
      projectAPI,
    }) => {
      const response = await projectAPI.create({
        name: "",
        description: "Description",
        status: "active",
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create project with empty description", async ({
      projectAPI,
    }) => {
      const response = await projectAPI.create({
        name: "Project Name",
        description: "",
        status: "active",
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("GET /projects/{projectId}", () => {
    test("should get project by ID", async ({ projectAPI }) => {
      // Create a project first
      const created = await projectAPI.create({
        name: "Get Project Test",
        description: "Test getting project",
        status: "active",
      });
      createdProjectIds.push(created.data!.id);

      // Get the project
      const response = await projectAPI.getById(created.data!.id);

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.id).toBe(created.data!.id);
      expect(response.data?.name).toBe("Get Project Test");
      expect(response.data?.description).toBe("Test getting project");
      expect(response.data?.status).toBe("active");
    });

    test("should return 404 for non-existent project", async ({
      projectAPI,
    }) => {
      const fakeId = "00000000-0000-0000-0000-000000000000";
      const response = await projectAPI.getById(fakeId);

      expect(response.status).toBe(404);
      expect(response.error).toBeDefined();
    });

    test("should return 400 for invalid UUID format", async ({
      projectAPI,
    }) => {
      const response = await projectAPI.getById("invalid-uuid");

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("GET /projects", () => {
    test("should get paginated list of projects", async ({ projectAPI }) => {
      // Create multiple projects
      for (let i = 1; i <= 3; i++) {
        const created = await projectAPI.create({
          name: `Paginated Project ${i}`,
          description: `Description ${i}`,
          status: "active",
        });
        createdProjectIds.push(created.data!.id);
      }

      // Get paginated list
      const response = await projectAPI.getPaginated({
        pageNumber: 1,
        pageSize: 10,
      });

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.items).toBeDefined();
      expect(Array.isArray(response.data?.items)).toBe(true);
      expect(response.data?.pageNumber).toBe(1);
      expect(response.data?.pageSize).toBe(10);
      expect(response.data?.totalCount).toBeGreaterThanOrEqual(3);
      expect(response.data?.totalPages).toBeGreaterThanOrEqual(1);
    });

    test("should filter projects by status", async ({ projectAPI }) => {
      // Create projects with different statuses
      const active = await projectAPI.create({
        name: "Active Project",
        description: "Active",
        status: "active",
      });
      const paused = await projectAPI.create({
        name: "Paused Project",
        description: "Paused",
        status: "paused",
      });
      createdProjectIds.push(active.data!.id, paused.data!.id);

      // Filter by active status
      const response = await projectAPI.getPaginated({
        status: ["active"],
      });

      expect(response.status).toBe(200);
      expect(response.data?.items).toBeDefined();

      // Verify all returned projects are active
      const allActive = response.data!.items!.every(
        (p) => p.status === "active"
      );
      expect(allActive).toBe(true);
    });

    test("should filter projects by multiple statuses", async ({
      projectAPI,
    }) => {
      // Create projects with different statuses
      const active = await projectAPI.create({
        name: "Active Project",
        description: "Active",
        status: "active",
      });
      const paused = await projectAPI.create({
        name: "Paused Project",
        description: "Paused",
        status: "paused",
      });
      const archived = await projectAPI.create({
        name: "Archived Project",
        description: "Archived",
        status: "archived",
      });
      createdProjectIds.push(
        active.data!.id,
        paused.data!.id,
        archived.data!.id
      );

      // Filter by active and paused statuses
      const response = await projectAPI.getPaginated({
        status: ["active", "paused"],
      });

      expect(response.status).toBe(200);
      expect(response.data?.items).toBeDefined();

      // Verify all returned projects are either active or paused
      const validStatuses = response.data!.items!.every(
        (p) => p.status === "active" || p.status === "paused"
      );
      expect(validStatuses).toBe(true);
    });

    test("should search projects by query", async ({ projectAPI }) => {
      const uniqueName = `Searchable Project ${Date.now()}`;
      const created = await projectAPI.create({
        name: uniqueName,
        description: "Find me by query",
        status: "active",
      });
      createdProjectIds.push(created.data!.id);

      // Search by name
      const response = await projectAPI.getPaginated({
        query: uniqueName,
      });

      expect(response.status).toBe(200);
      expect(response.data?.items).toBeDefined();
      expect(response.data!.items!.length).toBeGreaterThanOrEqual(1);

      // Verify our project is in the results
      const found = response.data!.items!.find(
        (p) => p.id === created.data!.id
      );
      expect(found).toBeDefined();
    });

    test("should sort projects by createdAt desc", async ({ projectAPI }) => {
      const response = await projectAPI.getPaginated({
        sortBy: "createdAt",
        sortOrder: "desc",
        pageSize: 5,
      });

      expect(response.status).toBe(200);
      expect(response.data?.items).toBeDefined();

      // Verify descending order
      if (response.data!.items!.length > 1) {
        for (let i = 0; i < response.data!.items!.length - 1; i++) {
          const current = new Date(response.data!.items![i].createdAt);
          const next = new Date(response.data!.items![i + 1].createdAt);
          expect(current.getTime()).toBeGreaterThanOrEqual(next.getTime());
        }
      }
    });

    test("should sort projects by createdAt asc", async ({ projectAPI }) => {
      const response = await projectAPI.getPaginated({
        sortBy: "createdAt",
        sortOrder: "asc",
        pageSize: 5,
      });

      expect(response.status).toBe(200);
      expect(response.data?.items).toBeDefined();

      // Verify ascending order
      if (response.data!.items!.length > 1) {
        for (let i = 0; i < response.data!.items!.length - 1; i++) {
          const current = new Date(response.data!.items![i].createdAt);
          const next = new Date(response.data!.items![i + 1].createdAt);
          expect(current.getTime()).toBeLessThanOrEqual(next.getTime());
        }
      }
    });

    test("should handle pagination correctly", async ({ projectAPI }) => {
      // Get first page
      const page1 = await projectAPI.getPaginated({
        pageNumber: 1,
        pageSize: 5,
      });

      expect(page1.status).toBe(200);
      expect(page1.data?.pageNumber).toBe(1);
      expect(page1.data?.pageSize).toBe(5);

      // If there are more pages, get the second page
      if (page1.data!.totalPages! > 1) {
        const page2 = await projectAPI.getPaginated({
          pageNumber: 2,
          pageSize: 5,
        });

        expect(page2.status).toBe(200);
        expect(page2.data?.pageNumber).toBe(2);

        // Verify no overlap between pages
        const page1Ids = page1.data!.items!.map((p) => p.id);
        const page2Ids = page2.data!.items!.map((p) => p.id);
        const overlap = page1Ids.filter((id) => page2Ids.includes(id));
        expect(overlap).toHaveLength(0);
      }
    });
  });

  test.describe("PATCH /projects/{projectId}", () => {
    test("should update project name", async ({ projectAPI }) => {
      const created = await projectAPI.create({
        name: "Original Name",
        description: "Description",
        status: "active",
      });
      createdProjectIds.push(created.data!.id);

      const response = await projectAPI.update(created.data!.id, {
        name: "Updated Name",
        description: "Description",
        status: "active",
      });

      expect(response.status).toBe(200);
      expect(response.data?.name).toBe("Updated Name");
      expect(response.data?.description).toBe("Description");
      expect(response.data?.status).toBe("active");
      expect(response.data?.updatedAt).not.toBe(created.data!.updatedAt);
    });

    test("should update project description", async ({ projectAPI }) => {
      const created = await projectAPI.create({
        name: "Project Name",
        description: "Original Description",
        status: "active",
      });
      createdProjectIds.push(created.data!.id);

      const response = await projectAPI.update(created.data!.id, {
        name: "Project Name",
        description: "Updated Description",
        status: "active",
      });

      expect(response.status).toBe(200);
      expect(response.data?.description).toBe("Updated Description");
      expect(response.data?.name).toBe("Project Name");
    });

    test("should update project status", async ({ projectAPI }) => {
      const created = await projectAPI.create({
        name: "Project Name",
        description: "Description",
        status: "active",
      });
      createdProjectIds.push(created.data!.id);

      const response = await projectAPI.update(created.data!.id, {
        status: "paused",
      });

      expect(response.status).toBe(200);
      expect(response.data?.status).toBe("paused");
    });

    test("should update multiple fields at once", async ({ projectAPI }) => {
      const created = await projectAPI.create({
        name: "Original Name",
        description: "Original Description",
        status: "active",
      });
      createdProjectIds.push(created.data!.id);

      const response = await projectAPI.update(created.data!.id, {
        name: "New Name",
        description: "New Description",
        status: "archived",
      });

      expect(response.status).toBe(200);
      expect(response.data?.name).toBe("New Name");
      expect(response.data?.description).toBe("New Description");
      expect(response.data?.status).toBe("archived");
    });

    test("should return 400 when updating with invalid data", async ({
      projectAPI,
    }) => {
      const created = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(created.data!.id);

      const response = await projectAPI.update(created.data!.id, {
        name: "New Name",
        // Missing required fields
      });

      expect(response.status).toBe(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to update with empty name", async ({ projectAPI }) => {
      const created = await projectAPI.create({
        name: "Original Name",
        description: "Description",
        status: "active",
      });
      createdProjectIds.push(created.data!.id);

      const response = await projectAPI.update(created.data!.id, {
        name: "",
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to update with invalid status", async ({
      projectAPI,
    }) => {
      const created = await projectAPI.create({
        name: "Project Name",
        description: "Description",
        status: "active",
      });
      createdProjectIds.push(created.data!.id);

      const response = await projectAPI.update(created.data!.id, {
        status: "invalid_status" as any,
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("DELETE /projects/{projectId}", () => {
    test("should delete project", async ({ projectAPI }) => {
      const created = await projectAPI.create({
        name: "Project to Delete",
        description: "This will be deleted",
        status: "active",
      });

      const deleteResponse = await projectAPI.remove(created.data!.id);

      expect(deleteResponse.status).toBe(204);

      // Verify project is deleted
      const getResponse = await projectAPI.getById(created.data!.id);
      expect(getResponse.status).toBe(404);
    });

    test("should return 404 when deleting non-existent project", async ({
      projectAPI,
    }) => {
      const fakeId = "00000000-0000-0000-0000-000000000000";
      const response = await projectAPI.remove(fakeId);

      expect(response.status).toBe(404);
    });

    test("should be idempotent - deleting twice should not error", async ({
      projectAPI,
    }) => {
      const created = await projectAPI.create({
        name: "Project to Delete Twice",
        description: "Testing idempotency",
        status: "active",
      });

      // First delete
      const firstDelete = await projectAPI.remove(created.data!.id);
      expect(firstDelete.status).toBe(204);

      // Second delete (should return 404 but not crash)
      const secondDelete = await projectAPI.remove(created.data!.id);
      expect(secondDelete.status).toBe(404);
    });
  });
});
