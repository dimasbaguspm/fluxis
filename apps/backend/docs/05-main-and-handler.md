# Main Command and Handler Documentation

## Overview

This document explains the application entry point and route registration system in the Fluxis backend. The main command initializes all components and starts the HTTP server, while the handler orchestrates route registration and dependency wiring.

---

## Application Entry Point

**Location**: `cmd/app/main.go`

### Purpose

The `main.go` file serves as the application's entry point, responsible for:

- Initializing the application context
- Loading configuration from environment variables
- Establishing database connections
- Running database migrations
- Setting up the HTTP server and API framework
- Registering all routes (public and private)
- Starting the HTTP server

### Complete Implementation

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"

    "github.com/danielgtaylor/huma/v2/adapters/humago"
    "github.com/dimasbaguspm/fluxis/internal"
    "github.com/dimasbaguspm/fluxis/internal/configs"
)

func main() {
    // 1. Create application context
    ctx := context.Background()

    // 2. Initialize HTTP router
    r := http.NewServeMux()

    // 3. Load environment configuration
    env := configs.NewEnvironment()

    // 4. Setup database configuration
    db := configs.NewDatabase(env)

    // 5. Establish database connection
    pool, err := db.Connect(ctx)
    if err != nil {
        slog.Error("Database is unreachable", "err", err)
        panic(err)
    }

    // 6. Initialize migration system
    migration := configs.Migration(env)

    // 7. Run database migrations
    slog.Info("Performing migration")
    if err := migration.Up(); err != nil {
        slog.Error("Failed to run migration", "error", err.Error())
        panic(err)
    }
    slog.Info("DB migration completed")

    // 8. Create Huma API instance
    humaApi := humago.New(r, configs.GetOpenapiConfig(env))

    // 9. Register public routes (no authentication)
    internal.RegisterPublicRoutes(humaApi, pool)

    // 10. Register private routes (with authentication)
    internal.RegisterPrivateRoutes(humaApi, pool)

    // 11. Start HTTP server
    slog.Info("All is ready! serving to port ", "Info", env.AppPort)
    http.ListenAndServe(fmt.Sprintf(":%s", env.AppPort), r)
}
```

### Step-by-Step Breakdown

#### Step 1: Application Context

```go
ctx := context.Background()
```

**Purpose**: Create a root context for managing application lifecycle

**Usage**:

- Passed to database connection (for connection timeout/cancellation)
- Can be used for graceful shutdown signals
- Propagated through request handlers

#### Step 2: HTTP Router Initialization

```go
r := http.NewServeMux()
```

**Purpose**: Create Go's standard HTTP request multiplexer

**Characteristics**:

- Built-in Go HTTP router
- Pattern matching for routes
- Wrapped by Huma for enhanced functionality

#### Step 3: Environment Configuration

```go
env := configs.NewEnvironment()
```

**Purpose**: Load and validate environment variables

**Loaded Variables**:

- `APP_PORT` - Server port (default: 3000)
- `IS_DEV_ENV` - Development mode flag
- `DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_PORT` - Database credentials
- `ADMIN_USERNAME`, `ADMIN_PASSWORD` - Admin authentication credentials

**Console Output**:

```
===============
Server:
Port: 3000
Stage: AppStageDev
Credential:
Username: admin
Password: password
===============
```

#### Step 4-5: Database Connection

```go
db := configs.NewDatabase(env)
pool, err := db.Connect(ctx)
```

**Purpose**: Establish PostgreSQL connection pool

**Process**:

1. Construct PostgreSQL connection string from environment variables
2. Create connection pool using `pgxpool.New()`
3. Ping database to verify connectivity
4. Return pool for use throughout application

**Error Handling**:

- Logs error and panics if database is unreachable
- Ensures application doesn't start with invalid database configuration

#### Step 6-7: Database Migrations

```go
migration := configs.Migration(env)
if err := migration.Up(); err != nil {
    slog.Error("Failed to run migration", "error", err.Error())
    panic(err)
}
```

**Purpose**: Apply pending database schema changes

**Process**:

1. Initialize migration runner with database URL
2. Execute all `.up.sql` files in `migrations/` directory
3. Skip if no new migrations (`migrate.ErrNoChange`)
4. Panic if migration fails

**Benefits**:

- Ensures database schema is up-to-date on every startup
- Supports version control of database changes
- Provides reproducible database state

#### Step 8: Huma API Initialization

```go
humaApi := humago.New(r, configs.GetOpenapiConfig(env))
```

**Purpose**: Create Huma API wrapper around HTTP router

**Configuration** (from `configs.GetOpenapiConfig()`):

```go
func GetOpenapiConfig(env Environment) huma.Config {
    config := huma.DefaultConfig("Fluxis", "1.0.0")

    // Set server URL based on environment
    url := "http://localhost:" + env.AppPort
    desc := "Development server"
    if env.AppStage == AppStageProd {
        url = "/api"
        desc = "Proxied server"
    }
    config.Servers = []*huma.Server{{URL: url, Description: desc}}

    // Configure JWT Bearer authentication
    config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
        "bearer": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
    }

    return config
}
```

**Features**:

- Auto-generates OpenAPI documentation
- Validates requests against schemas
- Provides consistent error responses

#### Step 9-10: Route Registration

```go
internal.RegisterPublicRoutes(humaApi, pool)
internal.RegisterPrivateRoutes(humaApi, pool)
```

See [Route Registration System](#route-registration-system) section below.

#### Step 11: HTTP Server Start

```go
http.ListenAndServe(fmt.Sprintf(":%s", env.AppPort), r)
```

**Purpose**: Start HTTP server and listen for requests

**Behavior**:

- Blocks until server shutdown
- Listens on configured port (default: 3000)
- Routes requests through the multiplexer

**Access**:

- Development: `http://localhost:3000`
- OpenAPI docs: `http://localhost:3000/docs` (auto-generated by Huma)

