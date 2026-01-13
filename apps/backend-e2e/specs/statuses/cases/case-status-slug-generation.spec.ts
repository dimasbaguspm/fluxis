import { test, expect } from "../../../fixtures";

/**
 * Case: Status Slug Generation
 * Tests automatic slug generation from status names
 */
test.describe("Case: Status Slug Generation", () => {
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

  test("should generate slug from simple name", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const status = await statusAPI.create({
      projectId: project.data!.id,
      name: "Todo",
    });
    createdStatusIds.push(status.data!.id);

    expect(status.data?.slug).toBeDefined();
    expect(status.data?.slug).toMatch(/^[a-z0-9_]+$/);
    expect(status.data?.slug.toLowerCase()).toContain("todo");
  });

  test("should generate slug from multi-word name", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const status = await statusAPI.create({
      projectId: project.data!.id,
      name: "In Progress",
    });
    createdStatusIds.push(status.data!.id);

    expect(status.data?.slug).toBeDefined();
    expect(status.data?.slug).toMatch(/^[a-z0-9_]+$/);
  });

  test("should handle special characters in name", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const status = await statusAPI.create({
      projectId: project.data!.id,
      name: "Ready to Deploy!",
    });
    createdStatusIds.push(status.data!.id);

    expect(status.data?.slug).toBeDefined();
    expect(status.data?.slug).toMatch(/^[a-z0-9_]+$/);
    // Slug should not contain special characters
    expect(status.data?.slug).not.toContain("!");
  });

  test("should handle numbers in name", async ({ projectAPI, statusAPI }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const status = await statusAPI.create({
      projectId: project.data!.id,
      name: "Phase 2 Testing",
    });
    createdStatusIds.push(status.data!.id);

    expect(status.data?.slug).toBeDefined();
    expect(status.data?.slug).toMatch(/^[a-z0-9_]+$/);
    expect(status.data?.slug).toContain("2");
  });

  test("should generate unique slugs for similar names", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const status1 = await statusAPI.create({
      projectId: project.data!.id,
      name: "Testing",
    });
    const status2 = await statusAPI.create({
      projectId: project.data!.id,
      name: "Testing Phase",
    });
    createdStatusIds.push(status1.data!.id, status2.data!.id);

    expect(status1.data?.slug).toBeDefined();
    expect(status2.data?.slug).toBeDefined();
    // Slugs should be different
    expect(status1.data?.slug).not.toBe(status2.data?.slug);
  });

  test("should update slug when name is updated", async ({
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
      name: "Original Name",
    });
    createdStatusIds.push(created.data!.id);

    const originalSlug = created.data!.slug;

    const updated = await statusAPI.update(created.data!.id, {
      name: "Updated Name",
    });

    expect(updated.data?.slug).toBeDefined();
    expect(updated.data?.slug).not.toBe(originalSlug);
  });

  test("should handle emoji in name", async ({ projectAPI, statusAPI }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const status = await statusAPI.create({
      projectId: project.data!.id,
      name: "Done ✅",
    });
    createdStatusIds.push(status.data!.id);

    expect(status.data?.slug).toBeDefined();
    expect(status.data?.slug).toMatch(/^[a-z0-9_]+$/);
    // Slug should not contain emoji
    expect(status.data?.slug).not.toContain("✅");
  });

  test("should allow same slug in different projects", async ({
    projectAPI,
    statusAPI,
  }) => {
    // Create two projects
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

    // Create statuses with same name in different projects
    const status1 = await statusAPI.create({
      projectId: project1.data!.id,
      name: "To Do",
    });
    const status2 = await statusAPI.create({
      projectId: project2.data!.id,
      name: "To Do",
    });
    createdStatusIds.push(status1.data!.id, status2.data!.id);

    // Slugs can be the same across different projects
    expect(status1.data?.slug).toBeDefined();
    expect(status2.data?.slug).toBeDefined();
  });
});
