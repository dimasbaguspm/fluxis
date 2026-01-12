# Code Flow Documentation

## Overview

This document explains the request lifecycle and code execution flow in the Fluxis backend application, from application startup to request processing and response generation.

## Application Startup Flow

### 1. Entry Point: `cmd/app/main.go`

The application starts in the `main()` function with the following sequence:

```go
func main() {
    // 1. Create context for lifecycle management
    ctx := context.Background()

    // 2. Initialize HTTP router
    r := http.NewServeMux()

    // 3. Load environment variables
    env := configs.NewEnvironment()

    // 4. Setup database configuration
    db := configs.NewDatabase(env)

    // 5. Establish database connection
    pool, err := db.Connect(ctx)

    // 6. Run database migrations
    migration := configs.Migration(env)
    migration.Up()

    // 7. Initialize Huma API with OpenAPI config
    humaApi := humago.New(r, configs.GetOpenapiConfig(env))

    // 8. Register public routes (no auth)
    internal.RegisterPublicRoutes(humaApi, pool)

    // 9. Register private routes (with auth middleware)
    internal.RegisterPrivateRoutes(humaApi, pool)

    // 10. Start HTTP server
    http.ListenAndServe(fmt.Sprintf(":%s", env.AppPort), r)
}
```

### Detailed Startup Steps

#### Step 1-2: Context & Router Initialization

```go
ctx := context.Background()
r := http.NewServeMux()
```

- Creates a background context for managing application lifecycle
- Initializes Go's standard HTTP multiplexer for routing

#### Step 3: Environment Configuration Loading

```go
env := configs.NewEnvironment()
```

**Location**: `internal/configs/environment.go`

Loads configuration from environment variables:

- `APP_PORT` - Server port (default: 3000)
- `IS_DEV_ENV` - Development/production flag
- `DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_PORT` - Database credentials
- `ADMIN_USERNAME`, `ADMIN_PASSWORD` - Admin credentials

Prints configuration summary to console.

#### Step 4-5: Database Connection Setup

```go
db := configs.NewDatabase(env)
pool, err := db.Connect(ctx)
```

**Location**: `internal/configs/database.go`

- Creates PostgreSQL connection string from environment variables
- Establishes connection pool using `pgxpool.New()`
- Pings database to verify connectivity
- Returns `*pgxpool.Pool` for use across the application

#### Step 6: Database Migrations

```go
migration := configs.Migration(env)
migration.Up()
```

**Location**: `internal/configs/migration.go`

- Initializes migration runner with connection string
- Executes pending `.up.sql` migration files from `migrations/` directory
- Skips if no new migrations (`migrate.ErrNoChange`)
- Panics on migration failure

#### Step 7: Huma API Initialization

```go
humaApi := humago.New(r, configs.GetOpenapiConfig(env))
```

**Location**: `internal/configs/openapi.go`

- Creates Huma API instance wrapping the HTTP router
- Configures OpenAPI documentation:
  - Title: "Fluxis"
  - Version: "1.0.0"
  - Server URL based on environment (localhost for dev, `/api` for prod)
  - JWT Bearer authentication scheme

#### Step 8-9: Route Registration

```go
internal.RegisterPublicRoutes(humaApi, pool)
internal.RegisterPrivateRoutes(humaApi, pool)
```

**Location**: `internal/handler.go`

See **Route Registration Flow** section below.

#### Step 10: HTTP Server Start

```go
http.ListenAndServe(fmt.Sprintf(":%s", env.AppPort), r)
```

- Starts HTTP server on configured port
- Blocks until server shutdown

---

## Route Registration Flow

### Public Routes Registration

**Function**: `RegisterPublicRoutes(api huma.API, pgx *pgxpool.Pool)`

```go
func RegisterPublicRoutes(api huma.API, pgx *pgxpool.Pool) {
    // 1. Create repository with database pool
    authRepo := repositories.NewAuthRepository(pgx)

    // 2. Create service with repository dependency
    authSrv := services.NewAuthService(authRepo)

    // 3. Create resource with service dependency
    authResource := resources.NewAuthResource(authSrv)

    // 4. Register routes with Huma API
    authResource.Routes(api)
}
```

**Registered Endpoints**:

- `POST /auth/login` - User authentication
- `POST /auth/refresh` - Token refresh

### Private Routes Registration

**Function**: `RegisterPrivateRoutes(api huma.API, pgx *pgxpool.Pool)`

```go
func RegisterPrivateRoutes(api huma.API, pgx *pgxpool.Pool) {
    // 1. Apply session middleware to all routes
    api.UseMiddleware(middlewares.SessionMiddleware(api))

    // 2. Create repository with database pool
    projectRepo := repositories.NewProjectRepository(pgx)

    // 3. Create service with repository dependency
    projectSrv := services.NewProjectService(projectRepo)

    // 4. Create resource with service dependency
    projectResource := resources.NewProjectResource(projectSrv)

    // 5. Register routes with Huma API
    projectResource.Routes(api)
}
```

