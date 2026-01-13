import { test, expect } from "../../../fixtures";

/**
 * Status Advanced Filtering Test Cases
 * Tests complex filtering for statuses (by project mainly)
 */
test.describe("Status Advanced Filtering", () => {
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

  test("should filter statuses by project correctly", async ({
    projectAPI,
    statusAPI,
  }) => {
    // Create two projects
    const project1 = await projectAPI.create({
      name: "Project 1",
      description: "First",
      status: "active",
    });

    const project2 = await projectAPI.create({
      name: "Project 2",
      description: "Second",
      status: "active",
    });

    createdProjectIds.push(project1.data!.id, project2.data!.id);

    const statuses1 = await statusAPI.getByProject(project1.data!.id);
    const statuses2 = await statusAPI.getByProject(project2.data!.id);

    createdStatusIds.push(...statuses1.data!.map((s) => s.id));
    createdStatusIds.push(...statuses2.data!.map((s) => s.id));

    // Statuses should be completely separate
    expect(statuses1.data?.length).toBeGreaterThan(0);
    expect(statuses2.data?.length).toBeGreaterThan(0);

    const ids1 = statuses1.data!.map((s) => s.id);
    const ids2 = statuses2.data!.map((s) => s.id);

    // No overlap
    const overlap = ids1.filter((id) => ids2.includes(id));
    expect(overlap.length).toBe(0);
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

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create multiple custom statuses
    for (let i = 0; i < 5; i++) {
      const status = await statusAPI.create({
        projectId: project.data!.id,
        name: `Status ${i}`,
      });
      createdStatusIds.push(status.data!.id);
    }

    const allStatuses = await statusAPI.getByProject(project.data!.id);

    expect(allStatuses.status).toBe(200);

    // Check positions are sequential
    const positions = allStatuses.data!.map((s) => s.position);
    for (let i = 1; i < positions.length; i++) {
      expect(positions[i]).toBeGreaterThan(positions[i - 1]);
    }
  });

  test("should filter by default status", async ({ projectAPI, statusAPI }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const statuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...statuses.data!.map((s) => s.id));

    // Create additional non-default statuses
    const custom1 = await statusAPI.create({
      projectId: project.data!.id,
      name: "Custom 1",
    });
    const custom2 = await statusAPI.create({
      projectId: project.data!.id,
      name: "Custom 2",
    });

    createdStatusIds.push(custom1.data!.id, custom2.data!.id);

    // Get all statuses
    const allStatuses = await statusAPI.getByProject(project.data!.id);

    // Should have exactly one default status
    const defaultStatuses = allStatuses.data!.filter((s) => s.isDefault);
    expect(defaultStatuses.length).toBe(1);

    // Non-default should be majority
    const nonDefaultStatuses = allStatuses.data!.filter((s) => !s.isDefault);
    expect(nonDefaultStatuses.length).toBeGreaterThan(1);
  });

  test.skip("should handle slug-based filtering/searching", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create status with specific name (will generate slug)
    const status = await statusAPI.create({
      projectId: project.data!.id,
      name: "In Review Process",
    });
    createdStatusIds.push(status.data!.id);

    expect(status.data?.slug).toBe("in_review_process");

    // Get all statuses and verify slug exists
    const allStatuses = await statusAPI.getByProject(project.data!.id);
    const reviewStatus = allStatuses.data!.find(
      (s) => s.slug === "in-review-process"
    );

    expect(reviewStatus).toBeDefined();
    expect(reviewStatus?.name).toBe("In Review Process");
  });

  test("should handle statuses with same name in different projects", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project1 = await projectAPI.create({
      name: "Project 1",
      description: "First",
      status: "active",
    });

    const project2 = await projectAPI.create({
      name: "Project 2",
      description: "Second",
      status: "active",
    });

    createdProjectIds.push(project1.data!.id, project2.data!.id);

    const defaultStatuses1 = await statusAPI.getByProject(project1.data!.id);
    const defaultStatuses2 = await statusAPI.getByProject(project2.data!.id);

    createdStatusIds.push(...defaultStatuses1.data!.map((s) => s.id));
    createdStatusIds.push(...defaultStatuses2.data!.map((s) => s.id));

    // Create statuses with same name in both projects
    const status1 = await statusAPI.create({
      projectId: project1.data!.id,
      name: "Testing",
    });

    const status2 = await statusAPI.create({
      projectId: project2.data!.id,
      name: "Testing",
    });

    createdStatusIds.push(status1.data!.id, status2.data!.id);

    // Both should have same name and slug
    expect(status1.data?.name).toBe("Testing");
    expect(status2.data?.name).toBe("Testing");
    expect(status1.data?.slug).toBe("testing");
    expect(status2.data?.slug).toBe("testing");

    // But different IDs
    expect(status1.data?.id).not.toBe(status2.data?.id);

    // Verify they're in different projects
    const statuses1 = await statusAPI.getByProject(project1.data!.id);
    const statuses2 = await statusAPI.getByProject(project2.data!.id);

    expect(statuses1.data!.some((s) => s.id === status1.data!.id)).toBe(true);
    expect(statuses1.data!.some((s) => s.id === status2.data!.id)).toBe(false);

    expect(statuses2.data!.some((s) => s.id === status2.data!.id)).toBe(true);
    expect(statuses2.data!.some((s) => s.id === status1.data!.id)).toBe(false);
  });

  test("should handle position gaps after deletions", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create several statuses
    const statuses = [];
    for (let i = 0; i < 5; i++) {
      const status = await statusAPI.create({
        projectId: project.data!.id,
        name: `Status ${i}`,
      });
      statuses.push(status.data!);
      createdStatusIds.push(status.data!.id);
    }

    // Delete middle ones
    await statusAPI.remove(statuses[1].id);
    await statusAPI.remove(statuses[3].id);

    // Get remaining statuses
    const remaining = await statusAPI.getByProject(project.data!.id);

    expect(remaining.status).toBe(200);

    // Positions should still be in ascending order
    const positions = remaining.data!.map((s) => s.position);
    for (let i = 1; i < positions.length; i++) {
      expect(positions[i]).toBeGreaterThan(positions[i - 1]);
    }
  });

  test("should handle unicode characters in status names", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create status with unicode
    const status = await statusAPI.create({
      projectId: project.data!.id,
      name: "完了 (Completed)",
    });
    createdStatusIds.push(status.data!.id);

    expect(status.status).toBe(200);
    expect(status.data?.name).toBe("完了 (Completed)");

    // Verify it appears in the list
    const allStatuses = await statusAPI.getByProject(project.data!.id);
    expect(allStatuses.data!.some((s) => s.id === status.data!.id)).toBe(true);
  });

  test("should handle emoji in status names", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create status with emoji
    const status = await statusAPI.create({
      projectId: project.data!.id,
      name: "🎉 Completed",
    });
    createdStatusIds.push(status.data!.id);

    expect(status.status).toBe(200);
    expect(status.data?.name).toBe("🎉 Completed");

    // Slug should handle emoji appropriately
    expect(status.data?.slug).toBeDefined();
  });

  test("should maintain consistency across concurrent reads", async ({
    projectAPI,
    statusAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const defaultStatuses = await statusAPI.getByProject(project.data!.id);
    createdStatusIds.push(...defaultStatuses.data!.map((s) => s.id));

    // Create some statuses
    for (let i = 0; i < 3; i++) {
      const status = await statusAPI.create({
        projectId: project.data!.id,
        name: `Status ${i}`,
      });
      createdStatusIds.push(status.data!.id);
    }

    // Read multiple times concurrently
    const readPromises = Array(10)
      .fill(null)
      .map(() => statusAPI.getByProject(project.data!.id));

    const results = await Promise.all(readPromises);

    // All reads should return same data
    const firstIds = results[0].data!.map((s) => s.id).sort();

    results.forEach((result) => {
      expect(result.status).toBe(200);
      const ids = result.data!.map((s) => s.id).sort();
      expect(ids).toEqual(firstIds);
    });
  });
});
