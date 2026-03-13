#!/usr/bin/env node

import { execSync } from "child_process";
import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const projectRoot = path.resolve(__dirname, "..");
const swaggerPath = path.resolve(projectRoot, "..", "api", "swagger.json");
const outputDir = path.resolve(projectRoot, "src", "types");
const outputFile = path.resolve(outputDir, "openapi.generated.ts");
const tempApiFile = path.resolve(outputDir, "Api.ts");

try {
  // Ensure output directory exists
  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }

  console.log("Generating types from api/swagger.json...");
  execSync(`npx swagger-typescript-api generate -p "${swaggerPath}" -o "${outputDir}"`, {
    stdio: "inherit",
    cwd: projectRoot,
  });

  // Rename generated file
  if (fs.existsSync(tempApiFile)) {
    fs.renameSync(tempApiFile, outputFile);
  }

  console.log(`✓ Generated types at ${path.relative(projectRoot, outputFile)}`);
} catch {
  console.error("✗ Failed to generate types");
  process.exit(1);
}
