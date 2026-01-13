import { test as base } from "@playwright/test";
import { AuthAPIClient } from "./auth-client";
import { ProjectAPIClient } from "./project-client";
import { StatusAPIClient } from "./status-client";
import { TaskAPIClient } from "./task-client";
import type { TestContext } from "../types/common";
import * as fs from "fs";
import * as path from "path";

/**
 * Extended test fixtures with API clients
 */
type APIFixtures = {
  testContext: TestContext;
  authAPI: AuthAPIClient;
  projectAPI: ProjectAPIClient;
  statusAPI: StatusAPIClient;
  taskAPI: TaskAPIClient;
  authenticatedContext: TestContext;
};

/**
 * Load authentication tokens from saved storage state
 */
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

/**
 * Extend Playwright test with custom fixtures
 */
export const test = base.extend<APIFixtures>({
  /**
   * Test context that carries shared state across fixtures
   */
  testContext: async ({ request }, use) => {
    // Load saved auth tokens from global setup
    const { accessToken, refreshToken } = loadAuthTokens();

    const context: TestContext = {
      baseURL: `http://localhost:8081`,
      accessToken,
      refreshToken,
    };
    await use(context);
  },

  /**
   * Auth API client
   */
  authAPI: async ({ request, testContext }, use) => {
    const client = new AuthAPIClient(request, testContext);
    await use(client);
    // Cleanup: logout after test
    client.logout();
  },

  /**
   * Project API client
   */
  projectAPI: async ({ request, testContext }, use) => {
    const client = new ProjectAPIClient(request, testContext);
    await use(client);
  },

  /**
   * Status API client
   */
  statusAPI: async ({ request, testContext }, use) => {
    const client = new StatusAPIClient(request, testContext);
    await use(client);
  },

  /**
   * Task API client
   */
  taskAPI: async ({ request, testContext }, use) => {
    const client = new TaskAPIClient(request, testContext);
    await use(client);
  },

  /**
   * Authenticated context - now automatically loaded from global setup
   * This fixture is kept for backward compatibility but tokens are
   * automatically available in testContext
   */
  authenticatedContext: async ({ testContext }, use) => {
    // Tokens are already loaded from storage state in testContext
    if (!testContext.accessToken) {
      throw new Error(
        "No authentication tokens found. Make sure global setup completed successfully."
      );
    }
    await use(testContext);
  },
});

export { expect } from "@playwright/test";
