import { test, expect } from "../../../fixtures";

/**
 * Case: Project Name Uniqueness
 * Tests scenarios related to duplicate project names and slug generation
 */
test.describe("Case: Project Name Uniqueness", () => {
  const createdProjectIds: string[] = [];

  test.afterEach(async ({ projectAPI }) => {
    for (const id of createdProjectIds) {
      await projectAPI.remove(id).catch(() => {});
    }
    createdProjectIds.length = 0;
  });

  test("should allow creating projects with different names", async ({
    projectAPI,
  }) => {
    const project1 = await projectAPI.create({
      name: "Project One",
      description: "First project",
      status: "active",
    });
    createdProjectIds.push(project1.data!.id);

    const project2 = await projectAPI.create({
      name: "Project Two",
      description: "Second project",
      status: "active",
    });
    createdProjectIds.push(project2.data!.id);

    expect(project1.status).toBe(200);
    expect(project2.status).toBe(200);
    expect(project1.data?.id).not.toBe(project2.data?.id);
  });

  test("should handle projects with similar names", async ({ projectAPI }) => {
    const timestamp = Date.now();

    const project1 = await projectAPI.create({
      name: `Test Project ${timestamp}`,
      description: "First",
      status: "active",
    });
    createdProjectIds.push(project1.data!.id);

    const project2 = await projectAPI.create({
      name: `Test Project ${timestamp} Copy`,
      description: "Second",
      status: "active",
    });
    createdProjectIds.push(project2.data!.id);

    expect(project1.status).toBe(200);
    expect(project2.status).toBe(200);
  });

  test("should preserve exact name casing", async ({ projectAPI }) => {
    const created = await projectAPI.create({
      name: "CamelCase Project Name",
      description: "Testing case preservation",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    expect(created.data?.name).toBe("CamelCase Project Name");

    // Verify after retrieval
    const fetched = await projectAPI.getById(created.data!.id);
    expect(fetched.data?.name).toBe("CamelCase Project Name");
  });

  test("should handle special characters in project name", async ({
    projectAPI,
  }) => {
    const created = await projectAPI.create({
      name: "Project with Special Chars: @#$%",
      description: "Testing special characters",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    expect(created.status).toBe(200);
    expect(created.data?.name).toBe("Project with Special Chars: @#$%");
  });

  test("should handle unicode characters in project name", async ({
    projectAPI,
  }) => {
    const created = await projectAPI.create({
      name: "Проект Unicode 项目 🚀",
      description: "Testing unicode",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    expect(created.status).toBe(200);
    expect(created.data?.name).toBe("Проект Unicode 项目 🚀");
  });

  test("should allow updating project name", async ({ projectAPI }) => {
    const created = await projectAPI.create({
      name: "Original Name",
      description: "Will be renamed",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    const updated = await projectAPI.update(created.data!.id, {
      name: "New Name",
      description: "Will be renamed",
      status: "active",
    });

    expect(updated.status).toBe(200);
    expect(updated.data?.name).toBe("New Name");
  });

  test("should handle name updates with same content", async ({
    projectAPI,
  }) => {
    const created = await projectAPI.create({
      name: "Unchanged Name",
      description: "Description",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    // Update with same name
    const updated = await projectAPI.update(created.data!.id, {
      name: "Unchanged Name",
      description: "Description",
      status: "active",
    });

    expect(updated.status).toBe(200);
    expect(updated.data?.name).toBe("Unchanged Name");
  });
});
