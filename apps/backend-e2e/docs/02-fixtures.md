# Custom Fixtures System

## Overview

Fixtures in Playwright extend the test context with reusable utilities, API clients, and shared state. The backend-e2e project uses **custom fixtures** to provide type-safe API clients to every test.

## Fixture Architecture

```
┌─────────────────────────────────────────┐
│      Playwright Base Test               │
│      (test, expect, request)            │
└──────────────┬──────────────────────────┘
               │
               │ base.extend<APIFixtures>()
               │
┌──────────────▼──────────────────────────┐
│      Custom Test Extension              │
│      (testContext, authAPI, ...)        │
└──────────────┬──────────────────────────┘
               │
               │ import { test } from "fixtures"
               │
┌──────────────▼──────────────────────────┐
│      Test Specs                         │
│      (specs/**/*.spec.ts)               │
└─────────────────────────────────────────┘
```

## Base Client (`fixtures/base-client.ts`)

The `BaseAPIClient` provides **common HTTP utilities** that all API clients inherit.

### Key Responsibilities

1. **HTTP Methods** - GET, POST, PATCH, PUT, DELETE
2. **Authentication** - Automatic bearer token injection
3. **Response Parsing** - JSON parsing and error handling
4. **Assertions** - Success/error response validation

### Implementation Details

#### Constructor

```typescript
export class BaseAPIClient {
  constructor(
    protected request: APIRequestContext,
    protected context: TestContext
  ) {}
}
```

**Parameters**:

- `request` - Playwright's API request context
- `context` - Shared test context (holds tokens, baseURL)

#### Authentication Headers

```typescript
protected getAuthHeaders(): Record<string, string> {
  if (this.context.accessToken) {
    return {
      Authorization: `Bearer ${this.context.accessToken}`,
    };
  }
  return {};
}
```

**Behavior**:

- If `accessToken` exists in context → Add `Authorization` header
- If no token → Return empty headers (for public endpoints)

#### HTTP Methods

All methods follow the same pattern:

1. Build full URL from path
2. Add authentication headers
3. Make request with Playwright's `request` context
4. Parse and return typed response

**Example: POST Request**

```typescript
protected async post<T>(path: string, body?: any): Promise<APIResponse<T>> {
  const url = new URL(path, this.context.baseURL);
  const response = await this.request.post(url.toString(), {
    headers: this.getAuthHeaders(),
    data: body,
  });
  return this.parseResponse<T>(response);
}
```

#### Response Parsing

```typescript
private async parseResponse<T>(response: any): Promise<APIResponse<T>> {
  const status = response.status();
  const headers: Record<string, string> = {};

  // Extract headers
  response.headersArray().forEach((header: { name: string; value: string }) => {
    headers[header.name] = header.value;
  });

  let data: T | undefined;
  let error: any;

  try {
    const body = await response.json();
    if (status >= 200 && status < 300) {
      data = body as T;
    } else {
      error = body;
    }
  } catch (e) {
    // No body or non-JSON response
    if (status >= 200 && status < 300) {
      data = undefined;
    }
  }

  return { data, error, status, headers };
}
```

**Key Features**:

- Handles both success (2xx) and error (4xx, 5xx) responses
- Gracefully handles missing body (e.g., 204 No Content)
- Returns structured response with data/error separation

#### Assertion Helpers

```typescript
protected assertSuccess<T>(
  response: APIResponse<T>
): asserts response is APIResponse<T> & { data: T } {
  expect(response.status).toBeGreaterThanOrEqual(200);
  expect(response.status).toBeLessThan(300);
  expect(response.data).toBeDefined();
}

protected assertError(
  response: APIResponse<any>,
  expectedStatus?: number
): void {
  if (expectedStatus) {
    expect(response.status).toBe(expectedStatus);
  } else {
    expect(response.status).toBeGreaterThanOrEqual(400);
  }
}
```

## Entity API Clients

Each API entity (auth, projects, tasks, etc.) has a dedicated client extending `BaseAPIClient`.

