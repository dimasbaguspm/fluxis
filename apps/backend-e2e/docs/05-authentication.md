# Authentication in E2E Tests

## Overview

Authentication is handled through **global setup** to avoid repeated login requests in every test. Tests automatically load authentication tokens from a shared storage state.

## Authentication Flow

```
┌──────────────────────────┐
│  Global Setup            │
│  (runs once)             │
│  1. Login via API        │
│  2. Save tokens to file  │
└──────────┬───────────────┘
           │
           │ Writes to
           ▼
┌──────────────────────────┐
│  .auth/user.json         │  ← Storage state with tokens
│  {                       │
│    origins: [{           │
│      localStorage: [     │
│        { accessToken }  │
│        { refreshToken } │
│      ]                   │
│    }]                    │
│  }                       │
└──────────┬───────────────┘
           │
           │ Loaded by fixtures
           ▼
┌──────────────────────────┐
│  testContext fixture     │
│  - Reads .auth/user.json │
│  - Extracts tokens       │
│  - Provides to tests     │
└──────────┬───────────────┘
           │
           │ Used by
           ▼
┌──────────────────────────┐
│  API Client fixtures     │
│  - authAPI               │
│  - projectAPI            │
│  - taskAPI (all authenticated)
└──────────┬───────────────┘
           │
           │ Injected into
           ▼
┌──────────────────────────┐
│  Test specs              │
│  - All tests start       │
│    authenticated         │
└──────────────────────────┘
```

## Global Setup

Located at: `global-setup.ts`

### Purpose

Perform **one-time authentication** before all tests run:

1. Login with default credentials
2. Extract access and refresh tokens
3. Save tokens to `.auth/user.json` storage state
4. All tests automatically load tokens from this file

### Implementation

```typescript
import { request, FullConfig } from "@playwright/test";
import dotenv from "dotenv";
import * as fs from "fs";
import * as path from "path";

async function globalSetup(config: FullConfig) {
  dotenv.config();

  const baseURL = `http://localhost:8081`;
  const username = "admin_username";
  const password = "admin_password";

  console.log("Performing global authentication setup...");

  const requestContext = await request.newContext({
    baseURL,
    extraHTTPHeaders: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
  });

  try {
    // Perform login
    const response = await requestContext.post("/auth/login", {
      data: {
        username,
        password,
      },
    });

    if (!response.ok()) {
      const error = await response.json();
      throw new Error(
        `Authentication failed: ${error.detail || response.statusText()}`
      );
    }

    const data = await response.json();

    if (!data.accessToken || !data.refreshToken) {
      throw new Error("Login response missing tokens");
    }

    console.log("Authentication successful");

    // Create .auth directory
    const authDir = path.join(process.cwd(), ".auth");
    if (!fs.existsSync(authDir)) {
      fs.mkdirSync(authDir, { recursive: true });
    }

    // Save tokens to storage state
    const storageState = {
      cookies: [],
      origins: [
        {
          origin: baseURL,
          localStorage: [
            {
              name: "accessToken",
              value: data.accessToken,
            },
            {
              name: "refreshToken",
              value: data.refreshToken,
            },
          ],
        },
      ],
    };

    const authStatePath = path.join(authDir, "user.json");
    fs.writeFileSync(authStatePath, JSON.stringify(storageState, null, 2));

    console.log("Auth tokens saved to .auth/user.json");
  } catch (error) {
    console.error("Global setup failed:", error);
    throw error;
  } finally {
    await requestContext.dispose();
  }
}

