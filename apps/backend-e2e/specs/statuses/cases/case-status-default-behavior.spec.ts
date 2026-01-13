import { test, expect } from "../../../fixtures";

/**
 * Case: Default Status Behavior
 * Tests behavior related to default status designation
 */
test.describe("Case: Default Status Behavior", () => {
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

  test("should have auto-created Todo as default status", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test default status",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    const defaultStatus = statuses.data!.find((s) => s.isDefault);

    expect(defaultStatus).toBeDefined();
    expect(defaultStatus?.name).toBe("Todo");
    expect(defaultStatus?.isDefault).toBe(true);
  });

  test("should not mark manually created statuses as default", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const firstStatus = await statusAPI.create({
      projectId: project.data!.id,
      name: "Custom First",
    });
    const secondStatus = await statusAPI.create({
      projectId: project.data!.id,
      name: "Custom Second",
    });
    const thirdStatus = await statusAPI.create({
      projectId: project.data!.id,
      name: "Custom Third",
    });
    createdStatusIds.push(
      firstStatus.data!.id,
      secondStatus.data!.id,
      thirdStatus.data!.id
    );

    // Manually created statuses should not be default
    expect(firstStatus.data?.isDefault).toBe(false);
    expect(secondStatus.data?.isDefault).toBe(false);
    expect(thirdStatus.data?.isDefault).toBe(false);
  });

  test("should have only one default status per project", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Create multiple statuses (in addition to 3 auto-created)
    for (let i = 1; i <= 5; i++) {
      const status = await statusAPI.create({
        projectId: project.data!.id,
        name: `Status ${i}`,
      });
      createdStatusIds.push(status.data!.id);
    }

    const statuses = await statusAPI.getByProject(project.data!.id);
    const defaultStatuses = statuses.data!.filter((s) => s.isDefault);

    expect(defaultStatuses).toHaveLength(1);
    // The auto-created "Todo" status is the default
    expect(defaultStatuses[0].name).toBe("Todo");
  });

  test("should maintain default status across different projects", async ({
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

    // Get auto-created statuses for each project
    const p1Statuses = await statusAPI.getByProject(project1.data!.id);
    const p2Statuses = await statusAPI.getByProject(project2.data!.id);

    // Each project should have its own default "Todo" status
    const p1Default = p1Statuses.data!.find((s) => s.isDefault);
    const p2Default = p2Statuses.data!.find((s) => s.isDefault);

    expect(p1Default).toBeDefined();
    expect(p2Default).toBeDefined();
    expect(p1Default?.name).toBe("Todo");
    expect(p2Default?.name).toBe("Todo");

    // Verify in list
    const project1Statuses = await statusAPI.getByProject(project1.data!.id);
    const project2Statuses = await statusAPI.getByProject(project2.data!.id);

    expect(project1Statuses.data!.filter((s) => s.isDefault)).toHaveLength(1);
    expect(project2Statuses.data!.filter((s) => s.isDefault)).toHaveLength(1);
  });

  test("should not change default status when creating new statuses", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Get the auto-created default status (Todo)
    const initialStatuses = await statusAPI.getByProject(project.data!.id);
    const defaultStatus = initialStatuses.data!.find((s) => s.isDefault);
    const defaultId = defaultStatus!.id;

    // Create more statuses
    for (let i = 1; i <= 5; i++) {
      const status = await statusAPI.create({
        projectId: project.data!.id,
        name: `Status ${i}`,
      });
      createdStatusIds.push(status.data!.id);
    }

    // Verify Todo is still the default
    const stillDefault = await statusAPI.getById(defaultId);
    expect(stillDefault.data?.isDefault).toBe(true);
    expect(stillDefault.data?.name).toBe("Todo");
  });
});
