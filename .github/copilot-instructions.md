# Fluxis - GitHub Copilot Instructions

## Project Overview

**Fluxis** is a project management system built as a **monorepo** containing multiple applications. This repository follows a clean architecture pattern with clear separation of concerns.

## Repository Structure

This is a **monorepo** containing **4 main applications**:

```
apps/
├── backend/          # Go REST API backend
├── backend-e2e/      # End-to-end testing for backend
├── frontend-tui/     # Terminal User Interface (TUI) frontend
└── frontend-web/     # Web frontend application
```

---

## Documentation by Application

### When Working on Backend (`apps/backend/`)

**DOCUMENTATION PATH:** `apps/backend/docs/`

**CRITICAL: Always retrieve and read documentation from `apps/backend/docs/` before implementing any changes.**

**Available documentation files (in order of reference):**

1. **`01-architecture.md`** - System architecture, tech stack, and design decisions
   - Reference when: Understanding overall system design, tech decisions
2. **`02-code-flow.md`** - Request lifecycle and execution flow
   - Reference when: Understanding how requests are processed end-to-end
3. **`03-rsrm-pattern.md`** - Resource-Service-Repository-Model pattern explanation
   - Reference when: Implementing new features, creating layers
4. **`04-middlewares.md`** - Middleware system and JWT authentication
   - Reference when: Adding authentication, cross-cutting concerns, middleware
5. **`05-main-and-handler.md`** - Application entry point and route registration
   - Reference when: Registering new routes, understanding initialization
6. **`06-migrations-workflow.md`** - Database migrations system and workflow
   - Reference when: Making database schema changes
7. **`07-concurrency.md`** - Goroutine patterns and concurrency
   - Reference when: Working with concurrent operations, parallel reads
8. **`08-workers.md`** - Worker pattern for async event processing
   - Reference when: Creating workers, processing entity changes, asynchronous logging

## Key Implementation Patterns:

- Follow the **RSRM** (Resource-Service-Repository-Model) pattern strictly
- Implementation order: Model → Repository → Service → Resource
- Use constructor functions for dependency injection
- Apply middleware for cross-cutting concerns

**Code Structure Reference:**

```
apps/backend/internal/
├── models/          # Data structures
├── repositories/    # Database access layer
├── services/        # Business logic layer
├── resources/       # HTTP handlers layer
├── handler.go       # Route registration
├── middlewares/     # Middleware implementations
└── workers/         # Worker logic for async processing
```

### When Working on Backend E2E (`apps/backend-e2e/`)

**DOCUMENTATION PATH:** `apps/backend-e2e/docs/`

**CRITICAL: Always retrieve and read documentation from `apps/backend-e2e/docs/` before implementing tests.**

**Available documentation files:**

1. **`01-overview.md`** - E2E testing architecture and setup overview
2. **`02-fixtures.md`** - Test fixtures and helper utilities
3. **`03-test-organization.md`** - How tests are structured and organized
4. **`04-type-generation.md`** - OpenAPI type generation for tests
5. **`05-authentication.md`** - Testing authentication flows
6. **`06-best-practices.md`** - E2E testing best practices

**Test Structure Reference:**

```
apps/backend-e2e/
├── fixtures/        # Test utilities and API clients
├── specs/           # Test specifications organized by feature
├── types/           # Generated OpenAPI types
└── scripts/         # Utility scripts (type generation)
```

### When Working on Frontend TUI (`apps/frontend-tui/`)

**DOCUMENTATION PATH:** `apps/frontend-tui/docs/`

**⚡ Always retrieve documentation from `apps/frontend-tui/docs/` before implementing features.**

Terminal User Interface (TUI) frontend application.

### When Working on Frontend Web (`apps/frontend-web/`)

**DOCUMENTATION PATH:** `apps/frontend-web/docs/`

**Always retrieve documentation from `apps/frontend-web/docs/` before implementing features.**

Web-based user interface application.

---

## General Guidelines

### Monorepo Navigation

- Each app is **independent** with its own structure and documentation
- Check the current working directory to determine which app you're working on
- Documentation is **app-specific** - always check the right `docs/` folder
- Shared types/contracts may exist between apps (coordinate when needed)

### Code Consistency

- Follow established patterns within each app
- Maintain consistent naming conventions
- Keep dependency injection explicit
- Write self-documenting code with clear function/variable names

### Development Workflow

1. **Read the docs first** - Check the relevant `apps/*/docs/` directory
2. **Follow existing patterns** - Look at similar implementations
3. **Keep it simple** - Don't over-engineer for solo development
4. **Commit often** - Small, focused commits are easier to debug
5. **Test manually** - Use tools like OpenAPI docs (`/docs`) for backend

### Architecture Principles

- **Separation of Concerns** - Each layer has one responsibility
- **Dependency Injection** - Use constructor functions
- **Clean Architecture** - Business logic independent of frameworks
- **Fail Fast** - Validate early, return errors clearly
- **Documentation** - Code should explain "why", comments explain "what"

## Context-Aware Guidelines

**BEFORE TAKING ANY ACTION - Always perform these steps:**

1. **Identify the target application** - Which app are you working in?

   - Backend? → Read from `apps/backend/docs/`
   - Backend E2E? → Read from `apps/backend-e2e/docs/`
   - Frontend TUI? → Read from `apps/frontend-tui/docs/`
   - Frontend Web? → Read from `apps/frontend-web/docs/`

2. **Retrieve relevant documentation** - Use semantic_search or read_file tools to pull context from the appropriate docs directory

3. **Match existing patterns** - Look at similar implementations in the same app before suggesting changes

4. **Reference the docs in your reasoning** - Cite which documentation files support your approach

### When Suggesting Code Changes:

- ✅ Match the existing code style in the current app
- ✅ Explicitly reference which doc file supports the pattern
- ✅ Follow established patterns (especially **RSRM for backend**)
- ✅ Consider the monorepo structure
- ✅ Suggest small, incremental changes for solo development

### When Explaining Architecture:

- ✅ Point to specific docs in `apps/*/docs/`
- ✅ Show examples from existing code in the same app
- ✅ Explain trade-offs for solo vs team development
- ✅ Keep explanations practical and actionable

### When Debugging Issues:

- ✅ **First:** Read relevant docs from `apps/*/docs/` directory
- ✅ **Second:** Look at similar working implementations in the same app
- ✅ **Third:** Suggest adding logs for troubleshooting
- ✅ **Fourth:** Consider the data flow through the layers

---

## Notes

- This is a **solo development project** - optimize for velocity and learning
- Tests, security hardening, and production features can be added incrementally
- Focus on getting features working end-to-end before optimization
- The architecture is solid - follow established patterns for consistency
