# Architecture Documentation

## Overview

The Fluxis backend is a RESTful API service built with **Go**, following clean architecture principles with clear separation of concerns. The application uses a layered architecture pattern to ensure maintainability, testability, and scalability.

## Tech Stack

### Core Framework & Libraries

- **Language**: Go 1.25.5
- **HTTP Framework**: Standard library `net/http` with `http.NewServeMux()`
- **API Framework**: [Huma v2](https://github.com/danielgtaylor/huma) - OpenAPI-driven REST API framework
  - Provides automatic OpenAPI documentation generation
  - Request/response validation
  - Error handling middleware
- **Database**: PostgreSQL
- **Database Driver**: [pgx/v5](https://github.com/jackc/pgx) - High-performance PostgreSQL driver
- **Authentication**: JWT tokens using [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
- **Database Migrations**: [golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate)
- **UUID Handling**: [google/uuid](https://github.com/google/uuid)

## Architectural Layers

The application follows a **4-layer architecture**:

```
┌─────────────────────────────────────┐
│         Resources Layer             │
│  (HTTP Handlers & Route Definition) │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│         Services Layer              │
│    (Business Logic & Validation)    │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│       Repositories Layer            │
│  (Data Access & Database Queries)   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│          Models Layer               │
│    (Data Structures & Schemas)      │
└─────────────────────────────────────┘
```

### 1. Resources Layer (`internal/resources/`)

**Purpose**: Handle HTTP requests and responses

- Defines API endpoints and route registration
- Parses HTTP requests
- Calls service layer methods
- Formats responses (success/error)
- Contains Huma operation definitions (OpenAPI metadata)

**Files**:

- `auth_resource.go` - Authentication endpoints (`/auth/login`, `/auth/refresh`)
- `project_resource.go` - Project CRUD endpoints (`/projects`, `/projects/{id}`)

### 2. Services Layer (`internal/services/`)

**Purpose**: Implement business logic and orchestration

- Contains domain business rules
- Performs data validation (e.g., UUID validation)
- Orchestrates calls to multiple repositories if needed
- Returns domain models or errors

**Files**:

- `auth_service.go` - Authentication logic (credential validation, token generation)
- `project_service.go` - Project business logic (validation, CRUD orchestration)

### 3. Repositories Layer (`internal/repositories/`)

**Purpose**: Database access and data persistence

- Executes SQL queries
- Manages database connections via pgx pool
- Transforms database rows into models
- Handles database-specific errors

**Files**:

- `auth_repository.go` - JWT token generation and validation
- `project_repository.go` - Project database operations

### 4. Models Layer (`internal/models/`)

**Purpose**: Data structures and contracts

- Input/Output DTOs (Data Transfer Objects)
- Validation rules (via struct tags)
- OpenAPI documentation (via struct tags)

**Files**:

- `auth_model.go` - Authentication request/response models
- `project_model.go` - Project entity models (CRUD, search, pagination)

## Supporting Components

### Configuration (`internal/configs/`)

Centralized configuration management:

- `environment.go` - Environment variables and application settings
- `database.go` - Database connection pool setup
- `migration.go` - Database migration runner
- `openapi.go` - OpenAPI/Swagger configuration

### Middleware (`internal/middlewares/`)

Cross-cutting concerns applied to HTTP requests:

- `session_middleware.go` - JWT token validation for protected routes

### Common Utilities (`internal/common/`)

Shared utility functions:

- `uuid.go` - UUID validation helper

### Database Migrations (`migrations/`)

SQL migration files for database schema versioning:

- `000001_setup.up.sql` - Initial schema creation
- `000001_setup.down.sql` - Schema rollback

## Directory Structure

```
apps/backend/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── internal/
│   ├── handler.go               # Route registration
│   ├── common/
│   │   └── uuid.go              # Shared utilities
│   ├── configs/
│   │   ├── database.go          # DB connection
│   │   ├── environment.go       # Environment config
│   │   ├── migration.go         # Migration runner
│   │   └── openapi.go           # API documentation config
│   ├── middlewares/
│   │   └── session_middleware.go # JWT auth middleware
│   ├── models/
│   │   ├── auth_model.go        # Auth DTOs
│   │   └── project_model.go     # Project DTOs
│   ├── repositories/
│   │   ├── auth_repository.go   # Auth data access
│   │   └── project_repository.go # Project data access
│   ├── resources/
│   │   ├── auth_resource.go     # Auth HTTP handlers
│   │   └── project_resource.go  # Project HTTP handlers
│   └── services/
│       ├── auth_service.go      # Auth business logic
│       └── project_service.go   # Project business logic
├── migrations/
│   ├── 000001_setup.up.sql      # Schema creation
│   └── 000001_setup.down.sql    # Schema teardown
├── docker-compose.yaml
├── Dockerfile
├── go.mod
└── README.md
```

## Key Architectural Decisions

### 1. **Dependency Injection via Constructor Functions**

Each layer uses constructor functions (e.g., `NewAuthService()`, `NewProjectResource()`) to inject dependencies, promoting loose coupling and testability.

```go
// Example from handler.go
authRepo := repositories.NewAuthRepository(pgx)
authSrv := services.NewAuthService(authRepo)
resources.NewAuthResource(authSrv).Routes(api)
```

### 2. **Separation of Public and Private Routes**

Routes are divided into:

- **Public Routes**: No authentication required (e.g., `/auth/login`)
- **Private Routes**: Protected by `SessionMiddleware` (e.g., `/projects`)

This is defined in `internal/handler.go`:

```go
func RegisterPublicRoutes(api huma.API, pgx *pgxpool.Pool)
func RegisterPrivateRoutes(api huma.API, pgx *pgxpool.Pool)
```

### 3. **Database Connection Pooling**

Uses pgx connection pooling (`pgxpool.Pool`) for efficient database connection management, passed down to all repositories.

### 4. **Migration-First Database Management**

Database schema changes are managed through migration files, ensuring version control and reproducibility across environments.

### 5. **Error Handling via Huma Error Types**

Standardized HTTP error responses using Huma's error helpers:

- `huma.Error400BadRequest()`
- `huma.Error401Unauthorized()`
- `huma.Error404NotFound()`

### 6. **Context Propagation**

Go's `context.Context` is threaded through all layers for request cancellation and timeout handling.

## Data Flow Example

A typical request flow through the architecture:

```
1. Client sends HTTP request
   ↓
2. HTTP Router (net/http.ServeMux)
   ↓
3. Middleware (SessionMiddleware) - validates JWT token
   ↓
4. Resource Layer - parses request, calls service
   ↓
5. Service Layer - validates business rules
   ↓
6. Repository Layer - executes SQL query
   ↓
7. Database - returns data
   ↓
8. Repository → Service → Resource
   ↓
9. HTTP Response sent to client
```

## Scalability Considerations

1. **Stateless Design**: JWT-based authentication allows horizontal scaling
2. **Connection Pooling**: Efficient database connection reuse
3. **Layered Architecture**: Easy to extract services into microservices if needed
4. **Configuration via Environment Variables**: Supports different deployment environments

## Security Features

- JWT-based authentication with access and refresh tokens
- Bearer token authentication in HTTP headers
- Password validation (currently hardcoded admin credentials)
- SQL injection prevention via parameterized queries
- Soft deletes for data retention

## Development & Deployment

- **Containerization**: Dockerfile provided for containerized deployments
- **Local Development**: docker-compose.yaml for local PostgreSQL setup
- **Migration Management**: Automatic migration execution on startup
- **OpenAPI Documentation**: Auto-generated at runtime

## Future Enhancement Opportunities

1. **Password Hashing**: Implement bcrypt for credential storage
2. **Database Secrets Management**: Move JWT secret to environment variables
3. **Observability**: Add structured logging, metrics, and tracing
4. **Testing**: Implement unit and integration tests
5. **API Versioning**: Support multiple API versions
6. **Rate Limiting**: Add rate limiting middleware
7. **CORS Configuration**: Configurable CORS policies
8. **Database Transactions**: Add transaction support for multi-step operations
