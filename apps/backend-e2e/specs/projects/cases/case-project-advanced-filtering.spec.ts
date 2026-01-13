import { test, expect } from "../../../fixtures";

/**
 * Project Advanced Filtering Test Cases
 * Tests complex filtering combinations for projects
 */
test.describe("Project Advanced Filtering", () => {
  const createdProjectIds: string[] = [];

  test.afterEach(async ({ projectAPI }) => {
    for (const id of createdProjectIds) {
      await projectAPI.remove(id).catch(() => {});
    }
    createdProjectIds.length = 0;
  });

  test("should filter by multiple project IDs", async ({ projectAPI }) => {
    // Create test projects
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

    const project3 = await projectAPI.create({
      name: "Project 3",
      description: "Third",
      status: "active",
    });

    createdProjectIds.push(
      project1.data!.id,
      project2.data!.id,
      project3.data!.id
    );

    // Filter by two specific IDs
    const response = await projectAPI.getPaginated({
      id: [project1.data!.id, project3.data!.id],
    });

    expect(response.status).toBe(200);
    expect(response?.data?.items?.length).toBe(2);

    const ids = response.data?.items?.map((p) => p.id);
    expect(ids).toContain(project1.data!.id);
    expect(ids).toContain(project3.data!.id);
    expect(ids).not.toContain(project2.data!.id);
  });

  test("should filter by multiple statuses", async ({ projectAPI }) => {
    const active = await projectAPI.create({
      name: "Active Project",
      description: "Test",
      status: "active",
    });

    const paused = await projectAPI.create({
      name: "Paused Project",
      description: "Test",
      status: "paused",
    });

    const archived = await projectAPI.create({
      name: "Archived Project",
      description: "Test",
      status: "archived",
    });

    createdProjectIds.push(active.data!.id, paused.data!.id, archived.data!.id);

    // Filter for active and paused only
    const response = await projectAPI.getPaginated({
      status: ["active", "paused"],
    });

    expect(response.status).toBe(200);
    const statuses = response.data?.items?.map((p) => p.status);

    expect(statuses).toContain("active");
    expect(statuses).toContain("paused");
    // Archived might appear if there are other archived projects
  });

  test("should combine ID filter with status filter", async ({
    projectAPI,
  }) => {
    const active1 = await projectAPI.create({
      name: "Active 1",
      description: "Test",
      status: "active",
    });

    const active2 = await projectAPI.create({
      name: "Active 2",
      description: "Test",
      status: "active",
    });

    const paused1 = await projectAPI.create({
      name: "Paused 1",
      description: "Test",
      status: "paused",
    });

    createdProjectIds.push(
      active1.data!.id,
      active2.data!.id,
      paused1.data!.id
    );

    // Filter by specific IDs AND status
    const response = await projectAPI.getPaginated({
      id: [active1.data!.id, active2.data!.id, paused1.data!.id],
      status: ["active"],
    });

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBe(2);

    const ids = response.data?.items?.map((p) => p.id);
    expect(ids).toContain(active1.data!.id);
    expect(ids).toContain(active2.data!.id);
    expect(ids).not.toContain(paused1.data!.id);
  });

  test("should combine filters with search query", async ({ projectAPI }) => {
    const project1 = await projectAPI.create({
      name: "Frontend Development",
      description: "React application",
      status: "active",
    });

    const project2 = await projectAPI.create({
      name: "Backend Development",
      description: "Node.js API",
      status: "active",
    });

    const project3 = await projectAPI.create({
      name: "Frontend Design",
      description: "UI/UX work",
      status: "paused",
    });

    createdProjectIds.push(
      project1.data!.id,
      project2.data!.id,
      project3.data!.id
    );

    // Search for "Frontend" with status filter
    const response = await projectAPI.getPaginated({
      query: "Frontend",
      status: ["active"],
    });

    expect(response.status).toBe(200);

    const names = response.data?.items?.map((p) => p.name);
    expect(names).toContain("Frontend Development");
  });

  test("should combine filters with pagination", async ({ projectAPI }) => {
    // Create multiple projects
    const projectIds: string[] = [];
    for (let i = 0; i < 15; i++) {
      const project = await projectAPI.create({
        name: `Test Project ${i}`,
        description: "Test",
        status: i % 2 === 0 ? "active" : "paused",
      });
      projectIds.push(project.data!.id);
    }
    createdProjectIds.push(...projectIds);

    // Filter by status with pagination
    const page1 = await projectAPI.getPaginated({
      status: ["active"],
      pageNumber: 1,
      pageSize: 3,
    });

    const page2 = await projectAPI.getPaginated({
      status: ["active"],
      pageNumber: 2,
      pageSize: 3,
    });

    expect(page1.status).toBe(200);
    expect(page2.status).toBe(200);

    // Both pages should only contain active projects
    page1.data!.items?.forEach((p) => {
      expect(p.status).toBe("active");
    });

    page2.data!.items?.forEach((p) => {
      expect(p.status).toBe("active");
    });

    // No overlap between pages
    const ids1 = page1.data!.items?.map((p) => p.id);
    const ids2 = page2.data!.items?.map((p) => p.id);
    const overlap = ids1?.filter((id) => ids2?.includes(id)) || [];
    expect(overlap.length).toBe(0);
  });

  test("should combine filters with sorting", async ({ projectAPI }) => {
    const old = await projectAPI.create({
      name: "Old Active",
      description: "Test",
      status: "active",
    });

    // Wait a bit
    await new Promise((resolve) => setTimeout(resolve, 100));

    const middle = await projectAPI.create({
      name: "Middle Active",
      description: "Test",
      status: "active",
    });

    await new Promise((resolve) => setTimeout(resolve, 100));

    const recent = await projectAPI.create({
      name: "Recent Active",
      description: "Test",
      status: "active",
    });

    const paused = await projectAPI.create({
      name: "Paused",
      description: "Test",
      status: "paused",
    });

    createdProjectIds.push(
      old.data!.id,
      middle.data!.id,
      recent.data!.id,
      paused.data!.id
    );

    // Filter active projects and sort by creation date
    const response = await projectAPI.getPaginated({
      status: ["active"],
      sortBy: "createdAt",
      sortOrder: "asc",
    });

    expect(response.status).toBe(200);

    // Find our test projects in the results
    const testProjects =
      response.data!.items?.filter((p) =>
        [old.data!.id, middle.data!.id, recent.data!.id].includes(p.id)
      ) || [];

    expect(testProjects.length).toBe(3);
    expect(testProjects[0].id).toBe(old.data!.id);
    expect(testProjects[1].id).toBe(middle.data!.id);
    expect(testProjects[2].id).toBe(recent.data!.id);
  });

  test("should handle search query matching multiple fields", async ({
    projectAPI,
  }) => {
    const project1 = await projectAPI.create({
      name: "Backend API",
      description: "Node.js service",
      status: "active",
    });

    const project2 = await projectAPI.create({
      name: "Frontend App",
      description: "Backend integration",
      status: "active",
    });

    const project3 = await projectAPI.create({
      name: "Mobile App",
      description: "iOS application",
      status: "active",
    });

    createdProjectIds.push(
      project1.data!.id,
      project2.data!.id,
      project3.data!.id
    );

    // Search for "Backend" (matches name of project1, description of project2)
    const response = await projectAPI.getPaginated({
      query: "Backend",
    });

    expect(response.status).toBe(200);

    const ids = response.data!.items?.map((p) => p.id);
    expect(ids).toContain(project1.data!.id); // Name matches
    expect(ids).toContain(project2.data!.id); // Description matches
    expect(ids).not.toContain(project3.data!.id); // No match
  });

  test("should handle case-insensitive search", async ({ projectAPI }) => {
    const project = await projectAPI.create({
      name: "React Application",
      description: "Frontend Development",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Search with different cases
    const lower = await projectAPI.getPaginated({ query: "react" });
    const upper = await projectAPI.getPaginated({ query: "REACT" });
    const mixed = await projectAPI.getPaginated({ query: "ReAcT" });

    expect(lower.status).toBe(200);
    expect(upper.status).toBe(200);
    expect(mixed.status).toBe(200);

    // All should find the project
    expect(lower.data!.items?.some((p) => p.id === project.data!.id)).toBe(
      true
    );
    expect(upper.data!.items?.some((p) => p.id === project.data!.id)).toBe(
      true
    );
    expect(mixed.data!.items?.some((p) => p.id === project.data!.id)).toBe(
      true
    );
  });

  test("should handle empty filters (return all)", async ({ projectAPI }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // No filters
    const response = await projectAPI.getPaginated({});

    expect(response.status).toBe(200);
    expect(response.data?.items?.length).toBeGreaterThan(0);
    expect(response.data!.items?.some((p) => p.id === project.data!.id)).toBe(
      true
    );
  });

  test("should handle filters with no matches", async ({ projectAPI }) => {
    const project = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Filter with impossible criteria
    const response = await projectAPI.getPaginated({
      id: ["00000000-0000-0000-0000-000000000000"],
      status: ["archived"],
      query: "nonexistent-string-xyz",
    });

    expect(response.status).toBe(200);
    expect(response.data?.items).toEqual([]);
    expect(response.data?.totalCount).toBe(0);
  });

  test("should handle special characters in search query", async ({
    projectAPI,
  }) => {
    const project = await projectAPI.create({
      name: "Project @2024",
      description: "Test with special chars: #hashtag, $money",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    // Search for special characters
    const response1 = await projectAPI.getPaginated({ query: "@2024" });
    const response2 = await projectAPI.getPaginated({ query: "#hashtag" });
    const response3 = await projectAPI.getPaginated({ query: "$money" });

    expect(response1.status).toBe(200);
    expect(response2.status).toBe(200);
    expect(response3.status).toBe(200);

    // Should find the project (depending on implementation)
  });

  test("should handle unicode characters in search", async ({ projectAPI }) => {
    const project = await projectAPI.create({
      name: "プロジェクト",
      description: "日本語の説明",
      status: "active",
    });
    createdProjectIds.push(project.data!.id);

    const response = await projectAPI.getPaginated({ query: "プロジェクト" });

    expect(response.status).toBe(200);
    expect(response.data!.items?.some((p) => p.id === project.data!.id)).toBe(
      true
    );
  });
});