export default globalSetup;
```

### Storage State Format

The `.auth/user.json` file contains:

```json
{
  "cookies": [],
  "origins": [
    {
      "origin": "http://localhost:8081",
      "localStorage": [
        {
          "name": "accessToken",
          "value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        },
        {
          "name": "refreshToken",
          "value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        }
      ]
    }
  ]
}
```

## Playwright Configuration

In `playwright.config.ts`:

```typescript
export default defineConfig({
  globalSetup: "./global-setup.ts", // Run before all tests

  projects: [
    {
      name: "api-tests",
      testMatch: "**/*.spec.ts",
      use: {
        storageState: ".auth/user.json", // Load authentication state
      },
    },
  ],
});
```

## Loading Tokens in Fixtures

Located at: `fixtures/index.ts`

### loadAuthTokens Function

```typescript
function loadAuthTokens(): { accessToken?: string; refreshToken?: string } {
  try {
    const authStatePath = path.join(process.cwd(), ".auth/user.json");
    if (fs.existsSync(authStatePath)) {
      const storageState = JSON.parse(fs.readFileSync(authStatePath, "utf-8"));

      // Extract tokens from localStorage in the storage state
      const origins = storageState.origins || [];
      for (const origin of origins) {
        if (origin.localStorage) {
          const accessToken = origin.localStorage.find(
            (item: any) => item.name === "accessToken"
          )?.value;
          const refreshToken = origin.localStorage.find(
            (item: any) => item.name === "refreshToken"
          )?.value;

          if (accessToken && refreshToken) {
            return { accessToken, refreshToken };
          }
        }
      }
    }
  } catch (error) {
    console.warn("Could not load auth tokens from storage state:", error);
  }
  return {};
}
```

### testContext Fixture

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

**Key points**:

- Runs for every test
- Loads tokens from `.auth/user.json`
- Provides authenticated context to all fixtures
- No additional login required

## Using Authentication in Tests

### Authenticated by Default

All tests automatically have access tokens:

```typescript
import { test, expect } from "../../fixtures";

test("should get projects", async ({ projectAPI }) => {
  // projectAPI automatically includes Bearer token in headers
  const response = await projectAPI.getPaginated();

  expect(response.status).toBe(200);
  // No login required - token already present
});
```

### Testing Unauthenticated Requests

To test endpoints without authentication:

```typescript
test("should fail without authentication", async ({ request, testContext }) => {
  // Create client WITHOUT token
  const unauthContext = { ...testContext, accessToken: undefined };
  const client = new ProjectAPIClient(request, unauthContext);

  const response = await client.getPaginated();

  expect(response.status).toBe(401);
  expect(response.error).toBeDefined();
});
```

### Re-authenticating in Tests

For testing authentication flows:

```typescript
test("should login and use new token", async ({ authAPI, projectAPI }) => {
  // Logout (clears stored tokens)
  authAPI.logout();

  // Login with new credentials
  const loginResponse = await authAPI.login("testuser", "testpass");
  expect(loginResponse.status).toBe(200);

  // Subsequent requests use new token
  const projectsResponse = await projectAPI.getPaginated();
  expect(projectsResponse.status).toBe(200);
});
```

## Token Management in API Clients

### Automatic Token Injection

In `BaseAPIClient`:

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

**Every HTTP request automatically includes**:

```http
GET /projects HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

### Token Updates After Login

In `AuthAPIClient`:

```typescript
async login(
  username: string,
  password: string
): Promise<APIResponse<LoginResponseModel>> {
  const response = await this.post<LoginResponseModel>("/auth/login", {
    username,
    password,
  });

  // Auto-update context with new tokens
  if (response.data) {
    this.context.accessToken = response.data.accessToken;
    this.context.refreshToken = response.data.refreshToken;
  }

  return response;
}
```

### Token Refresh

```typescript
async refresh(
  refreshToken?: string
): Promise<APIResponse<RefreshResponseModel>> {
  const token = refreshToken || this.context.refreshToken;
  if (!token) {
    throw new Error("No refresh token available");
  }

  const response = await this.post<RefreshResponseModel>("/auth/refresh", {
    refreshToken: token,
  });

  // Auto-update access token
  if (response.data) {
    this.context.accessToken = response.data.accessToken;
  }

  return response;
}
```

## Testing Authentication Scenarios

### 1. Valid Login

```typescript
test("should login with valid credentials", async ({ authAPI }) => {
  const response = await authAPI.login("admin_username", "admin_password");

  expect(response.status).toBe(200);
  expect(response.data?.accessToken).toBeDefined();
  expect(response.data?.refreshToken).toBeDefined();
});
```

### 2. Invalid Credentials

```typescript
test("should reject invalid credentials", async ({ authAPI }) => {
  const response = await authAPI.login("invalid", "credentials");

  expect(response.status).toBeGreaterThanOrEqual(400);
  expect(response.error).toBeDefined();
});
```

### 3. Token Refresh

```typescript
test("should refresh access token", async ({ authAPI }) => {
  // Login first
  const loginResponse = await authAPI.login("admin_username", "admin_password");
  const oldAccessToken = loginResponse.data!.accessToken;
  const refreshToken = loginResponse.data!.refreshToken;

  // Refresh token
  const refreshResponse = await authAPI.refresh(refreshToken);

  expect(refreshResponse.status).toBe(200);
  expect(refreshResponse.data?.accessToken).toBeDefined();
  expect(refreshResponse.data?.accessToken).not.toBe(oldAccessToken);
});
```

### 4. Expired Token

```typescript
test("should reject expired token", async ({ projectAPI, testContext }) => {
  // Set an expired token
  testContext.accessToken = "expired.token.here";

  const response = await projectAPI.getPaginated();

  expect(response.status).toBe(401);
  expect(response.error?.detail).toContain("expired");
});
```

### 5. Missing Token

```typescript
test("should reject request without token", async ({
  request,
  testContext,
}) => {
  // Create unauthenticated client
  const unauthContext = { ...testContext, accessToken: undefined };
  const client = new ProjectAPIClient(request, unauthContext);

  const response = await client.getPaginated();

  expect(response.status).toBe(401);
});
```

### 6. Invalid Token Format

```typescript
test("should reject malformed token", async ({ projectAPI, testContext }) => {
  testContext.accessToken = "not-a-valid-jwt";

  const response = await projectAPI.getPaginated();

  expect(response.status).toBe(401);
});
```

## Authorization Testing

Test role-based access control:

```typescript
test.describe("Authorization", () => {
  test("admin should access all projects", async ({ authAPI, projectAPI }) => {
    await authAPI.login("admin", "admin_password");

    const response = await projectAPI.getPaginated();

    expect(response.status).toBe(200);
    expect(response.data?.items).toBeDefined();
  });

  test("regular user should only see own projects", async ({
    authAPI,
    projectAPI,
  }) => {
    await authAPI.login("user1", "user_password");

    const response = await projectAPI.getPaginated();

    expect(response.status).toBe(200);
    expect(response.data?.items.every((p) => p.ownerId === "user1")).toBe(true);
  });

  test("guest should not access projects", async ({ authAPI, projectAPI }) => {
    await authAPI.login("guest", "guest_password");

    const response = await projectAPI.getPaginated();

    expect(response.status).toBe(403); // Forbidden
  });
});
```

## Troubleshooting

### Issue: All tests fail with 401

**Symptom**: Tests immediately fail with unauthorized errors

**Cause**: Global setup didn't run or failed

**Solution**:

```bash
# Check if .auth/user.json exists
ls .auth/user.json

# Run global setup manually
npx playwright test --project=setup

# Or re-run all tests (global setup will run)
npm test
```

### Issue: Token expired during test run

**Symptom**: Tests fail midway through with 401 errors

**Cause**: Access token expired (typical lifetime: 15-30 minutes)

**Solution 1** - Use refresh token:

```typescript
test.beforeEach(async ({ authAPI, testContext }) => {
  // Refresh token before each test
  if (testContext.refreshToken) {
    await authAPI.refresh(testContext.refreshToken);
  }
});
```

**Solution 2** - Extend token lifetime in backend (for testing):

```go
// backend/internal/repositories/auth_repository.go
const AccessTokenExpiry = 24 * time.Hour  // Longer expiry for tests
```

### Issue: Tests work locally but fail in CI

**Symptom**: Authentication works on dev machine, fails in CI/CD

**Cause**: Environment variables or credentials not set in CI

**Solution**:

```yaml
# .github/workflows/test.yml
env:
  BASE_URL: http://localhost:8081
  ADMIN_USERNAME: ${{ secrets.ADMIN_USERNAME }}
  ADMIN_PASSWORD: ${{ secrets.ADMIN_PASSWORD }}

steps:
  - name: Run Tests
    run: npm test
```

## Best Practices

### ✅ DO

- **Use global setup** for initial authentication
- **Load tokens from storage state** in fixtures
- **Auto-inject Bearer tokens** in API clients
- **Test both authenticated and unauthenticated** scenarios
- **Refresh tokens** for long-running test suites
- **Clear tokens** in test cleanup if needed

### ❌ DON'T

- **Don't login in every test** - Wastes time and API calls
- **Don't hardcode tokens** - Use storage state
- **Don't share tokens across users** - Isolate test contexts
- **Don't ignore 401 errors** - They indicate auth issues
- **Don't skip authentication tests** - Critical for security

## Summary

Authentication in E2E tests:

- **Global setup** - Login once, save tokens to `.auth/user.json`
- **Automatic loading** - Fixtures load tokens from storage state
- **Transparent injection** - API clients automatically include Bearer tokens
- **Test isolation** - Each test can re-authenticate if needed
- **Type-safe** - Auth responses use generated OpenAPI types

**Benefits**:

- **Fast test execution** - No repeated login calls
- **Reliable** - Tokens managed automatically
- **Flexible** - Easy to test different auth scenarios
- **Maintainable** - Auth logic centralized in fixtures

This architecture ensures tests are fast, reliable, and easy to maintain while thoroughly testing authentication and authorization flows.
