import { test, expect } from "../../fixtures";

/**
 * Status API - Common CRUD Tests
 * Tests standard operations: Create, Read, Update, Delete, List
 */
test.describe("Status API", () => {
  const createdProjectIds: string[] = [];
  const createdStatusIds: string[] = [];

  // Cleanup after each test
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

  test.describe("POST /statuses", () => {
    test("should create status with required fields", async ({
      projectAPI,
      statusAPI,
    }) => {
      // Create a project first
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For status tests",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const response = await statusAPI.create({
        projectId: project.data!.id,
        name: "To Do",
      });

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.id).toBeDefined();
      expect(response.data?.projectId).toBe(project.data!.id);
      expect(response.data?.name).toBe("To Do");
      expect(response.data?.slug).toBeDefined();
      expect(response.data?.position).toBeDefined();
      expect(response.data?.isDefault).toBeDefined();
      expect(response.data?.createdAt).toBeDefined();
      expect(response.data?.updatedAt).toBeDefined();

      createdStatusIds.push(response.data!.id);
    });

    test("should auto-generate slug from name", async ({
      projectAPI,
      statusAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For slug test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const response = await statusAPI.create({
        projectId: project.data!.id,
        name: "In Progress",
      });

      expect(response.status).toBe(200);
      expect(response.data?.slug).toBeDefined();
      expect(response.data?.slug).toMatch(/^[a-z0-9_]+$/);

      createdStatusIds.push(response.data!.id);
    });

    test("should set position automatically", async ({
      projectAPI,
      statusAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "For position test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      // Get auto-created statuses first
      const existingStatuses = await statusAPI.getByProject(project.data!.id);
      const maxPosition = Math.max(
        ...existingStatuses.data!.map((s) => s.position)
      );

      const status1 = await statusAPI.create({
        projectId: project.data!.id,
        name: "First Status",
      });
      const status2 = await statusAPI.create({
        projectId: project.data!.id,
        name: "Second Status",
      });

      expect(status1.data?.position).toBeDefined();
      expect(status2.data?.position).toBeDefined();
      expect(status1.data!.position).toBeGreaterThan(maxPosition);
      expect(status2.data!.position).toBeGreaterThan(status1.data!.position);

      createdStatusIds.push(status1.data!.id, status2.data!.id);
    });

    test("should fail to create status with missing projectId", async ({
      statusAPI,
    }) => {
      const response = await statusAPI.create({
        name: "Invalid Status",
      } as any);

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create status with missing name", async ({
      projectAPI,
      statusAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const response = await statusAPI.create({
        projectId: project.data!.id,
      } as any);

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create status with empty name", async ({
      projectAPI,
      statusAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const response = await statusAPI.create({
        projectId: project.data!.id,
        name: "",
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create status with invalid projectId", async ({
      statusAPI,
    }) => {
      const response = await statusAPI.create({
        projectId: "invalid-uuid",
        name: "Test Status",
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to create status for non-existent project", async ({
      statusAPI,
    }) => {
      const fakeProjectId = "00000000-0000-0000-0000-000000000000";
      const response = await statusAPI.create({
        projectId: fakeProjectId,
        name: "Test Status",
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("GET /statuses", () => {
    test("should get statuses by project ID", async ({
      projectAPI,
      statusAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      // Create multiple statuses
      const status1 = await statusAPI.create({
        projectId: project.data!.id,
        name: "To Do",
      });
      const status2 = await statusAPI.create({
        projectId: project.data!.id,
        name: "In Progress",
      });
      createdStatusIds.push(status1.data!.id, status2.data!.id);

      const response = await statusAPI.getByProject(project.data!.id);

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(Array.isArray(response.data)).toBe(true);
      // 3 auto-created + 2 manually created = 5 total
      expect(response.data!.length).toBeGreaterThanOrEqual(5);

      // Verify all statuses belong to the project
      const allBelongToProject = response.data!.every(
        (s) => s.projectId === project.data!.id
      );
      expect(allBelongToProject).toBe(true);
    });

    test("should return statuses ordered by position", async ({
      projectAPI,
      statusAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      // Create multiple statuses
      const status1 = await statusAPI.create({
        projectId: project.data!.id,
        name: "First",
      });
      const status2 = await statusAPI.create({
        projectId: project.data!.id,
        name: "Second",
      });
      const status3 = await statusAPI.create({
        projectId: project.data!.id,
        name: "Third",
      });
      createdStatusIds.push(
        status1.data!.id,
        status2.data!.id,
        status3.data!.id
      );

      const response = await statusAPI.getByProject(project.data!.id);

      expect(response.status).toBe(200);
      // Verify ascending position order
      for (let i = 0; i < response.data!.length - 1; i++) {
        expect(response.data![i].position).toBeLessThanOrEqual(
          response.data![i + 1].position
        );
      }
    });

    test("should return default statuses for new project", async ({
      projectAPI,
      statusAPI,
    }) => {
      const project = await projectAPI.create({
        name: "New Project",
        description: "Has default statuses",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const response = await statusAPI.getByProject(project.data!.id);

      expect(response.status).toBe(200);
      expect(Array.isArray(response.data)).toBe(true);
      // Should have 3 auto-created statuses
      expect(response.data).toHaveLength(3);
      expect(response.data?.map((s) => s.name)).toEqual([
        "Todo",
        "In Progress",
        "Done",
      ]);
    });

    test("should fail to get statuses without projectId", async ({
      statusAPI,
    }) => {
      const response = await statusAPI.getByProject("");

      expect(response.status).toBeGreaterThanOrEqual(400);
    });
  });

  test.describe("GET /statuses/{statusId}", () => {
    test("should get status by ID", async ({ projectAPI, statusAPI }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const created = await statusAPI.create({
        projectId: project.data!.id,
        name: "Test Status",
      });
      createdStatusIds.push(created.data!.id);

      const response = await statusAPI.getById(created.data!.id);

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.id).toBe(created.data!.id);
      expect(response.data?.name).toBe("Test Status");
      expect(response.data?.projectId).toBe(project.data!.id);
    });

    test("should return 404 for non-existent status", async ({ statusAPI }) => {
      const fakeId = "00000000-0000-0000-0000-000000000000";
      const response = await statusAPI.getById(fakeId);

      expect(response.status).toBe(404);
      expect(response.error).toBeDefined();
    });

    test("should return 400 for invalid UUID format", async ({ statusAPI }) => {
      const response = await statusAPI.getById("invalid-uuid");

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("PATCH /statuses/{statusId}", () => {
    test("should update status name", async ({ projectAPI, statusAPI }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const created = await statusAPI.create({
        projectId: project.data!.id,
        name: "Original Name",
      });
      createdStatusIds.push(created.data!.id);

      const response = await statusAPI.update(created.data!.id, {
        name: "Updated Name",
      });

      expect(response.status).toBe(200);
      expect(response.data?.name).toBe("Updated Name");
      expect(response.data?.slug).toBeDefined();
      expect(response.data?.updatedAt).not.toBe(created.data!.updatedAt);
    });

    test("should fail to update with empty name", async ({
      projectAPI,
      statusAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const created = await statusAPI.create({
        projectId: project.data!.id,
        name: "Test Status",
      });
      createdStatusIds.push(created.data!.id);

      const response = await statusAPI.update(created.data!.id, {
        name: "",
      });

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail to update non-existent status", async ({ statusAPI }) => {
      const fakeId = "00000000-0000-0000-0000-000000000000";
      const response = await statusAPI.update(fakeId, {
        name: "New Name",
      });

      expect(response.status).toBe(404);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("DELETE /statuses/{statusId}", () => {
    test("should delete status", async ({ projectAPI, statusAPI }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const created = await statusAPI.create({
        projectId: project.data!.id,
        name: "Status to Delete",
      });

      const deleteResponse = await statusAPI.remove(created.data!.id);

      expect(deleteResponse.status).toBe(204);

      // Verify status is deleted
      const getResponse = await statusAPI.getById(created.data!.id);
      expect(getResponse.status).toBe(404);
    });

    test("should return 404 when deleting non-existent status", async ({
      statusAPI,
    }) => {
      const fakeId = "00000000-0000-0000-0000-000000000000";
      const response = await statusAPI.remove(fakeId);

      expect(response.status).toBe(404);
    });

    test("should be idempotent - deleting twice", async ({
      projectAPI,
      statusAPI,
    }) => {
      const project = await projectAPI.create({
        name: "Test Project",
        description: "Test",
        status: "active",
      });
      createdProjectIds.push(project.data!.id);

      const created = await statusAPI.create({
        projectId: project.data!.id,
        name: "Status to Delete",
      });

      // First delete
      const firstDelete = await statusAPI.remove(created.data!.id);
      expect(firstDelete.status).toBe(204);

      // Second delete
      const secondDelete = await statusAPI.remove(created.data!.id);
      expect(secondDelete.status).toBe(404);
    });
  });
});
