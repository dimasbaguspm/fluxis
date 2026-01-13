# Test Organization

## Overview

Tests are organized into **two distinct categories** based on their purpose and scope:

1. **Common Tests** - Standard CRUD operations
2. **Case Tests** - Business logic, edge cases, and feature-specific scenarios

## File Structure

```
specs/
├── auth/
│   ├── auth.spec.ts              # Common: Login, refresh, logout
│   └── cases/
│       ├── case-token-expiry.spec.ts
│       └── case-concurrent-sessions.spec.ts
│
├── projects/
│   ├── projects.spec.ts          # Common: CRUD operations
│   └── cases/
│       ├── case-project-cascade-delete.spec.ts
│       ├── case-project-slug-uniqueness.spec.ts
│       └── case-archived-project-restrictions.spec.ts
│
├── statuses/
│   ├── statuses.spec.ts          # Common: CRUD operations
│   └── cases/
│       ├── case-status-reordering.spec.ts
│       └── case-status-task-migration.spec.ts
│
└── tasks/
    ├── tasks.spec.ts             # Common: CRUD operations
    └── cases/
        ├── case-task-priority-ordering.spec.ts
        └── case-task-status-transitions.spec.ts
```

## Common Tests (`{entity}.spec.ts`)

### Purpose

Test **standard CRUD operations** that apply to most entities:

- **C**reate - POST new entities
- **R**ead - GET single entity and paginated lists
- **U**pdate - PATCH existing entities
- **D**elete - DELETE entities

### Location

Directly in the entity folder: `specs/{entity}/{entity}.spec.ts`

### Structure

```typescript
import { test, expect } from "../../fixtures";

test.describe("{Entity} API", () => {
  test.describe("POST /{entities}", () => {
    test("should create entity with valid data", async ({ entityAPI }) => {
      // Test creation with valid input
    });

    test("should fail to create with invalid data", async ({ entityAPI }) => {
      // Test validation errors
    });

    test("should fail to create with missing required fields", async ({
      entityAPI,
    }) => {
      // Test missing field errors
    });
  });

  test.describe("GET /{entities}/{id}", () => {
    test("should get entity by id", async ({ entityAPI }) => {
      // Test retrieval by ID
    });

    test("should return 404 for non-existent entity", async ({ entityAPI }) => {
      // Test not found errors
    });
  });

  test.describe("GET /{entities}", () => {
    test("should get paginated list", async ({ entityAPI }) => {
      // Test pagination
    });

    test("should filter by query parameters", async ({ entityAPI }) => {
      // Test filtering
    });
  });

  test.describe("PATCH /{entities}/{id}", () => {
    test("should update entity", async ({ entityAPI }) => {
      // Test update
    });

    test("should fail to update non-existent entity", async ({ entityAPI }) => {
      // Test not found errors
    });
  });

  test.describe("DELETE /{entities}/{id}", () => {
    test("should delete entity", async ({ entityAPI }) => {
      // Test deletion
    });

    test("should return 404 when deleting non-existent entity", async ({
      entityAPI,
    }) => {
      // Test not found errors
    });
  });
});
```

### Example: Auth Common Tests