---

## Route Registration System

**Location**: `internal/handler.go`

### Purpose

The handler file provides two functions for registering routes:

1. **Public Routes**: Accessible without authentication
2. **Private Routes**: Protected by JWT authentication middleware

### Complete Implementation

```go
package internal

import (
    "github.com/danielgtaylor/huma/v2"
    "github.com/dimasbaguspm/fluxis/internal/middlewares"
    "github.com/dimasbaguspm/fluxis/internal/repositories"
    "github.com/dimasbaguspm/fluxis/internal/resources"
    "github.com/dimasbaguspm/fluxis/internal/services"
    "github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(api huma.API, pgx *pgxpool.Pool) {
    authRepo := repositories.NewAuthRepository(pgx)
    authSrv := services.NewAuthService(authRepo)

    resources.NewAuthResource(authSrv).Routes(api)
}

func RegisterPrivateRoutes(api huma.API, pgx *pgxpool.Pool) {
    api.UseMiddleware(middlewares.SessionMiddleware(api))

    projectRepo := repositories.NewProjectRepository(pgx)
    projectSrv := services.NewProjectService(projectRepo)

    resources.NewProjectResource(projectSrv).Routes(api)
}
```

---

### Public Routes Registration

```go
func RegisterPublicRoutes(api huma.API, pgx *pgxpool.Pool) {
    // 1. Create repository with database pool
    authRepo := repositories.NewAuthRepository(pgx)

    // 2. Create service with repository dependency
    authSrv := services.NewAuthService(authRepo)

    // 3. Create resource and register routes
    resources.NewAuthResource(authSrv).Routes(api)
}
```

#### Dependency Chain

```
Database Pool (pgxpool.Pool)
    ↓
AuthRepository (data access)
    ↓
AuthService (business logic)
    ↓
AuthResource (HTTP handlers)
    ↓
Registered Routes
```

#### Registered Endpoints

| Method | Path          | Handler              | Description          |
| ------ | ------------- | -------------------- | -------------------- |
| POST   | /auth/login   | authResource.login   | Authenticate user    |
| POST   | /auth/refresh | authResource.refresh | Refresh access token |

**Characteristics**:

- No authentication required
- Accessible to all clients
- Used for obtaining JWT tokens

#### Example Request Flow

```
1. Client → POST /auth/login
2. HTTP Router → authResource.login()
3. authResource → authService.Login()
4. authService → authRepo.GenerateFreshTokens()
5. authRepo → JWT token generation
6. Response ← { accessToken, refreshToken, username }
```

---

### Private Routes Registration

```go
func RegisterPrivateRoutes(api huma.API, pgx *pgxpool.Pool) {
    // 1. Apply authentication middleware to all subsequent routes
    api.UseMiddleware(middlewares.SessionMiddleware(api))

    // 2. Create repository with database pool
    projectRepo := repositories.NewProjectRepository(pgx)

    // 3. Create service with repository dependency
    projectSrv := services.NewProjectService(projectRepo)

    // 4. Create resource and register routes
    resources.NewProjectResource(projectSrv).Routes(api)
}
```

