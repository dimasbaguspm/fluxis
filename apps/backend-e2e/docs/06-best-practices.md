# Best Practices for API Testing

## Overview

This document outlines recommended patterns, anti-patterns, and guidelines for writing effective, maintainable E2E API tests.

## General Testing Principles

### 1. Test Independence

**✅ DO**: Each test should be completely independent

```typescript
test("should create project", async ({ projectAPI }) => {
  // Create fresh data for this test
  const response = await projectAPI.create({
    name: "Test Project",
    status: "active",
  });

  expect(response.status).toBe(200);

  // Cleanup at the end
  await projectAPI.delete(response.data!.id);
});

test("should update project", async ({ projectAPI }) => {
  // Create its own data
  const created = await projectAPI.create({
    name: "Another Project",
    status: "active",
  });

  const updated = await projectAPI.update(created.data!.id, {
    name: "Updated Name",
  });

  expect(updated.status).toBe(200);

  await projectAPI.delete(created.data!.id);
});
```

**❌ DON'T**: Share state between tests

```typescript
// BAD - Tests depend on each other
let sharedProjectId: string;

test("create project", async ({ projectAPI }) => {
  const response = await projectAPI.create({ name: "Shared" });
  sharedProjectId = response.data!.id; // ❌ Shared state
});

test("update project", async ({ projectAPI }) => {
  // ❌ Depends on previous test
  await projectAPI.update(sharedProjectId, { name: "Updated" });
});
```

### 2. Arrange-Act-Assert Pattern

Structure tests with clear sections:

```typescript
test("should handle project status transition", async ({ projectAPI }) => {
  // ARRANGE: Setup test data
  const created = await projectAPI.create({
    name: "Test Project",
    status: "active",
  });
  const projectId = created.data!.id;

  // ACT: Perform the action being tested
  const response = await projectAPI.update(projectId, {
    status: "archived",
  });

  // ASSERT: Verify the outcome
  expect(response.status).toBe(200);
  expect(response.data?.status).toBe("archived");

  // Verify persistence (optional but thorough)
  const fetched = await projectAPI.getById(projectId);
  expect(fetched.data?.status).toBe("archived");

  // CLEANUP
  await projectAPI.delete(projectId);
});
```

### 3. Descriptive Test Names

**✅ DO**: Use clear, behavior-driven test names

```typescript
test("should return 404 when getting non-existent project", async () => {});
test("should reject project creation with missing required fields", async () => {});
test("should cascade delete all tasks when project is deleted", async () => {});
test("should maintain status order after reordering", async () => {});
```

**❌ DON'T**: Use vague or technical names

```typescript
test("test 1", async () => {}); // ❌ What does this test?
test("projects", async () => {}); // ❌ Too vague
test("DELETE /projects/{id}", async () => {}); // ❌ Describes endpoint, not behavior
```

### 4. Test One Thing

Each test should verify **one specific behavior**:

**✅ DO**: Single responsibility

```typescript
test("should create project with valid data", async ({ projectAPI }) => {
  const response = await projectAPI.create({
    name: "Test",
    status: "active",
  });

  expect(response.status).toBe(200);
  expect(response.data?.name).toBe("Test");
  expect(response.data?.id).toBeDefined();
});

test("should reject project with invalid status", async ({ projectAPI }) => {
  const response = await projectAPI.create({
    name: "Test",
    status: "invalid_status", // Invalid value
  });

  expect(response.status).toBe(400);
  expect(response.error).toBeDefined();
});
```

**❌ DON'T**: Test multiple unrelated behaviors

```typescript
test("project CRUD operations", async ({ projectAPI }) => {
  // ❌ Testing create, read, update, delete in one test
  const created = await projectAPI.create({ name: "Test" });
  const fetched = await projectAPI.getById(created.data!.id);
  const updated = await projectAPI.update(created.data!.id, {
    name: "Updated",
  });
  const deleted = await projectAPI.delete(created.data!.id);
  // Too many concerns in one test
});
```

## API Testing Patterns

### 1. Response Validation

Always validate both structure and content:

```typescript
test("should return complete project data", async ({ projectAPI }) => {
  const response = await projectAPI.create({
    name: "Test Project",
    description: "Test Description",
    status: "active",
  });

  // Status code
  expect(response.status).toBe(200);

  // Data exists
  expect(response.data).toBeDefined();

  // Required fields
  expect(response.data?.id).toBeDefined();
  expect(response.data?.name).toBe("Test Project");
  expect(response.data?.status).toBe("active");

  // Optional fields
  expect(response.data?.description).toBe("Test Description");

  // Auto-generated fields
  expect(response.data?.createdAt).toBeDefined();
  expect(response.data?.updatedAt).toBeDefined();

  // Field types
  expect(typeof response.data?.id).toBe("string");
  expect(typeof response.data?.createdAt).toBe("string");
});
```

### 2. Error Handling

Test error scenarios thoroughly:

```typescript
test.describe("Error Handling", () => {
  test("should return 400 for missing required field", async ({
    projectAPI,
  }) => {
    const response = await projectAPI.create({
      // Missing "name" field
      status: "active",
    } as any);

    expect(response.status).toBe(400);
    expect(response.error?.detail).toContain("name");
  });

  test("should return 404 for non-existent resource", async ({
    projectAPI,
  }) => {
    const fakeId = "00000000-0000-0000-0000-000000000000";
    const response = await projectAPI.getById(fakeId);

    expect(response.status).toBe(404);
  });

  test("should return 409 for duplicate slug", async ({ projectAPI }) => {
    const name = "Duplicate Name";

    // Create first project
    const first = await projectAPI.create({ name, status: "active" });
    expect(first.status).toBe(200);

    // Try to create duplicate
    const second = await projectAPI.create({ name, status: "active" });
    expect(second.status).toBe(409); // Conflict
    expect(second.error?.detail).toContain("already exists");

    // Cleanup
    await projectAPI.delete(first.data!.id);
  });
});
```

### 3. Pagination Testing

Verify pagination works correctly:

```typescript
test("should paginate projects correctly", async ({ projectAPI }) => {
  // Arrange: Create multiple projects
  const projectIds: string[] = [];
  for (let i = 1; i <= 15; i++) {
    const response = await projectAPI.create({
      name: `Project ${i}`,
      status: "active",
    });
    projectIds.push(response.data!.id);
  }

  // Act & Assert: Test first page
  const page1 = await projectAPI.getPaginated({ page: 1, limit: 10 });
  expect(page1.status).toBe(200);
  expect(page1.data?.items).toHaveLength(10);
  expect(page1.data?.page).toBe(1);
  expect(page1.data?.limit).toBe(10);
  expect(page1.data?.totalCount).toBeGreaterThanOrEqual(15);
  expect(page1.data?.totalPages).toBeGreaterThanOrEqual(2);

  // Test second page
  const page2 = await projectAPI.getPaginated({ page: 2, limit: 10 });
  expect(page2.status).toBe(200);
  expect(page2.data?.items.length).toBeGreaterThanOrEqual(5);

  // Verify no overlap between pages
  const page1Ids = page1.data!.items.map((p) => p.id);
  const page2Ids = page2.data!.items.map((p) => p.id);
  const overlap = page1Ids.filter((id) => page2Ids.includes(id));
  expect(overlap).toHaveLength(0);

  // Cleanup
  for (const id of projectIds) {
    await projectAPI.delete(id);
  }
});
```

### 4. Data Relationships

Test cascading operations and foreign keys:

```typescript
test("should cascade delete related entities", async ({
  projectAPI,
  statusAPI,
  taskAPI,
}) => {
  // Create project
  const project = await projectAPI.create({
    name: "Parent Project",
    status: "active",
  });
  const projectId = project.data!.id;

  // Create status
  const status = await statusAPI.create({
    projectId,
    name: "Todo",
  });
  const statusId = status.data!.id;

  // Create tasks
  const task1 = await taskAPI.create({
    projectId,
    statusId,
    title: "Task 1",
    priority: "high",
  });

  const task2 = await taskAPI.create({
    projectId,
    statusId,
    title: "Task 2",
    priority: "medium",
  });

  // Verify all created
  const statuses = await statusAPI.getByProject(projectId);
  expect(statuses.data?.items).toHaveLength(1);

  const tasks = await taskAPI.getPaginated({ projectId });
  expect(tasks.data?.items).toHaveLength(2);

  // Delete project
  await projectAPI.delete(projectId);

  // Verify cascade delete
  const statusesAfter = await statusAPI.getByProject(projectId);
  expect(statusesAfter.data?.items).toHaveLength(0);

  const tasksAfter = await taskAPI.getPaginated({ projectId });
  expect(tasksAfter.data?.items).toHaveLength(0);
});
```