```typescript
// specs/auth/auth.spec.ts
import { test, expect } from "../../fixtures";

test.describe("Authentication API", () => {
  test.describe("POST /auth/login", () => {
    test("should successfully login with valid credentials", async ({
      authAPI,
    }) => {
      const username = "admin_username";
      const password = "admin_password";

      const response = await authAPI.login(username, password);

      expect(response.status).toBe(200);
      expect(response.data).toBeDefined();
      expect(response.data?.accessToken).toBeDefined();
      expect(response.data?.refreshToken).toBeDefined();
      expect(typeof response.data?.accessToken).toBe("string");
      expect(typeof response.data?.refreshToken).toBe("string");
    });

    test("should fail login with invalid credentials", async ({ authAPI }) => {
      const response = await authAPI.login("invalid_user", "invalid_password");

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.status).toBeLessThan(500);
      expect(response.error).toBeDefined();
    });

    test("should fail login with empty username", async ({ authAPI }) => {
      const response = await authAPI.login("", "somepassword");

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });

    test("should fail login with empty password", async ({ authAPI }) => {
      const username = "admin_username";
      const response = await authAPI.login(username, "");

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });

  test.describe("POST /auth/refresh", () => {
    test("should successfully refresh token with valid refresh token", async ({
      authAPI,
    }) => {
      // First login to get tokens
      const username = "admin_username";
      const password = "admin_password";
      const loginResponse = await authAPI.login(username, password);

      expect(loginResponse.data).toBeDefined();
      const refreshToken = loginResponse.data!.refreshToken;

      // Now refresh the token
      const refreshResponse = await authAPI.refresh(refreshToken);

      expect(refreshResponse.status).toBe(200);
      expect(refreshResponse.data).toBeDefined();
      expect(refreshResponse.data?.accessToken).toBeDefined();
      expect(typeof refreshResponse.data?.accessToken).toBe("string");
    });

    test("should fail refresh with invalid token", async ({ authAPI }) => {
      const response = await authAPI.refresh("invalid_refreshToken");

      expect(response.status).toBeGreaterThanOrEqual(400);
      expect(response.error).toBeDefined();
    });
  });
});
```

### Common Test Checklist

For each entity, ensure these tests exist:

- [ ] Create with valid data
- [ ] Create with invalid data (validation errors)
- [ ] Create with missing required fields
- [ ] Get by ID (success)
- [ ] Get by ID (404 not found)
- [ ] Get paginated list
- [ ] Get with query filters
- [ ] Update with valid data
- [ ] Update with invalid data
- [ ] Update non-existent entity (404)
- [ ] Delete entity (success)
- [ ] Delete non-existent entity (404)

## Case Tests (`case-{name}.spec.ts`)

### Purpose

Test **business requirements, edge cases, and complex scenarios**:

- **Business rules** - Domain-specific logic
- **Edge cases** - Boundary conditions, null values
- **Integration** - Multi-entity workflows
- **Security** - Authorization, access control
- **Performance** - Large datasets, rate limiting
- **Error handling** - Specific error scenarios

### Location

In the `cases/` subfolder: `specs/{entity}/cases/case-{name}.spec.ts`

### Naming Convention

Use descriptive names that explain the scenario:

- `case-token-expiry.spec.ts` - Test JWT expiration behavior
- `case-project-cascade-delete.spec.ts` - Test cascade deletion
- `case-status-reordering.spec.ts` - Test drag-and-drop reordering
- `case-task-priority-ordering.spec.ts` - Test priority-based sorting

### Structure

```typescript
import { test, expect } from "../../../fixtures";

test.describe("Case: {Scenario Description}", () => {
  test("should {behavior} when {condition}", async ({
    entityAPI,
    otherAPI,
  }) => {
    // Arrange: Setup test data
    // Act: Perform actions
    // Assert: Verify behavior
  });

  test("should {behavior} when {edge case}", async ({ entityAPI }) => {
    // Test edge case
  });
});
```

### Example: Project Cascade Delete