#### Dependency Chain

```
Database Pool (pgxpool.Pool)
    ↓
ProjectRepository (data access)
    ↓
ProjectService (business logic)
    ↓
ProjectResource (HTTP handlers)
    ↓
SessionMiddleware (JWT validation) → Registered Routes
```

#### Registered Endpoints

| Method | Path                  | Handler                      | Description               |
| ------ | --------------------- | ---------------------------- | ------------------------- |
| GET    | /projects             | projectResource.getPaginated | List projects (paginated) |
| GET    | /projects/{projectId} | projectResource.getDetail    | Get single project        |
| POST   | /projects/{projectId} | projectResource.create       | Create new project        |
| PATCH  | /projects/{projectId} | projectResource.update       | Update project            |
| DELETE | /projects/{projectId} | projectResource.delete       | Delete project (soft)     |

**Characteristics**:

- **Authentication Required**: All routes protected by `SessionMiddleware`
- **JWT Token**: Must include `Authorization: Bearer <token>` header
- **Validated Access**: Token validated before reaching handler

#### Example Request Flow (Protected)

```
1. Client → GET /projects (with Bearer token)
2. HTTP Router → SessionMiddleware
3. SessionMiddleware → Validate JWT token
   ├─ Invalid → Return 400 error (stop processing)
   └─ Valid → Continue to handler
4. projectResource.getPaginated()
5. projectService.GetPaginated()
6. projectRepo.GetPaginated()
7. Database query execution
8. Response ← { items: [...], pageNumber: 1, ... }
```

---

## Dependency Injection Pattern

### Constructor Functions

Each layer uses constructor functions for dependency injection:

```go
// Repository layer
func NewAuthRepository(pgx *pgxpool.Pool) AuthRepository {
    return AuthRepository{pgx}
}

func NewProjectRepository(pgx *pgxpool.Pool) ProjectRepository {
    return ProjectRepository{pgx}
}

// Service layer
func NewAuthService(authRepo AuthRepository) AuthService {
    return AuthService{authRepo}
}

func NewProjectService(projectRepo ProjectRepository) ProjectService {
    return ProjectService{projectRepo}
}

// Resource layer
func NewAuthResource(authSrv AuthService) AuthResource {
    return AuthResource{authSrv}
}

func NewProjectResource(projectSrv ProjectService) ProjectResource {
    return ProjectResource{projectSrv}
}
```

### Benefits

1. **Loose Coupling**: Components depend on interfaces, not implementations
2. **Testability**: Easy to mock dependencies for unit testing
3. **Flexibility**: Can swap implementations without changing consumers
4. **Clear Dependencies**: Explicit declaration of what each component needs

### Example: Testing with Mocks

```go
// Production: Real repository
authRepo := repositories.NewAuthRepository(dbPool)
authService := services.NewAuthService(authRepo)

// Testing: Mock repository
mockRepo := &MockAuthRepository{
    GenerateFreshTokensFunc: func(...) (string, string, error) {
        return "mock_access", "mock_refresh", nil
    },
}
authService := services.NewAuthService(mockRepo)
```

---

## Route Organization

### Current Structure

```
Public Routes (No Auth):
├── POST /auth/login    (Login)
└── POST /auth/refresh  (Token Refresh)

Private Routes (JWT Required):
├── GET    /projects              (List projects)
├── GET    /projects/{projectId}  (Get project)
├── POST   /projects/{projectId}  (Create project)
├── PATCH  /projects/{projectId}  (Update project)
└── DELETE /projects/{projectId}  (Delete project)
```

### Adding New Entities

To add a new entity (e.g., "tasks"), follow this pattern:

#### 1. Create Models

```go
// internal/models/task_model.go
type TaskModel struct { ... }
type TaskCreateModel struct { ... }
// etc.
```

#### 2. Create Repository

```go
// internal/repositories/task_repository.go
type TaskRepository struct { pgx *pgxpool.Pool }
func NewTaskRepository(pgx *pgxpool.Pool) TaskRepository { ... }
```

#### 3. Create Service

```go
// internal/services/task_service.go
type TaskService struct { taskRepo TaskRepository }
func NewTaskService(taskRepo TaskRepository) TaskService { ... }
```

#### 4. Create Resource

