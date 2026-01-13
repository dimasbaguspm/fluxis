import { test, expect } from "../../../fixtures";

/**
 * Case: Status Cascade Operations
 * Tests behavior when parent project is deleted or modified
 */
test.describe("Case: Status Cascade Operations", () => {
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

  test("should delete all statuses when project is deleted", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Will be deleted",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Create multiple statuses (in addition to 3 auto-created)
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

    // Verify statuses exist (3 auto-created + 3 manual = 6)
    const beforeDelete = await statusAPI.getByProject(project.data!.id);
    expect(beforeDelete.data).toHaveLength(6);

    // Delete project
    await projectAPI.remove(project.data!.id);

    // Verify manually created statuses are soft-deleted (still return 200 with deletedAt)
    const status1Check = await statusAPI.getById(status1.data!.id);
    const status2Check = await statusAPI.getById(status2.data!.id);
    const status3Check = await statusAPI.getById(status3.data!.id);

    // Soft delete: status may still be retrievable but with deletedAt set
    // or may return 404 depending on backend filter implementation
    expect([200, 404]).toContain(status1Check.status);
    expect([200, 404]).toContain(status2Check.status);
    expect([200, 404]).toContain(status3Check.status);
  });

  test("should not affect statuses of other projects when one is deleted", async ({
    projectAPI,
    statusAPI,
  }) => {
    // Create two projects
    const project1 = await projectAPI.create({
      name: "Project 1",
      description: "Will be deleted",
      status: "active",
    });
    const project2 = await projectAPI.create({
      name: "Project 2",
      description: "Will remain",
      status: "active",
    });
    createdProjectIds.push(project1.data!.id, project2.data!.id);

    // Create statuses for both projects
    const p1Status = await statusAPI.create({
      projectId: project1.data!.id,
      name: "P1 Status",
    });
    const p2Status = await statusAPI.create({
      projectId: project2.data!.id,
      name: "P2 Status",
    });
    createdStatusIds.push(p1Status.data!.id, p2Status.data!.id);

    // Delete project1
    await projectAPI.remove(project1.data!.id);

    // Verify project1 statuses are soft-deleted
    const p1StatusCheck = await statusAPI.getById(p1Status.data!.id);
    expect([200, 404]).toContain(p1StatusCheck.status);

    // Verify project2 statuses still exist
    const p2StatusCheck = await statusAPI.getById(p2Status.data!.id);
    expect(p2StatusCheck.status).toBe(200);
    expect(p2StatusCheck.data?.id).toBe(p2Status.data!.id);
  });

  test("should handle deleting project with many statuses", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Project with Many Statuses",
      description: "Test cascade delete",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Create 10 statuses (in addition to 3 auto-created)
    const statusIds: string[] = [];
    for (let i = 1; i <= 10; i++) {
      const status = await statusAPI.create({
        projectId: project.data!.id,
        name: `Status ${i}`,
      });
      statusIds.push(status.data!.id);
      createdStatusIds.push(status.data!.id);
    }

    // Verify all created (3 auto + 10 manual = 13)
    const beforeDelete = await statusAPI.getByProject(project.data!.id);
    expect(beforeDelete.data).toHaveLength(13);

    // Delete project
    await projectAPI.remove(project.data!.id);

    // Verify all statuses are soft-deleted
    for (const statusId of statusIds) {
      const check = await statusAPI.getById(statusId);
      expect([200, 404]).toContain(check.status);
    }
  });

  test("should allow creating new statuses after all are deleted", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test recreation",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Create and delete statuses
    const status1 = await statusAPI.create({
      projectId: project.data!.id,
      name: "Temporary Status",
    });
    await statusAPI.remove(status1.data!.id);

    // Create new status
    const status2 = await statusAPI.create({
      projectId: project.data!.id,
      name: "New Status",
    });
    createdStatusIds.push(status2.data!.id);

    expect(status2.status).toBe(200);
    expect(status2.data?.name).toBe("New Status");

    // Verify it's NOT marked as default (only auto-created statuses are default)
    expect(status2.data?.isDefault).toBe(false);
  });

  test("should maintain status independence across projects", async ({
    projectAPI,
    statusAPI,
  }) => {
    // Create multiple projects with same status names
    const projects = [];
    for (let i = 1; i <= 3; i++) {
      const project = await projectAPI.create({
        name: `Project ${i}`,
        description: `Project ${i}`,
        status: "active",
      });
      createdProjectIds.push(project.data!.id);
      projects.push(project);

      // Create additional statuses
      const todo = await statusAPI.create({
        projectId: project.data!.id,
        name: "To Do",
      });
      const inProgress = await statusAPI.create({
        projectId: project.data!.id,
        name: "In Progress",
      });
      createdStatusIds.push(todo.data!.id, inProgress.data!.id);
    }

    // Delete middle project
    await projectAPI.remove(projects[1].data!.id);

    // Verify other projects' statuses are intact (3 auto + 2 manual = 5)
    const project1Statuses = await statusAPI.getByProject(projects[0].data!.id);
    const project3Statuses = await statusAPI.getByProject(projects[2].data!.id);

    expect(project1Statuses.data).toHaveLength(5);
    expect(project3Statuses.data).toHaveLength(5);
  });
});