```typescript
// specs/projects/cases/case-project-cascade-delete.spec.ts
import { test, expect } from "../../../fixtures";

test.describe("Case: Project Cascade Delete", () => {
  test("should delete all statuses when project is deleted", async ({
    projectAPI,
    statusAPI,
  }) => {
    // Arrange: Create project with statuses
    const projectResponse = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });

    const projectId = projectResponse.data!.id;

    await statusAPI.create({
      projectId,
      name: "Todo",
    });

    await statusAPI.create({
      projectId,
      name: "In Progress",
    });

    // Verify statuses exist
    const statusesResponse = await statusAPI.getByProject(projectId);
    expect(statusesResponse.data?.items).toHaveLength(2);

    // Act: Delete project
    const deleteResponse = await projectAPI.delete(projectId);
    expect(deleteResponse.status).toBe(204);

    // Assert: Statuses should also be deleted
    const statusesAfterDelete = await statusAPI.getByProject(projectId);
    expect(statusesAfterDelete.data?.items).toHaveLength(0);
  });

  test("should delete all tasks when project is deleted", async ({
    projectAPI,
    statusAPI,
    taskAPI,
  }) => {
    // Arrange: Create project with status and tasks
    const projectResponse = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });

    const projectId = projectResponse.data!.id;

    const statusResponse = await statusAPI.create({
      projectId,
      name: "Todo",
    });

    const statusId = statusResponse.data!.id;

    await taskAPI.create({
      projectId,
      statusId,
      title: "Task 1",
      priority: "medium",
    });

    await taskAPI.create({
      projectId,
      statusId,
      title: "Task 2",
      priority: "high",
    });

    // Verify tasks exist
    const tasksResponse = await taskAPI.getPaginated({ projectId });
    expect(tasksResponse.data?.items).toHaveLength(2);

    // Act: Delete project
    await projectAPI.delete(projectId);

    // Assert: Tasks should also be deleted
    const tasksAfterDelete = await taskAPI.getPaginated({ projectId });
    expect(tasksAfterDelete.data?.items).toHaveLength(0);
  });
});
```

### Example: Status Reordering

```typescript
// specs/statuses/cases/case-status-reordering.spec.ts
import { test, expect } from "../../../fixtures";

test.describe("Case: Status Reordering", () => {
  test("should reorder statuses maintaining order field", async ({
    projectAPI,
    statusAPI,
  }) => {
    // Arrange: Create project with multiple statuses
    const projectResponse = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });

    const projectId = projectResponse.data!.id;

    const todo = await statusAPI.create({ projectId, name: "Todo" });
    const inProgress = await statusAPI.create({
      projectId,
      name: "In Progress",
    });
    const done = await statusAPI.create({ projectId, name: "Done" });

    // Initial order: [Todo, In Progress, Done]
    const initialStatuses = await statusAPI.getByProject(projectId);
    expect(initialStatuses.data?.items.map((s) => s.name)).toEqual([
      "Todo",
      "In Progress",
      "Done",
    ]);

    // Act: Reorder to [Done, Todo, In Progress]
    const reorderResponse = await statusAPI.reorder({
      projectId,
      ids: [done.data!.id, todo.data!.id, inProgress.data!.id],
    });

    expect(reorderResponse.status).toBe(200);

    // Assert: Verify new order
    const reorderedStatuses = await statusAPI.getByProject(projectId);
    expect(reorderedStatuses.data?.items.map((s) => s.name)).toEqual([
      "Done",
      "Todo",
      "In Progress",
    ]);
  });

  test("should fail to reorder with invalid status IDs", async ({
    projectAPI,
    statusAPI,
  }) => {
    const projectResponse = await projectAPI.create({
      name: "Test Project",
      description: "Test",
      status: "active",
    });

    const projectId = projectResponse.data!.id;

    // Act: Try to reorder with invalid IDs
    const response = await statusAPI.reorder({
      projectId,
      ids: ["invalid-uuid", "another-invalid-uuid"],
    });

    // Assert: Should fail
    expect(response.status).toBeGreaterThanOrEqual(400);
    expect(response.error).toBeDefined();
  });
});
```

### Case Test Guidelines

#### When to Create a Case Test

Create a case test when:

- ✅ Testing **business rules** specific to the domain
- ✅ Testing **multi-entity interactions** (e.g., cascade delete)
- ✅ Testing **complex workflows** (e.g., state transitions)
- ✅ Testing **edge cases** not covered by common tests
- ✅ Testing **security/authorization** beyond basic 401/403
- ✅ Testing **performance requirements** (e.g., pagination limits)
- ✅ Reproducing and preventing **bug regressions**

#### What NOT to Include in Case Tests

Avoid duplicating common tests:

- ❌ Basic CRUD operations (already in common tests)
- ❌ Simple validation errors (covered in common tests)
- ❌ Standard 404 errors (covered in common tests)

