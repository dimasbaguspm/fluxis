import { defineConfig } from "@playwright/test";
import dotenv from "dotenv";

dotenv.config();

/**
 * Playwright configuration for E2E API testing
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: "./specs",
  timeout: 30 * 1000,
  fullyParallel: false, // API tests should run sequentially to avoid conflicts
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: [
    ["html", { outputFolder: "playwright-report" }],
    ["json", { outputFile: "test-results/results.json" }],
    ["list"],
  ],

  use: {
    baseURL: `http://localhost:8081`,
    trace: "on-first-retry",
    extraHTTPHeaders: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    screenshot: "only-on-failure",
    video: "retain-on-failure",
  },

  projects: [
    {
      name: "api-tests",
      testMatch: "**/*.spec.ts",
      use: {
        storageState: ".auth/user.json",
      },
    },
  ],
  globalSetup: "./global-setup.ts",
});