### Example: Auth API Client (`fixtures/auth-client.ts`)

```typescript
export class AuthAPIClient extends BaseAPIClient {
  constructor(request: APIRequestContext, context: TestContext) {
    super(request, context);
  }

  async login(
    username: string,
    password: string
  ): Promise<APIResponse<LoginResponseModel>> {
    const response = await this.post<LoginResponseModel>("/auth/login", {
      username,
      password,
    });

    // Auto-save tokens to context
    if (response.data) {
      this.context.accessToken = response.data.accessToken;
      this.context.refreshToken = response.data.refreshToken;
    }

    return response;
  }

  async refresh(
    refreshToken?: string
  ): Promise<APIResponse<RefreshResponseModel>> {
    const token = refreshToken || this.context.refreshToken;
    if (!token) {
      throw new Error("No refresh token available");
    }

    const response = await this.post<RefreshResponseModel>("/auth/refresh", {
      refresh_token: token,
    });

    // Auto-update access token
    if (response.data) {
      this.context.accessToken = response.data.accessToken;
    }

    return response;
  }

  logout(): void {
    this.context.accessToken = undefined;
    this.context.refreshToken = undefined;
  }
}
```

**Key Patterns**:

1. **Type-safe requests** - Uses generated OpenAPI types
2. **Automatic token management** - Saves tokens to context after login
3. **Domain-specific methods** - `login()`, `refresh()`, `logout()`
4. **Consistent error handling** - Inherits from base client

## Fixture Registration (`fixtures/index.ts`)

This file **extends Playwright's test** with custom fixtures.

### Test Context Fixture

```typescript
testContext: async ({ request }, use) => {
  // Load saved auth tokens from global setup
  const { accessToken, refreshToken } = loadAuthTokens();

  const context: TestContext = {
    baseURL: `http://localhost:8081`,
    accessToken,
    refreshToken,
  };

  await use(context);
};
```

**Purpose**: Create shared context with authentication tokens

**Lifecycle**:

1. Load tokens from `.auth/user.json` (saved by global setup)
2. Create context object with baseURL and tokens
3. Provide context to fixtures and tests
4. Context persists throughout test execution

### API Client Fixtures

```typescript
authAPI: async ({ request, testContext }, use) => {
  const client = new AuthAPIClient(request, testContext);
  await use(client);
  // Cleanup: logout after test
  client.logout();
};
```

**Pattern for all API clients**:

1. **Setup**: Create client with `request` and `testContext`
2. **Use**: Provide client to test
3. **Teardown**: Cleanup (logout, clear state)

### Authenticated Context Fixture

```typescript
authenticatedContext: async ({ testContext }, use) => {
  // Tokens are already loaded from storage state
  if (!testContext.accessToken) {
    throw new Error(
      "No authentication tokens found. Make sure global setup completed successfully."
    );
  }
  await use(testContext);
};
```

**Purpose**: Validate that authentication is present

**Use case**: Tests that require guaranteed authentication

## Creating New API Clients

Follow this pattern to add new entity clients:

### Step 1: Define Types

```typescript
import type { components } from "../types/openapi";

export type ProjectCreateRequest = components["schemas"]["ProjectCreateModel"];
export type ProjectResponse = components["schemas"]["ProjectModel"];
export type ProjectPaginatedResponse =
  components["schemas"]["ProjectPaginatedModel"];
export type ProjectUpdateRequest = components["schemas"]["ProjectUpdateModel"];
```

### Step 2: Create Client Class

```typescript
export class ProjectAPIClient extends BaseAPIClient {
  constructor(request: APIRequestContext, context: TestContext) {
    super(request, context);
  }

  async create(
    data: ProjectCreateRequest
  ): Promise<APIResponse<ProjectResponse>> {
    return this.post<ProjectResponse>("/projects", data);
  }

  async getById(id: string): Promise<APIResponse<ProjectResponse>> {
    return this.get<ProjectResponse>(`/projects/${id}`);
  }