## Test Organization Best Practices

### 1. Descriptive Test Names

Use clear, behavior-driven test names:

```typescript
// ✅ GOOD
test("should delete all related statuses when project is deleted", async () => {});

// ❌ BAD
test("test project delete", async () => {});
```

### 2. Arrange-Act-Assert Pattern

```typescript
test("should ...", async ({ entityAPI }) => {
  // Arrange: Setup test data
  const created = await entityAPI.create({
    /* data */
  });

  // Act: Perform the action
  const response = await entityAPI.update(created.data!.id, {
    /* changes */
  });

  // Assert: Verify the outcome
  expect(response.status).toBe(200);
  expect(response.data?.someField).toBe("expected value");
});
```

### 3. Test Independence

Each test should be **completely independent**:

```typescript
// ✅ GOOD - Each test creates its own data
test("test 1", async ({ projectAPI }) => {
  const project = await projectAPI.create({ name: "Project 1" });
  // Test logic
});

test("test 2", async ({ projectAPI }) => {
  const project = await projectAPI.create({ name: "Project 2" });
  // Test logic
});

// ❌ BAD - Tests share state
let sharedProject;

test("test 1", async ({ projectAPI }) => {
  sharedProject = await projectAPI.create({ name: "Shared" });
});

test("test 2", async () => {
  // Depends on test 1 running first
  await projectAPI.update(sharedProject.id, {
    /* ... */
  });
});
```

### 4. Cleanup After Tests

Use `afterEach` or `afterAll` for cleanup:

```typescript
test.describe("Projects", () => {
  const createdProjectIds: string[] = [];

  test.afterEach(async ({ projectAPI }) => {
    // Cleanup projects created in tests
    for (const id of createdProjectIds) {
      await projectAPI.delete(id).catch(() => {
        // Ignore errors (project may already be deleted)
      });
    }
    createdProjectIds.length = 0; // Clear array
  });

  test("should create project", async ({ projectAPI }) => {
    const response = await projectAPI.create({ name: "Test" });
    createdProjectIds.push(response.data!.id);
    // Test logic
  });
});
```

### 5. Group Related Tests

Use nested `test.describe` blocks:

```typescript
test.describe("Projects API", () => {
  test.describe("POST /projects", () => {
    test.describe("Success Cases", () => {
      test("should create with minimum fields", async () => {});
      test("should create with all fields", async () => {});
    });

    test.describe("Validation Errors", () => {
      test("should fail with missing name", async () => {});
      test("should fail with invalid status", async () => {});
    });
  });
});
```

## Directory Structure Example

Complete example for the `projects` entity:

```
specs/projects/
├── projects.spec.ts                          # Common CRUD tests
└── cases/
    ├── case-project-cascade-delete.spec.ts   # Cascade deletion behavior
    ├── case-project-slug-uniqueness.spec.ts  # Slug generation and uniqueness
    ├── case-archived-projects.spec.ts        # Archived project restrictions
    ├── case-project-status-count.spec.ts     # Status count aggregation
    └── case-project-permissions.spec.ts      # Role-based access control
```

## Summary

### Common Tests (`{entity}.spec.ts`)

- **Location**: `specs/{entity}/{entity}.spec.ts`
- **Purpose**: Standard CRUD operations
- **Scope**: Create, Read, Update, Delete, List
- **Coverage**: Basic validation, 404 errors, success cases

### Case Tests (`case-{name}.spec.ts`)

- **Location**: `specs/{entity}/cases/case-{name}.spec.ts`
- **Purpose**: Business logic, edge cases, integrations
- **Scope**: Domain-specific rules, multi-entity workflows
- **Coverage**: Complex scenarios, security, performance

This organization provides:

- **Clear separation** - Common vs. specialized tests
- **Easy navigation** - Find tests by entity or scenario
- **Maintainability** - Update common tests, add cases independently
- **Comprehensive coverage** - Both CRUD and business logic
