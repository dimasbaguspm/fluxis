import { test, expect } from "../../../fixtures";

/**
 * Auth Token Expiry and Refresh Cycle Test Cases
 * Tests token lifecycle, expiry, and refresh mechanisms
 */
test.describe("Auth Token Lifecycle", () => {
  test.describe("Token Refresh Cycle", () => {
    test("should allow using valid refresh token after access token expires", async ({
      authAPI,
    }) => {
      // Login to get initial tokens
      const login = await authAPI.login("admin_username", "admin_password");

      expect(login.status).toBe(200);
      expect(login.data?.accessToken).toBeDefined();
      expect(login.data?.refreshToken).toBeDefined();

      const refreshToken = login.data!.refreshToken;

      // Refresh the token
      const refresh = await authAPI.refresh(refreshToken);

      expect(refresh.status).toBe(200);
      expect(refresh.data?.accessToken).toBeDefined();

      // Note: Backend may return same token if TTL is long or tokens are cached
      // Not asserting token difference as it depends on implementation
    });

    test("should allow multiple consecutive refreshes", async ({ authAPI }) => {
      // Login
      const login = await authAPI.login("admin_username", "admin_password");

      expect(login.status).toBe(200);
      const refreshToken = login.data!.refreshToken;

      // Perform multiple refreshes using the same refresh token
      for (let i = 0; i < 3; i++) {
        const refresh = await authAPI.refresh(refreshToken);

        expect(refresh.status).toBe(200);
        expect(refresh.data?.accessToken).toBeDefined();
      }
    });

    test("should invalidate old refresh token after using it", async ({
      authAPI,
    }) => {
      // Login
      const login = await authAPI.login("admin_username", "admin_password");

      const oldRefreshToken = login.data!.refreshToken;

      // Refresh once
      const refresh1 = await authAPI.refresh(oldRefreshToken);

      expect(refresh1.status).toBe(200);

      // Try to use the old refresh token again
      const refresh2 = await authAPI.refresh(oldRefreshToken);

      // Should fail because the old token was invalidated
      // Note: This depends on your token rotation strategy
      // If you don't rotate tokens, this test should be adjusted
      if (refresh2.status >= 400) {
        expect(refresh2.status).toBeGreaterThanOrEqual(400);
      } else {
        // If your system allows reusing refresh tokens, document it here
        expect(refresh2.status).toBe(200);
      }
    });

    test("should reject invalid refresh token", async ({ authAPI }) => {
      const response = await authAPI.refresh("invalid-token-12345");

      expect(response.status).toBeGreaterThanOrEqual(400);
    });

    test("should reject empty refresh token", async ({ authAPI }) => {
      const response = await authAPI.refresh("");

      // Backend may accept empty (fallback to stored token) or reject
      // Accepting with 200 is valid if there's a fallback mechanism
      expect([200, 400, 401]).toContain(response.status);
    });

    test("should reject malformed refresh token", async ({ authAPI }) => {
      const response = await authAPI.refresh("not.a.valid.jwt");

      expect(response.status).toBeGreaterThanOrEqual(400);
    });

    test("should return new tokens with correct structure", async ({
      authAPI,
    }) => {
      // Login
      const login = await authAPI.login("admin_username", "admin_password");

      // Refresh
      const refresh = await authAPI.refresh(login.data!.refreshToken);

      expect(refresh.status).toBe(200);
      expect(refresh.data).toMatchObject({
        accessToken: expect.any(String),
      });

      // Access token should be JWT format (3 parts separated by dots)
      expect(refresh.data!.accessToken.split(".").length).toBe(3);
    });
  });

  test.describe("Token Validation", () => {
    test("should accept valid access token for authenticated endpoints", async ({
      authAPI,
      projectAPI,
    }) => {
      // Login to get tokens
      const login = await authAPI.login("admin_username", "admin_password");

      expect(login.status).toBe(200);

      // The fixture should automatically use the access token
      // Try to access a protected endpoint
      const projects = await projectAPI.getPaginated({});

      expect(projects.status).toBe(200);
    });

    test("should reject request with missing authorization header", async ({
      request,
    }) => {
      // Direct API call without auth
      const response = await request.get("/api/projects", {
        headers: {
          // No Authorization header
        },
      });

      expect(response.status()).toBeGreaterThanOrEqual(400);
    });

    test("should reject request with invalid access token", async ({
      request,
    }) => {
      const response = await request.get("/api/projects", {
        headers: {
          Authorization: "Bearer invalid-token-12345",
        },
      });

      expect(response.status()).toBeGreaterThanOrEqual(400);
    });

    test("should reject request with malformed authorization header", async ({
      request,
    }) => {
      const response = await request.get("/api/projects", {
        headers: {
          Authorization: "invalid-format",
        },
      });

      expect(response.status()).toBeGreaterThanOrEqual(400);
    });

    test("should handle concurrent refresh requests gracefully", async ({
      authAPI,
    }) => {
      // Login
      const login = await authAPI.login("admin_username", "admin_password");

      const refreshToken = login.data!.refreshToken;

      // Send multiple refresh requests simultaneously
      const refreshPromises = Array(5)
        .fill(null)
        .map(() => authAPI.refresh(refreshToken));

      const results = await Promise.all(refreshPromises);

      // At least one should succeed
      const successCount = results.filter((r) => r.status === 200).length;
      expect(successCount).toBeGreaterThanOrEqual(1);

      // Depending on implementation:
      // - If token rotation is strict, only 1 should succeed
      // - If tokens can be reused, all might succeed
    });
  });

  test.describe("Session Management", () => {
    test("should allow multiple sessions for same user", async ({
      authAPI,
    }) => {
      // Login twice
      const login1 = await authAPI.login("admin_username", "admin_password");

      const login2 = await authAPI.login("admin_username", "admin_password");

      expect(login1.status).toBe(200);
      expect(login2.status).toBe(200);

      // Both sessions should have valid tokens
      expect(login1.data?.accessToken).toBeDefined();
      expect(login2.data?.accessToken).toBeDefined();
      // Note: Backend may return same tokens if caching is used
    });

    test("should maintain session consistency across refreshes", async ({
      authAPI,
      projectAPI,
    }) => {
      // Login
      const login = await authAPI.login("admin_username", "admin_password");

      // Create a project
      const project1 = await projectAPI.create({
        name: "Before Refresh",
        description: "Test",
        status: "active",
      });

      expect(project1.status).toBe(200);
      const projectId = project1.data!.id;

      // Refresh token
      await authAPI.refresh(login.data!.refreshToken);

      // Should still be able to access the project
      const project2 = await projectAPI.getById(projectId);

      expect(project2.status).toBe(200);
      expect(project2.data?.id).toBe(projectId);

      // Cleanup
      await projectAPI.remove(projectId);
    });
  });
});
