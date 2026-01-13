import { test, expect } from "../../../fixtures";

/**
 * Case: Project Status Transitions
 * Tests business rules around project status changes
 */
test.describe("Case: Project Status Transitions", () => {
  const createdProjectIds: string[] = [];

  test.afterEach(async ({ projectAPI }) => {
    for (const id of createdProjectIds) {
      await projectAPI.remove(id).catch(() => {});
    }
    createdProjectIds.length = 0;
  });

  test("should transition project from active to paused", async ({
    projectAPI,
  }) => {
    // Create active project
    const created = await projectAPI.create({
      name: "Active Project",
      description: "Will be paused",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    expect(created.data?.status).toBe("active");

    // Pause the project
    const paused = await projectAPI.update(created.data!.id, {
      status: "paused",
    });

    expect(paused.status).toBe(200);
    expect(paused.data?.status).toBe("paused");

    // Verify the change persisted
    const fetched = await projectAPI.getById(created.data!.id);
    expect(fetched.data?.status).toBe("paused");
  });

  test("should transition project from paused to active", async ({
    projectAPI,
  }) => {
    // Create paused project
    const created = await projectAPI.create({
      name: "Paused Project",
      description: "Will be activated",
      status: "paused",
    });
    createdProjectIds.push(created.data!.id);

    // Activate the project
    const activated = await projectAPI.update(created.data!.id, {
      status: "active",
    });

    expect(activated.status).toBe(200);
    expect(activated.data?.status).toBe("active");
  });

  test("should transition project from active to archived", async ({
    projectAPI,
  }) => {
    // Create active project
    const created = await projectAPI.create({
      name: "Active Project",
      description: "Will be archived",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    // Archive the project
    const archived = await projectAPI.update(created.data!.id, {
      status: "archived",
    });

    expect(archived.status).toBe(200);
    expect(archived.data?.status).toBe("archived");

    // Verify archived projects are still retrievable
    const fetched = await projectAPI.getById(created.data!.id);
    expect(fetched.status).toBe(200);
    expect(fetched.data?.status).toBe("archived");
  });

  test("should allow unarchiving a project", async ({ projectAPI }) => {
    // Create archived project
    const created = await projectAPI.create({
      name: "Archived Project",
      description: "Will be unarchived",
      status: "archived",
    });
    createdProjectIds.push(created.data!.id);

    // Unarchive to active
    const unarchived = await projectAPI.update(created.data!.id, {
      status: "active",
    });

    expect(unarchived.status).toBe(200);
    expect(unarchived.data?.status).toBe("active");
  });

  test("should handle multiple status transitions", async ({ projectAPI }) => {
    const created = await projectAPI.create({
      name: "Multi-transition Project",
      description: "Testing multiple transitions",
      status: "active",
    });
    createdProjectIds.push(created.data!.id);

    // Active -> Paused
    const paused = await projectAPI.update(created.data!.id, {
      status: "paused",
    });
    expect(paused.data?.status).toBe("paused");

    // Paused -> Active
    const reactivated = await projectAPI.update(created.data!.id, {
      status: "active",
    });
    expect(reactivated.data?.status).toBe("active");

    // Active -> Archived
    const archived = await projectAPI.update(created.data!.id, {
      status: "archived",
    });
    expect(archived.data?.status).toBe("archived");

    // Archived -> Active
    const finalActive = await projectAPI.update(created.data!.id, {
      status: "active",
    });
    expect(finalActive.data?.status).toBe("active");
  });

  test("should filter archived projects separately", async ({ projectAPI }) => {
    // Create projects with different statuses
    const active = await projectAPI.create({
      name: "Active for Filter",
      description: "Active",
      status: "active",
    });
    const archived = await projectAPI.create({
      name: "Archived for Filter",
      description: "Archived",
      status: "archived",
    });
    createdProjectIds.push(active.data!.id, archived.data!.id);

    // Get only archived projects
    const archivedList = await projectAPI.getPaginated({
      status: ["archived"],
    });

    expect(archivedList.status).toBe(200);
    const allArchived = archivedList.data!.items!.every(
      (p) => p.status === "archived"
    );
    expect(allArchived).toBe(true);

    // Get only active projects
    const activeList = await projectAPI.getPaginated({
      status: ["active"],
    });

    expect(activeList.status).toBe(200);
    const allActive = activeList.data!.items!.every(
      (p) => p.status === "active"
    );
    expect(allActive).toBe(true);
  });
});
