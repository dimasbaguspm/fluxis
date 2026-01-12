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

**Always refer to:** `apps/backend/docs/`

Available documentation:

- `01-architecture.md` - System architecture, tech stack, and design decisions
- `02-code-flow.md` - Request lifecycle and execution flow
- `03-rsrm-pattern.md` - Resource-Service-Repository-Model pattern explanation
- `04-middlewares.md` - Middleware system and JWT authentication
- `05-main-and-handler.md` - Application entry point and route registration
- `06-migrations-workflow.md` - Database migrations system and workflow

**Key Patterns:**

- Follow the **RSRM** (Resource-Service-Repository-Model) pattern
- All new features should implement: Model → Repository → Service → Resource
- Use constructor functions for dependency injection
- Apply middleware for cross-cutting concerns

### When Working on Backend E2E (`apps/backend-e2e/`)

**Always refer to:** `apps/backend-e2e/docs/`

This app contains end-to-end tests for the backend API.

### When Working on Frontend TUI (`apps/frontend-tui/`)

**Always refer to:** `apps/frontend-tui/docs/`

Terminal-based user interface for the project management system.

### When Working on Frontend Web (`apps/frontend-web/`)

**Always refer to:** `apps/frontend-web/docs/`

Web-based user interface for the project management system.

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

---

## Backend-Specific Instructions

### Adding New Endpoints

1. Create/update **Model** in `internal/models/`
2. Implement **Repository** in `internal/repositories/`
3. Implement **Service** in `internal/services/`
4. Create **Resource** (HTTP handlers) in `internal/resources/`
5. Register routes in `internal/handler.go`

### Database Changes

1. Create new migration files: `migrations/NNNN_name.{up,down}.sql`
2. Write both UP and DOWN migrations
3. Test migrations before committing
4. Never edit existing migrations (create new ones instead)

### Authentication

- Public routes: Register in `RegisterPublicRoutes()`
- Private routes: Register in `RegisterPrivateRoutes()` (auto-protected by SessionMiddleware)
- JWT tokens managed in `repositories/auth_repository.go`

---

## Quick Reference

### Current Tech Stack (Backend)

- **Language:** Go 1.25.5
- **Framework:** Huma v2 (OpenAPI-first)
- **Database:** PostgreSQL with pgx driver
- **Authentication:** JWT (access + refresh tokens)
- **Migrations:** golang-migrate
- **API Docs:** Auto-generated OpenAPI at `/docs`

---

## Context-Aware

**When suggesting code:**

- ✅ Match the existing code style in the current app
- ✅ Reference the appropriate documentation
- ✅ Follow established patterns (especially RSRM for backend)
- ✅ Consider the monorepo structure
- ✅ Suggest small, incremental changes for solo development

**When explaining architecture:**

- ✅ Point to specific docs in `apps/*/docs/`
- ✅ Show examples from existing code
- ✅ Explain trade-offs for solo vs team development
- ✅ Keep explanations practical and actionable

**When debugging:**

- ✅ Check relevant docs first
- ✅ Look at similar working implementations
- ✅ Suggest adding logs for troubleshooting
- ✅ Consider the data flow through layers

---

## Notes

- This is a **solo development project** - optimize for velocity and learning
- Tests, security hardening, and production features can be added incrementally
- Focus on getting features working end-to-end before optimization
- The architecture is solid - follow established patterns for consistency

---

## Key Documentation Entry Points

| Working On      | Read First                                    | Key Concepts               |
| --------------- | --------------------------------------------- | -------------------------- |
| Backend API     | `apps/backend/docs/01-architecture.md`        | Layers, RSRM pattern       |
| Backend Routes  | `apps/backend/docs/05-main-and-handler.md`    | Route registration, DI     |
| Backend Auth    | `apps/backend/docs/04-middlewares.md`         | JWT validation, middleware |
| Database Schema | `apps/backend/docs/06-migrations-workflow.md` | Migration workflow         |
| Request Flow    | `apps/backend/docs/02-code-flow.md`           | Full lifecycle             |
| Design Patterns | `apps/backend/docs/03-rsrm-pattern.md`        | Layer responsibilities     |
| Concurrency     | `apps/backend/docs/07-concurrency.md`         | Goroutine patterns         |

---

**Remember:** The documentation is your friend! Check the relevant `apps/*/docs/` directory before suggesting changes or implementations.
