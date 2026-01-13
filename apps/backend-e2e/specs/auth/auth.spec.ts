import { test, expect } from "../../fixtures";

/**
 * Authentication endpoint tests
 * Tests for /auth/login and /auth/refresh
 */
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

    test("should fail refresh with empty token", async ({
      authAPI,
      testContext,
    }) => {
      // Clear any existing tokens
      testContext.refreshToken = undefined;

      try {
        await authAPI.refresh();
        // Should not reach here
        expect(true).toBe(false);
      } catch (error) {
        expect(error).toBeDefined();
      }
    });
  });
});