  async getPaginated(params?: {
    page?: number;
    limit?: number;
  }): Promise<APIResponse<ProjectPaginatedResponse>> {
    return this.get<ProjectPaginatedResponse>("/projects", params);
  }

  async update(
    id: string,
    data: ProjectUpdateRequest
  ): Promise<APIResponse<ProjectResponse>> {
    return this.patch<ProjectResponse>(`/projects/${id}`, data);
  }

  async delete(id: string): Promise<APIResponse<void>> {
    return this.delete<void>(`/projects/${id}`);
  }
}
```

### Step 3: Register Fixture

```typescript
// In fixtures/index.ts
type APIFixtures = {
  testContext: TestContext;
  authAPI: AuthAPIClient;
  projectAPI: ProjectAPIClient; // Add new fixture type
};

export const test = base.extend<APIFixtures>({
  // ... existing fixtures ...

  projectAPI: async ({ request, testContext }, use) => {
    const client = new ProjectAPIClient(request, testContext);
    await use(client);
  },
});
```

### Step 4: Use in Tests

```typescript
import { test, expect } from "../../fixtures";

test("should create project", async ({ projectAPI }) => {
  const response = await projectAPI.create({
    name: "Test Project",
    description: "Test description",
    status: "active",
  });

  expect(response.status).toBe(200);
  expect(response.data?.name).toBe("Test Project");
});
```

## Fixture Benefits

### ✅ Type Safety

```typescript
// TypeScript enforces correct types
const response = await authAPI.login("user", "pass");
//    ^? APIResponse<LoginResponseModel>

// Autocomplete available
response.data?.accessToken;
//              ^? string
```

### ✅ Reusability

```typescript
// Same client used across all tests
test("test 1", async ({ authAPI }) => {
  await authAPI.login("user1", "pass1");
});

test("test 2", async ({ authAPI }) => {
  await authAPI.login("user2", "pass2");
});
```

### ✅ Automatic Cleanup

```typescript
// Fixtures handle setup/teardown
authAPI: async ({ request, testContext }, use) => {
  const client = new AuthAPIClient(request, testContext);
  await use(client); // Test runs here
  client.logout(); // Automatic cleanup
};
```

### ✅ Dependency Injection

```typescript
// Tests receive only what they need
test("auth test", async ({ authAPI }) => {
  // Only authAPI is injected
});

test("project test", async ({ projectAPI, authAPI }) => {
  // Both clients available
});
```

### ✅ Consistent Error Handling

All clients inherit the same error handling logic from `BaseAPIClient`:

```typescript
const response = await projectAPI.create({
  /* invalid data */
});

if (response.status >= 400) {
  console.log(response.error); // Structured error object
}
```

## Best Practices

### DO ✅

- **Extend BaseAPIClient** for all entity clients
- **Use generated types** from `types/openapi.ts`
- **Return APIResponse<T>** for consistent response structure
- **Handle tokens automatically** in auth-related methods
- **Provide cleanup** in fixture teardown

### DON'T ❌

- **Don't bypass fixtures** - Use API clients, not direct `request.post()`
- **Don't hardcode URLs** - Use paths relative to `baseURL`
- **Don't ignore errors** - Always check `response.status` or use assertions
- **Don't share state** between tests via fixtures
- **Don't manually manage tokens** - Let fixtures handle it

## Summary

The fixture system provides:

- **BaseAPIClient** - Reusable HTTP methods and utilities
- **Entity clients** - Domain-specific API clients (AuthAPIClient, ProjectAPIClient, etc.)
- **Fixture registration** - Extend Playwright test with custom fixtures
- **Type safety** - Generated OpenAPI types throughout
- **Automatic cleanup** - Setup/teardown lifecycle management

This architecture ensures tests are:

- **Type-safe** - Compile-time checks for API contracts
- **Maintainable** - Reusable clients, DRY principle
- **Consistent** - Same patterns across all tests
- **Reliable** - Automatic authentication and cleanup