**Registered Endpoints** (all require JWT authentication):

- `GET /projects` - List projects (paginated, with search/filter)
- `GET /projects/{projectId}` - Get project details
- `POST /projects/{projectId}` - Create new project
- `PATCH /projects/{projectId}` - Update project
- `DELETE /projects/{projectId}` - Soft delete project

---

## Request Processing Flow

### Example: Login Request Flow

Let's trace a `POST /auth/login` request through the entire system.

#### 1. HTTP Request Arrives

```
POST /auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}
```

#### 2. HTTP Router Matching

- `http.ServeMux` matches route to Huma handler
- Huma validates request against OpenAPI schema
- Parses JSON body into `models.AuthLoginInputModel`

#### 3. Resource Layer: `auth_resource.go`

```go
func (ar AuthResource) login(ctx context.Context,
    input *struct{ Body models.AuthLoginInputModel })
    (*struct{ Body models.AuthLoginOutputModel }, error) {

    // 1. Load environment (for admin credentials)
    env := configs.NewEnvironment()

    // 2. Call service layer
    svcResp, err := ar.authService.Login(input.Body, env)
    if err != nil {
        return nil, err  // Huma handles error serialization
    }

    // 3. Wrap response
    resp := &struct{ Body models.AuthLoginOutputModel }{
        Body: svcResp,
    }

    return resp, nil
}
```

#### 4. Service Layer: `auth_service.go`

```go
func (as *AuthService) Login(data models.AuthLoginInputModel,
    env configs.Environment) (models.AuthLoginOutputModel, error) {

    // 1. Validate credentials against environment config
    isValid := env.Admin.Username != "" &&
               env.Admin.Password != "" &&
               data.Username == env.Admin.Username &&
               data.Password == env.Admin.Password

    if !isValid {
        return models.AuthLoginOutputModel{},
            huma.Error401Unauthorized("Invalid credentials")
    }

    // 2. Generate JWT tokens via repository
    accessToken, refreshToken, err :=
        as.authRepo.GenerateFreshTokens(data)
    if err != nil {
        return models.AuthLoginOutputModel{}, err
    }

    // 3. Return output model
    return models.AuthLoginOutputModel{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        Username:     data.Username,
    }, nil
}
```

#### 5. Repository Layer: `auth_repository.go`

```go
func (ar AuthRepository) GenerateFreshTokens(m models.AuthLoginInputModel)
    (accessToken, refreshToken string, err error) {

    // Generate access token (7 day expiry)
    accessToken, err = generateToken(accessTokenType)
    if err != nil {
        return "", "", err
    }

    // Generate refresh token (30 day expiry)
    refreshToken, err = generateToken(refreshTokenType)
    if err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil
}

func generateToken(sub string) (string, error) {
    now := time.Now()

    // Set expiry based on token type
    var expiredAt time.Time
    switch sub {
    case accessTokenType:
        expiredAt = now.Add(7 * 24 * time.Hour)
    case refreshTokenType:
        expiredAt = now.Add(30 * 24 * time.Hour)
    }

    // Create JWT claims
    claims := jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(expiredAt),
        IssuedAt:  jwt.NewNumericDate(now),
        Subject:   sub,
    }

    // Sign and return token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretJWT))
}
```

#### 6. Response Flow Back

```
Repository → Service → Resource → Huma → HTTP Response
```

#### 7. HTTP Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "username": "admin"
}
```

---

### Example: Protected Project List Request Flow

Let's trace a `GET /projects?status=active&pageNumber=1` request.

#### 1. HTTP Request Arrives

```
GET /projects?status=active&pageNumber=1
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

#### 2. Middleware Execution: `session_middleware.go`

```go
func SessionMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
    return func(ctx huma.Context, next func(huma.Context)) {
        // 1. Extract Authorization header
        cat := ctx.Header("Authorization")

        // 2. Remove "Bearer " prefix
        if len(cat) > 7 && cat[:7] == "Bearer " {
            cat = cat[7:]
        }

        // 3. Validate JWT token
        isValid := authRepo.IsTokenValid(cat)

        if !isValid {
            // Return 400 error if invalid
            huma.WriteErr(api, ctx,
                repositories.AuthErrorInvalidAccessToken.GetStatus(),
                repositories.AuthErrorInvalidAccessToken.Error())
            return
        }

        // 4. Token valid - continue to next handler
        next(ctx)
    }
}
```

**Token Validation in Repository**:

```go
func (ar AuthRepository) IsTokenValid(token string) bool {
    // Parse JWT token
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, AuthErrorInvalidSigningMethod
        }
        return []byte(secretJWT), nil
    })

    if err != nil || !parsedToken.Valid {
        return false
    }

    return true
}
```

#### 3. Resource Layer: `project_resource.go`