### 5. Idempotency Testing

Test that operations are idempotent when expected:

```typescript
test("DELETE should be idempotent", async ({ projectAPI }) => {
  const created = await projectAPI.create({
    name: "Test",
    status: "active",
  });
  const projectId = created.data!.id;

  // First delete
  const firstDelete = await projectAPI.delete(projectId);
  expect(firstDelete.status).toBe(204);

  // Second delete (should return 404, but not crash)
  const secondDelete = await projectAPI.delete(projectId);
  expect(secondDelete.status).toBe(404);
});
```

## Test Data Management

### 1. Use Unique Test Data

Avoid conflicts with concurrent tests:

```typescript
test("should create unique project", async ({ projectAPI }) => {
  // Generate unique name
  const uniqueName = `Project ${Date.now()}-${Math.random()}`;

  const response = await projectAPI.create({
    name: uniqueName,
    status: "active",
  });

  expect(response.status).toBe(200);
  expect(response.data?.name).toBe(uniqueName);

  await projectAPI.delete(response.data!.id);
});
```

### 2. Cleanup Test Data

Always clean up created resources:

**Option 1: Inline cleanup**

```typescript
test("should update project", async ({ projectAPI }) => {
  const created = await projectAPI.create({ name: "Test", status: "active" });
  const projectId = created.data!.id;

  try {
    const updated = await projectAPI.update(projectId, { name: "Updated" });
    expect(updated.status).toBe(200);
  } finally {
    // Always cleanup, even if test fails
    await projectAPI.delete(projectId);
  }
});
```

**Option 2: afterEach hook**

```typescript
test.describe("Projects", () => {
  const projectIds: string[] = [];

  test.afterEach(async ({ projectAPI }) => {
    // Cleanup all created projects
    for (const id of projectIds) {
      await projectAPI.delete(id).catch(() => {
        // Ignore errors (already deleted or not found)
      });
    }
    projectIds.length = 0;
  });

  test("test 1", async ({ projectAPI }) => {
    const response = await projectAPI.create({
      name: "Test",
      status: "active",
    });
    projectIds.push(response.data!.id);
    // Test logic
  });

  test("test 2", async ({ projectAPI }) => {
    const response = await projectAPI.create({
      name: "Test",
      status: "active",
    });
    projectIds.push(response.data!.id);
    // Test logic
  });
});
```

### 3. Test Data Factories

Create reusable data factories:

```typescript
// test-helpers.ts
export async function createTestProject(
  projectAPI: ProjectAPIClient,
  overrides?: Partial<ProjectCreateRequest>
): Promise<ProjectResponse> {
  const response = await projectAPI.create({
    name: `Test Project ${Date.now()}`,
    status: "active",
    ...overrides,
  });

  if (!response.data) {
    throw new Error("Failed to create test project");
  }

  return response.data;
}

export async function createTestTask(
  taskAPI: TaskAPIClient,
  projectId: string,
  statusId: string,
  overrides?: Partial<TaskCreateRequest>
): Promise<TaskResponse> {
  const response = await taskAPI.create({
    projectId,
    statusId,
    title: `Test Task ${Date.now()}`,
    priority: "medium",
    ...overrides,
  });

  if (!response.data) {
    throw new Error("Failed to create test task");
  }

  return response.data;
}

// Usage in tests
test("should list tasks by project", async ({
  projectAPI,
  statusAPI,
  taskAPI,
}) => {
  const project = await createTestProject(projectAPI);
  const status = await createTestStatus(statusAPI, project.id);
  const task1 = await createTestTask(taskAPI, project.id, status.id);
  const task2 = await createTestTask(taskAPI, project.id, status.id);

  const response = await taskAPI.getPaginated({ projectId: project.id });

  expect(response.data?.items).toHaveLength(2);

  // Cleanup
  await projectAPI.delete(project.id);
});
```

