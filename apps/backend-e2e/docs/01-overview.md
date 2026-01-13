# Backend E2E Testing - Overview

## Introduction

The **backend-e2e** project provides comprehensive End-to-End (E2E) API testing for the Fluxis backend API. This test suite ensures API contracts, business logic, and data integrity across all endpoints.

## Tech Stack

### Core Testing Framework

- **Test Framework**: [Playwright Test](https://playwright.dev/docs/api/class-test) v1.48.0
- **Language**: TypeScript 5.6
- **Runtime**: Node.js with ESM modules
- **Type Generation**: [openapi-typescript](https://github.com/drwpow/openapi-typescript) v7.4.0
- **Environment**: [dotenv](https://github.com/motdotla/dotenv) v16.4.5

### Why Playwright for API Testing?

While Playwright is primarily known for browser automation, it provides excellent API testing capabilities:

- **Built-in API request context** - No need for additional HTTP libraries
- **Type-safe fixtures** - Extend test context with custom fixtures
- **Parallel execution** - Run tests concurrently with worker pools
- **Powerful assertions** - Rich assertion library with helpful error messages
- **Global setup/teardown** - Authenticate once, reuse tokens across tests
- **HTML/JSON reporting** - Visual test reports out of the box

## Project Structure

```
apps/backend-e2e/
├── fixtures/              # Custom test fixtures (API clients)
│   ├── base-client.ts     # Base API client with HTTP methods
│   ├── auth-client.ts     # Authentication API client
│   └── index.ts           # Fixture registration & test extension
│
├── specs/                 # Test specifications
│   ├── auth/              # Authentication tests
│   │   ├── auth.spec.ts   # Common CRUD tests
│   │   └── cases/         # Business logic & edge cases
│   ├── projects/          # Project entity tests
│   ├── statuses/          # Status entity tests
│   └── tasks/             # Task entity tests
│
├── types/                 # TypeScript type definitions
│   ├── common.ts          # Shared types (APIResponse, TestContext)
│   └── openapi.ts         # Generated types from OpenAPI spec
│
├── scripts/               # Utility scripts
│   └── generate-openapi-types.ts  # Type generator from openapi.yaml
│
├── global-setup.ts        # Playwright global setup (authentication)
├── playwright.config.ts   # Playwright configuration
├── package.json           # Dependencies & scripts
└── tsconfig.json          # TypeScript configuration
```

## Source of Truth: OpenAPI Specification

The **single source of truth** for API contracts is:

```
/api/openapi.yaml
```

This OpenAPI specification defines:

- **Endpoints** - Paths, HTTP methods, parameters
- **Request/Response schemas** - Data structures
- **Authentication** - Security schemes (JWT bearer tokens)
- **Validation rules** - Required fields, data types, constraints

### Type Generation Workflow

```
┌───────────────────────┐
│  api/openapi.yaml     │  ← Backend generates this
│  (Source of Truth)    │
└──────────┬────────────┘
           │
           │ npm run generate:types
           │
           ▼
┌───────────────────────┐
│  types/openapi.ts     │  ← TypeScript types generated
│  (components, ops)    │
└──────────┬────────────┘
           │
           │ import types
           │
           ▼
┌───────────────────────┐
│  fixtures/*.ts        │  ← API clients use types
│  (Type-safe clients)  │
└──────────┬────────────┘
           │
           │ test.use()
           │
           ▼
┌───────────────────────┐
│  specs/**/*.spec.ts   │  ← Tests are fully type-safe
│  (E2E test suites)    │
└───────────────────────┘
```

### Generating Types

```bash
npm run generate:types
```

This script:

1. Reads `/api/openapi.yaml`
2. Uses `openapi-typescript` to generate TypeScript types
3. Outputs to `/types/openapi.ts`

**Generated types include**:

- `components["schemas"]["{ModelName}"]` - Request/response models
- `operations["{operationId}"]` - Endpoint operation types
- `paths["{path}"]["{method}"]` - Path-specific types

## Test Organization

### Two Types of Tests

Tests are organized into **two categories**:

#### 1. Common Tests (`{entity}.spec.ts`)

Located directly in the entity folder: `specs/{entity}/{entity}.spec.ts`

**Purpose**: Test standard CRUD operations

- **Create** - POST new entities
- **Read** - GET single & paginated entities
- **Update** - PATCH existing entities
- **Delete** - DELETE entities
- **List** - GET paginated lists

**Example**: `specs/auth/auth.spec.ts`

- Login with valid credentials
- Login with invalid credentials
- Refresh token with valid token
- Refresh token with invalid token

#### 2. Case Tests (`case-{name}.spec.ts`)

Located in the `cases/` subfolder: `specs/{entity}/cases/case-{name}.spec.ts`

**Purpose**: Test business requirements, edge cases, and feature-specific scenarios

- **Business rules** - Domain-specific validation
- **Edge cases** - Boundary conditions, null handling
- **Integration scenarios** - Multi-entity workflows
- **Security** - Authorization, access control
- **Error handling** - Specific error conditions

**Example**: `specs/projects/cases/case-project-status-cascade.spec.ts`

- When deleting a project, all statuses should be deleted
- When deleting a status, all tasks should be reassigned to default status

### Naming Conventions

| Type          | Location                                   | Naming Pattern               |
| ------------- | ------------------------------------------ | ---------------------------- |
| Common Tests  | `specs/{entity}/{entity}.spec.ts`          | `{entity}.spec.ts`           |
| Case Tests    | `specs/{entity}/cases/case-{name}.spec.ts` | `case-{description}.spec.ts` |
| API Fixtures  | `fixtures/{entity}-client.ts`              | `{entity}-client.ts`         |
| Base Fixtures | `fixtures/base-client.ts`                  | `base-client.ts`             |
| Custom Types  | `types/common.ts`                          | Descriptive names            |

## Key Concepts

### 1. Type Safety

All API requests and responses are **fully typed** using generated OpenAPI types:

```typescript
import type { components } from "../types/openapi";

// Type-safe request model
type CreateProjectRequest = components["schemas"]["ProjectCreateModel"];

// Type-safe response model
type ProjectResponse = components["schemas"]["ProjectModel"];

// Usage in client
async createProject(data: CreateProjectRequest): Promise<APIResponse<ProjectResponse>> {
  return this.post<ProjectResponse>("/projects", data);
}
```

### 2. Custom Fixtures

Playwright fixtures extend the test context with reusable API clients:

```typescript
import { test, expect } from "../../fixtures";

test("should create project", async ({ projectAPI }) => {
  //                                   ^^^^^^^^^^^ Custom fixture
  const response = await projectAPI.create({ name: "Test Project" });
  expect(response.status).toBe(200);
});
```

### 3. Authentication Flow

Tests use **global setup** to authenticate once:

1. `global-setup.ts` logs in and saves tokens to `.auth/user.json`
2. All tests load tokens from storage state
3. Individual tests can re-authenticate if needed

### 4. Test Isolation

- Each test should be **independent**
- Use `beforeEach` for test-specific setup
- Use `afterEach` for cleanup
- Avoid shared state between tests

## Configuration

### Playwright Configuration (`playwright.config.ts`)

Key settings:

```typescript
export default defineConfig({
  testDir: "./specs", // Test location
  timeout: 30 * 1000, // 30s per test
  fullyParallel: false, // Sequential execution
  workers: 1, // Single worker (avoid conflicts)
  use: {
    baseURL: "http://localhost:8081",
    storageState: ".auth/user.json", // Authenticated state
  },
  globalSetup: "./global-setup.ts", // Pre-test authentication
});
```

### Environment Variables

Configure via `.env` file or environment:

```env
BASE_URL=http://localhost:8081
```

## Running Tests

```bash
# Run all tests
npm test

# Run with UI mode (visual debugger)
npm run test:debug

# Run in headed mode (see browser)
npm run test:headed

# View HTML report
npm run test:report

# Run specific test file
npx playwright test specs/auth/auth.spec.ts

# Run tests matching pattern
npx playwright test --grep "should create"
```

## Benefits of This Architecture

### ✅ Type Safety

- Compile-time errors for invalid requests/responses
- Autocomplete for API fields
- Refactoring safety (rename fields, update all usages)

### ✅ Single Source of Truth

- Backend changes → regenerate types → tests show errors
- No manual type definitions
- API contract is always in sync

### ✅ Reusable Fixtures

- DRY principle (Don't Repeat Yourself)
- Shared authentication logic
- Consistent error handling

### ✅ Clear Test Organization

- Common tests → Quick verification of CRUD
- Case tests → Deep coverage of business logic
- Easy to find and maintain tests

### ✅ Fast Execution

- Global authentication (login once)
- Parallel execution (when safe)
- Isolated test contexts

## Next Steps

Continue reading:

- [02-fixtures.md](./02-fixtures.md) - Custom fixture system explained
- [03-test-organization.md](./03-test-organization.md) - Writing effective tests
- [04-type-generation.md](./04-type-generation.md) - OpenAPI type generation details
- [05-authentication.md](./05-authentication.md) - Authentication patterns
- [06-best-practices.md](./06-best-practices.md) - Testing best practices

## Summary

This E2E testing framework provides:

- **Type-safe API testing** using generated types from OpenAPI spec
- **Reusable fixtures** for consistent API client usage
- **Clear test organization** (common vs cases)
- **Fast execution** with global authentication
- **Comprehensive coverage** of CRUD and business logic

The architecture ensures that tests are maintainable, reliable, and always in sync with the backend API contract.