```go
func (pr ProjectResource) getPaginated(ctx context.Context,
    input *models.ProjectSearchModel)
    (*struct{ Body models.ProjectPaginatedModel }, error) {

    // 1. Query parameters automatically parsed by Huma into model
    // input = {
    //   Status: ["active"],
    //   PageNumber: 1,
    //   PageSize: 25 (default),
    //   SortBy: "createdAt" (default),
    //   SortOrder: "desc" (default)
    // }

    // 2. Call service layer
    respSrv, err := pr.projectSrv.GetPaginated(ctx, *input)
    if err != nil {
        return nil, err
    }

    // 3. Return wrapped response
    return &struct{ Body models.ProjectPaginatedModel }{
        Body: respSrv,
    }, nil
}
```

#### 4. Service Layer: `project_service.go`

```go
func (ps *ProjectService) GetPaginated(ctx context.Context,
    q models.ProjectSearchModel) (models.ProjectPaginatedModel, error) {

    // 1. Validate UUIDs in ID filter (if provided)
    for _, id := range q.ID {
        if !common.ValidateUUID(id) {
            return models.ProjectPaginatedModel{},
                huma.Error400BadRequest("Must provide UUID format")
        }
    }

    // 2. Delegate to repository
    return ps.projectRepo.GetPaginated(ctx, q)
}
```

#### 5. Repository Layer: `project_repository.go`

```go
func (pr ProjectRepository) GetPaginated(ctx context.Context,
    query models.ProjectSearchModel) (models.ProjectPaginatedModel, error) {

    // 1. Map query parameters to SQL columns
    sortColumn := "created_at"  // from query.SortBy
    sortOrder := "DESC"          // from query.SortOrder
    offset := (query.PageNumber - 1) * query.PageSize
    searchPattern := "%" + query.Query + "%"

    // 2. Execute complex SQL query with CTE
    sql := `
        WITH filtered AS (
            SELECT id, name, description, status, created_at, updated_at
            FROM projects
            WHERE deleted_at IS NULL
                AND ($1::uuid[] IS NULL OR id = ANY($1))
                AND ($2::text[] IS NULL OR status::text = ANY($2))
                AND ($3 = '' OR name ILIKE $3 OR description ILIKE $3)
        ),
        counted AS (
            SELECT COUNT(*) as total FROM filtered
        )
        SELECT f.*, c.total
        FROM filtered f
        CROSS JOIN counted c
        ORDER BY f.` + sortColumn + ` ` + sortOrder + `
        LIMIT $4 OFFSET $5
    `

    // 3. Execute query
    rows, err := pr.pgx.Query(ctx, sql,
        query.ID, query.Status, searchPattern, query.PageSize, offset)

    // 4. Scan rows into models
    var items []models.ProjectModel
    var totalCount int

    for rows.Next() {
        var item models.ProjectModel
        err := rows.Scan(&item.ID, &item.Name, &item.Description,
            &item.Status, &item.CreatedAt, &item.UpdatedAt, &totalCount)
        items = append(items, item)
    }

    // 5. Calculate pagination metadata
    totalPages := (totalCount + query.PageSize - 1) / query.PageSize

    // 6. Return paginated result
    return models.ProjectPaginatedModel{
        Items:      items,
        PageNumber: query.PageNumber,
        PageSize:   query.PageSize,
        TotalPages: totalPages,
        TotalCount: totalCount,
    }, nil
}
```

#### 6. Response Flow Back

```
Repository → Service → Resource → Huma → HTTP Response
```

#### 7. HTTP Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "items": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Project Alpha",
      "description": "First project",
      "status": "active",
      "createdAt": "2026-01-10T10:30:00Z",
      "updatedAt": "2026-01-10T10:30:00Z"
    }
  ],
  "pageNumber": 1,
  "pageSize": 25,
  "totalPages": 1,
  "totalCount": 1
}
```

---

## Error Handling Flow

### Example: Invalid UUID Error

```go
// Request
GET /projects/invalid-id

// Flow
Resource → Service (UUID validation fails) → Error returned

// Response
HTTP/1.1 400 Bad Request
{
  "status": 400,
  "title": "Bad Request",
  "detail": "Must provide UUID format"
}
```

### Example: Unauthorized Access

```go
// Request
GET /projects
Authorization: Bearer invalid_token

// Flow
Middleware (token validation fails) → Error returned immediately

// Response
HTTP/1.1 400 Bad Request
{
  "status": 400,
  "title": "Bad Request",
  "detail": "Invalid access token"
}
```

---

## Key Observations

1. **Context Propagation**: `context.Context` is passed through all layers for cancellation/timeout
2. **Error Short-Circuit**: Errors return immediately, unwinding the stack
3. **Dependency Injection**: Each layer receives its dependencies via constructors
4. **Stateless Handlers**: No shared state between requests (except DB pool)
5. **Automatic Validation**: Huma validates requests against OpenAPI schema
6. **Structured Responses**: Consistent error format via Huma error helpers
