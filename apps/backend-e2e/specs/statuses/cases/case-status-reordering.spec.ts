import { test, expect } from "../../../fixtures";

/**
 * Case: Status Reordering
 * Tests drag-and-drop reordering functionality for statuses
 */
test.describe("Case: Status Reordering", () => {
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

  test("should reorder statuses maintaining correct positions", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Kanban Project",
      description: "Testing reorder",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Create statuses
    const todo = await statusAPI.create({
      projectId: project.data!.id,
      name: "To Do",
    });
    const inProgress = await statusAPI.create({
      projectId: project.data!.id,
      name: "In Progress",
    });
    const done = await statusAPI.create({
      projectId: project.data!.id,
      name: "Done",
    });
    createdStatusIds.push(todo.data!.id, inProgress.data!.id, done.data!.id);

    // Get all statuses (including 3 auto-created)
    const initialStatuses = await statusAPI.getByProject(project.data!.id);
    const allStatusIds = initialStatuses.data!.map((s) => s.id);

    // Find our manually created statuses
    const todoStatus = initialStatuses.data!.find((s) => s.name === "To Do");
    const inProgressStatus = initialStatuses.data!.find(
      (s) => s.name === "In Progress"
    );
    const doneStatus = initialStatuses.data!.find((s) => s.name === "Done");

    // Reorder ALL statuses (backend requires all IDs)
    // Put manually created ones first: [To Do, In Progress, Done, auto-created...]
    const autoCreatedIds = allStatusIds.filter(
      (id) => ![todo.data!.id, inProgress.data!.id, done.data!.id].includes(id)
    );

    const reorderResponse = await statusAPI.reorder({
      projectId: project.data!.id,
      ids: [
        done.data!.id,
        todo.data!.id,
        inProgress.data!.id,
        ...autoCreatedIds,
      ],
    });

    expect(reorderResponse.status).toBe(200);
    expect(reorderResponse.data).toBeDefined();

    // Verify new order (manually created should be reordered)
    const reorderedStatuses = await statusAPI.getByProject(project.data!.id);
    const manualStatuses = reorderedStatuses.data!.filter((s) =>
      [todo.data!.id, inProgress.data!.id, done.data!.id].includes(s.id)
    );
    expect(manualStatuses.map((s) => s.name)).toEqual([
      "Done",
      "To Do",
      "In Progress",
    ]);

    // Verify positions are sequential
    for (let i = 0; i < reorderedStatuses.data!.length - 1; i++) {
      expect(reorderedStatuses.data![i].position).toBeLessThan(
        reorderedStatuses.data![i + 1].position
      );
    }
  });

  test("should handle reordering with all status IDs", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Reorder Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const status = await statusAPI.create({
      projectId: project.data!.id,
      name: "Custom Status",
    });
    createdStatusIds.push(status.data!.id);

    // Get all statuses and include all IDs in reorder
    const allStatuses = await statusAPI.getByProject(project.data!.id);
    const allIds = allStatuses.data!.map((s) => s.id);

    const response = await statusAPI.reorder({
      projectId: project.data!.id,
      ids: allIds,
    });

    expect(response.status).toBe(200);
    expect(response.data!.length).toBeGreaterThanOrEqual(4); // 3 auto + 1 manual
  });

  test("should handle multiple reorders", async ({ projectAPI, statusAPI }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Multiple reorders",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

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

    // Get all statuses to include auto-created ones
    let allStatuses = await statusAPI.getByProject(project.data!.id);
    const autoCreatedIds = allStatuses
      .data!.filter(
        (s) =>
          ![status1.data!.id, status2.data!.id, status3.data!.id].includes(s.id)
      )
      .map((s) => s.id);

    // First reorder: [1, 2, 3] -> [3, 2, 1] (+ auto-created at end)
    await statusAPI.reorder({
      projectId: project.data!.id,
      ids: [
        status3.data!.id,
        status2.data!.id,
        status1.data!.id,
        ...autoCreatedIds,
      ],
    });

    let statuses = await statusAPI.getByProject(project.data!.id);
    let manualStatuses = statuses.data!.filter((s) =>
      [status1.data!.id, status2.data!.id, status3.data!.id].includes(s.id)
    );
    expect(manualStatuses.map((s) => s.name)).toEqual([
      "Status 3",
      "Status 2",
      "Status 1",
    ]);

    // Second reorder: [3, 2, 1] -> [2, 1, 3] (+ auto-created at end)
    await statusAPI.reorder({
      projectId: project.data!.id,
      ids: [
        status2.data!.id,
        status1.data!.id,
        status3.data!.id,
        ...autoCreatedIds,
      ],
    });

    statuses = await statusAPI.getByProject(project.data!.id);
    manualStatuses = statuses.data!.filter((s) =>
      [status1.data!.id, status2.data!.id, status3.data!.id].includes(s.id)
    );
    expect(manualStatuses.map((s) => s.name)).toEqual([
      "Status 2",
      "Status 1",
      "Status 3",
    ]);
  });

  test("should fail to reorder with invalid status IDs", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const response = await statusAPI.reorder({
      projectId: project.data!.id,
      ids: ["invalid-uuid-1", "invalid-uuid-2"],
    });

    expect(response.status).toBeGreaterThanOrEqual(400);
    expect(response.error).toBeDefined();
  });

  test("should fail to reorder with missing projectId", async ({
    statusAPI,
  }) => {
    const response = await statusAPI.reorder({
      ids: ["some-id"],
    } as any);

    expect(response.status).toBeGreaterThanOrEqual(400);
    expect(response.error).toBeDefined();
  });

  test("should fail to reorder with non-existent project", async ({
    statusAPI,
  }) => {
    const fakeProjectId = "00000000-0000-0000-0000-000000000000";
    const response = await statusAPI.reorder({
      projectId: fakeProjectId,
      ids: ["some-id"],
    });

    expect(response.status).toBeGreaterThanOrEqual(400);
    expect(response.error).toBeDefined();
  });

  test("should fail to reorder statuses from different projects", async ({
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

    // Create statuses in different projects
    const status1 = await statusAPI.create({
      projectId: project1.data!.id,
      name: "Status 1",
    });
    const status2 = await statusAPI.create({
      projectId: project2.data!.id,
      name: "Status 2",
    });
    createdStatusIds.push(status1.data!.id, status2.data!.id);

    // Try to reorder statuses from project2 using status from project1
    const response = await statusAPI.reorder({
      projectId: project2.data!.id,
      ids: [status1.data!.id, status2.data!.id], // status1 belongs to project1
    });

    expect(response.status).toBeGreaterThanOrEqual(400);
    expect(response.error).toBeDefined();
  });

  test("should fail to reorder with incomplete list", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Empty list should fail because backend expects all status IDs
    const response = await statusAPI.reorder({
      projectId: project.data!.id,
      ids: [],
    });

    expect(response.status).toBeGreaterThanOrEqual(400);
  });
});