## Assertion Best Practices

### 1. Specific Assertions

Be specific about what you're testing:

**✅ DO**: Specific checks

```typescript
expect(response.status).toBe(200);
expect(response.data?.name).toBe("Expected Name");
expect(response.data?.items).toHaveLength(5);
expect(response.data?.createdAt).toMatch(/^\d{4}-\d{2}-\d{2}/);
```

**❌ DON'T**: Vague checks

```typescript
expect(response.status).toBeTruthy(); // ❌ Any truthy value passes
expect(response.data).toBeDefined(); // ❌ Doesn't check contents
expect(response.data?.items.length > 0); // ❌ Doesn't verify exact count
```

### 2. Multiple Assertions

It's okay to have multiple assertions in one test:

```typescript
test("should return complete user profile", async ({ userAPI }) => {
  const response = await userAPI.getProfile();

  // All these assertions relate to the same behavior
  expect(response.status).toBe(200);
  expect(response.data?.id).toBeDefined();
  expect(response.data?.username).toBeDefined();
  expect(response.data?.email).toMatch(/^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/);
  expect(response.data?.createdAt).toBeDefined();
  expect(response.data?.password).toBeUndefined(); // Should not expose password
});
```

### 3. Custom Matchers

Create custom matchers for common patterns:

```typescript
// test-matchers.ts
expect.extend({
  toBeValidUUID(received: string) {
    const uuidRegex =
      /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    const pass = uuidRegex.test(received);

    return {
      pass,
      message: () =>
        pass
          ? `Expected ${received} not to be a valid UUID`
          : `Expected ${received} to be a valid UUID`,
    };
  },

  toBeISO8601(received: string) {
    const date = new Date(received);
    const pass = !isNaN(date.getTime()) && received === date.toISOString();

    return {
      pass,
      message: () =>
        pass
          ? `Expected ${received} not to be a valid ISO8601 date`
          : `Expected ${received} to be a valid ISO8601 date`,
    };
  },
});

// Usage
test("should return valid timestamps", async ({ projectAPI }) => {
  const response = await projectAPI.create({ name: "Test", status: "active" });

  expect(response.data?.id).toBeValidUUID();
  expect(response.data?.createdAt).toBeISO8601();
  expect(response.data?.updatedAt).toBeISO8601();
});
```

## Performance Best Practices

### 1. Minimize API Calls

Reduce unnecessary requests:

**✅ DO**: Efficient testing

```typescript
test("should verify project structure", async ({ projectAPI }) => {
  const response = await projectAPI.create({
    name: "Test",
    description: "Description",
    status: "active",
  });

  // Single request validates everything
  expect(response.status).toBe(200);
  expect(response.data).toMatchObject({
    name: "Test",
    description: "Description",
    status: "active",
  });
});
```

**❌ DON'T**: Excessive requests

```typescript
test("should handle project", async ({ projectAPI }) => {
  const created = await projectAPI.create({ name: "Test", status: "active" });
  const fetched1 = await projectAPI.getById(created.data!.id); // ❌ Unnecessary
  const fetched2 = await projectAPI.getById(created.data!.id); // ❌ Duplicate
  const listed = await projectAPI.getPaginated(); // ❌ Not needed
  // Too many calls for simple verification
});
```

### 2. Parallel Execution

Run independent tests in parallel:

```typescript
// playwright.config.ts
export default defineConfig({
  fullyParallel: true, // Enable parallel execution
  workers: 4, // Number of parallel workers
});
```

**Note**: Only enable for truly independent tests. Disable for tests that share database state.

### 3. Batch Operations

Test batch operations when available:

```typescript
test("should create multiple projects efficiently", async ({ projectAPI }) => {
  // If API supports batch creation
  const response = await projectAPI.createBatch([
    { name: "Project 1", status: "active" },
    { name: "Project 2", status: "active" },
    { name: "Project 3", status: "active" },
  ]);

  expect(response.status).toBe(200);
  expect(response.data?.created).toHaveLength(3);
});
```

## Debugging Tips

### 1. Verbose Logging

Add logging for debugging:

```typescript
test("should handle complex workflow", async ({ projectAPI, taskAPI }) => {
  console.log("Creating project...");
  const project = await projectAPI.create({ name: "Test", status: "active" });
  console.log("Created project:", project.data?.id);

  console.log("Creating task...");
  const task = await taskAPI.create({
    projectId: project.data!.id,
    statusId: "some-status-id",
    title: "Test Task",
    priority: "high",
  });
  console.log("Created task:", task.data?.id);

  // ... more logic
});
```

### 2. Use test.only for Focused Testing

```typescript
// Run only this test
test.only("should fix this specific issue", async ({ projectAPI }) => {
  // Debug this test in isolation
});
```

### 3. Inspect Full Response

```typescript
test("debugging response structure", async ({ projectAPI }) => {
  const response = await projectAPI.create({ name: "Test", status: "active" });

  // Log entire response for inspection
  console.log("Full response:", JSON.stringify(response, null, 2));

  // Continue with assertions
});
```

## Common Anti-Patterns

### ❌ Hardcoded IDs

```typescript
// BAD - Assumes ID exists in database
test("should get project", async ({ projectAPI }) => {
  const response = await projectAPI.getById("hardcoded-id-12345");
  expect(response.status).toBe(200);
});
```

**Solution**: Create test data dynamically

### ❌ Sleep/Wait Instead of Polling

```typescript
// BAD - Arbitrary wait
test("should process async operation", async ({ projectAPI }) => {
  await projectAPI.triggerAsyncOperation();
  await new Promise((resolve) => setTimeout(resolve, 5000)); // ❌ Arbitrary 5s wait
  const result = await projectAPI.getResult();
});
```

**Solution**: Poll until ready

```typescript
// GOOD - Poll until ready
async function waitForResult(projectAPI: ProjectAPIClient, maxAttempts = 10) {
  for (let i = 0; i < maxAttempts; i++) {
    const response = await projectAPI.getResult();
    if (response.status === 200 && response.data?.status === "completed") {
      return response;
    }
    await new Promise((resolve) => setTimeout(resolve, 500));
  }
  throw new Error("Operation did not complete in time");
}
```

### ❌ Testing Implementation Details

```typescript
// BAD - Tests internal structure instead of behavior
test("should use SHA-256 for password hashing", async () => {
  // ❌ Don't test internal implementation
});
```

**Solution**: Test behavior and contracts

```typescript
// GOOD - Test observable behavior
test("should not expose password in user response", async ({ userAPI }) => {
  const response = await userAPI.getProfile();
  expect(response.data?.password).toBeUndefined();
});
```

## Summary

### Key Principles

- ✅ **Independence** - Tests don't depend on each other
- ✅ **Clarity** - Descriptive names, clear structure
- ✅ **Focus** - One behavior per test
- ✅ **Completeness** - Test success, failure, and edge cases
- ✅ **Efficiency** - Minimize unnecessary API calls
- ✅ **Cleanup** - Always remove test data

### Test Structure

```typescript
test("should {behavior} when {condition}", async ({ fixtures }) => {
  // ARRANGE: Setup test data
  // ACT: Perform action
  // ASSERT: Verify outcome
  // CLEANUP: Remove test data
});
```

### Coverage Checklist

For comprehensive API testing:

- [ ] Happy path (valid inputs)
- [ ] Validation errors (invalid inputs)
- [ ] Not found errors (missing resources)
- [ ] Conflict errors (duplicates, constraints)
- [ ] Authorization (access control)
- [ ] Edge cases (boundaries, nulls)
- [ ] Relationships (foreign keys, cascades)
- [ ] Pagination (multiple pages, limits)
- [ ] Idempotency (repeated operations)
- [ ] Performance (large datasets)

Following these practices ensures your E2E tests are reliable, maintainable, and provide comprehensive coverage of API behavior.