```go
// internal/resources/task_resource.go
type TaskResource struct { taskSrv TaskService }
func NewTaskResource(taskSrv TaskService) TaskResource { ... }
func (tr TaskResource) Routes(api huma.API) { ... }
```

#### 5. Register Routes in Handler

```go
// internal/handler.go
func RegisterPrivateRoutes(api huma.API, pgx *pgxpool.Pool) {
    api.UseMiddleware(middlewares.SessionMiddleware(api))

    // Existing project routes
    projectRepo := repositories.NewProjectRepository(pgx)
    projectSrv := services.NewProjectService(projectRepo)
    resources.NewProjectResource(projectSrv).Routes(api)

    // NEW: Task routes
    taskRepo := repositories.NewTaskRepository(pgx)
    taskSrv := services.NewTaskService(taskRepo)
    resources.NewTaskResource(taskSrv).Routes(api)
}
```

---

## Error Handling on Startup

### Database Connection Failure

```go
pool, err := db.Connect(ctx)
if err != nil {
    slog.Error("Database is unreachable", "err", err)
    panic(err)  // Application stops immediately
}
```

**Reason**: No point starting server without database connectivity.

### Migration Failure

```go
if err := migration.Up(); err != nil {
    slog.Error("Failed to run migration", "error", err.Error())
    panic(err)  // Application stops immediately
}
```

**Reason**: Running with incorrect schema can cause runtime errors.

### Why Panic?

During startup, encountering critical errors (database unreachable, migration failures) should prevent the application from starting. Panicking ensures:

- Clear failure indication
- No partial/broken application state
- Forces operator to fix configuration issues

---

## Environment-Specific Behavior

### Development Mode (`IS_DEV_ENV=true`)

```go
if env.AppStage == AppStageDev {
    // OpenAPI server URL
    url = "http://localhost:" + env.AppPort
    desc = "Development server"
}
```

**Characteristics**:

- Server URL: `http://localhost:3000`
- OpenAPI docs accessible: `http://localhost:3000/docs`
- Verbose logging
- Hot-reload friendly (if using tools like `air`)

### Production Mode (`IS_DEV_ENV=false` or unset)

```go
if env.AppStage == AppStageProd {
    // OpenAPI server URL
    url = "/api"
    desc = "Proxied server"
}
```

**Characteristics**:

- Server URL: `/api` (expects reverse proxy)
- Assumes API behind proxy (e.g., Nginx, Traefik)
- Reduced logging
- Production-optimized settings

---

## OpenAPI Documentation

### Auto-Generated Documentation

Huma automatically generates OpenAPI documentation accessible at:

```
GET http://localhost:3000/docs
```

**Includes**:

- All registered routes
- Request/response schemas (from model struct tags)
- Authentication requirements
- Enum values and validation rules
- Try-it-out functionality

### Example OpenAPI Output

```yaml
openapi: 3.0.0
info:
  title: Fluxis
  version: 1.0.0
servers:
  - url: http://localhost:3000
    description: Development server
paths:
  /auth/login:
    post:
      operationId: login
      summary: Login
      tags:
        - Authentication
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                username:
                  type: string
                  minLength: 1
                password:
                  type: string
                  minLength: 1
  /projects:
    get:
      operationId: project-get-paginated
      summary: Get Projects
      tags:
        - Project
      security:
        - bearer: []
```

---

## Summary

### Main Command Responsibilities

1. ✅ Initialize application context
2. ✅ Load environment configuration
3. ✅ Establish database connection
4. ✅ Run database migrations
5. ✅ Setup HTTP server and API framework
6. ✅ Register all routes (public and private)
7. ✅ Start HTTP server

### Handler Responsibilities

1. ✅ Separate public and private routes
2. ✅ Apply authentication middleware to protected routes
3. ✅ Wire dependencies (Repository → Service → Resource)
4. ✅ Register routes with Huma API

### Key Design Principles

- **Fail Fast**: Panic on critical startup errors
- **Dependency Injection**: Explicit dependency wiring
- **Separation of Concerns**: Clear boundaries between route types
- **Configuration-Driven**: Environment variables control behavior
- **Documentation-First**: Auto-generated OpenAPI specs

### Adding New Features

To add new endpoints:

1. Create models, repository, service, resource
2. Add registration in `RegisterPublicRoutes()` or `RegisterPrivateRoutes()`
3. Restart application
4. OpenAPI documentation automatically updated
