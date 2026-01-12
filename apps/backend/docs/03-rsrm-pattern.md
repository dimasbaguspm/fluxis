# Resource-Service-Repository-Model Pattern

## Overview

The Fluxis backend implements the **Resource-Service-Repository-Model (RSRM)** pattern, a clean architecture approach that separates concerns across four distinct layers. This pattern provides clear boundaries between different responsibilities, making the codebase maintainable, testable, and scalable.

## The Four Layers

```
┌─────────────────────────────────────────────────────────┐
│                    Resource Layer                       │
│  • HTTP endpoint handlers                               │
│  • Request/response parsing and formatting              │
│  • Route definitions and OpenAPI metadata               │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ DTO (Input Models)
                     ↓
┌─────────────────────────────────────────────────────────┐
│                    Service Layer                        │
│  • Business logic and rules                             │
│  • Data validation                                      │
│  • Orchestration across repositories                    │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ Domain Models
                     ↓
┌─────────────────────────────────────────────────────────┐
│                  Repository Layer                       │
│  • Data access and persistence                          │
│  • SQL query execution                                  │
│  • Database connection management                       │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ Database Entities
                     ↓
┌─────────────────────────────────────────────────────────┐
│                     Model Layer                         │
│  • Data structures (DTOs, entities)                     │
│  • Validation rules (struct tags)                       │
│  • OpenAPI documentation (struct tags)                  │
└─────────────────────────────────────────────────────────┘
```

---

## 1. Model Layer

**Location**: `internal/models/`

### Purpose

Define data structures used across the application, including:

- **Input Models**: Request payloads from clients
- **Output Models**: Response structures to clients
- **Entity Models**: Database record representations
- **Search/Filter Models**: Query parameters for list operations

### Characteristics

- **Struct Tags for Validation**: Uses Huma's validation tags (`minLength`, `enum`, `required`)
- **OpenAPI Documentation**: `doc` tags generate API documentation
- **Immutable Data Structures**: Simple Go structs with no methods
- **Type Safety**: Strong typing for all fields

### Example: Auth Models (`auth_model.go`)

```go
type AuthLoginInputModel struct {
    Username string `json:"username" minLength:"1" doc:"Your username"`
    Password string `json:"password" minLength:"1" doc:"Your password"`
}

type AuthLoginOutputModel struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
    Username     string `json:"username"`
}

type AuthRefreshInputModel struct {
    RefreshToken string `json:"refreshToken"`
}

type AuthRefreshOutputModel struct {
    AccessToken string `json:"acessToken"`
}
```

**Key Points**:

- Input models have validation constraints (`minLength:"1"`)
- Output models have no validation (server-controlled)
- Clear separation between input and output structures

### Example: Project Models (`project_model.go`)

```go
// Entity model (database record)
type ProjectModel struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    Description string     `json:"description"`
    Status      string     `json:"status" enum:"active,paused,archived"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
    DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

// Pagination wrapper
type ProjectPaginatedModel struct {
    Items      []ProjectModel `json:"items"`
    PageNumber int            `json:"pageNumber"`
    PageSize   int            `json:"pageSize"`
    TotalPages int            `json:"totalPages"`
    TotalCount int            `json:"totalCount"`
}

// Search/filter model (query parameters)
type ProjectSearchModel struct {
    ID         []string `query:"id"`
    Query      string   `query:"query"`
    Status     []string `query:"status" enum:"active,paused,archived"`
    PageNumber int      `query:"pageNumber" default:"1"`
    PageSize   int      `query:"pageSize" default:"25"`
    SortBy     string   `query:"sortBy" enum:"createdAt,updatedAt,status" default:"createdAt"`
    SortOrder  string   `query:"sortOrder" enum:"asc,desc" default:"desc"`
}

// Create model (POST request)
type ProjectCreateModel struct {
    Name        string `json:"name" minLength:"1"`
    Description string `json:"description" minLength:"1"`
    Status      string `json:"status" enum:"active,paused,archived"`
}

// Update model (PATCH request - all fields optional)
type ProjectUpdateModel struct {
    Name        string `json:"name,omitempty" required:"false" minLength:"1"`
    Description string `json:"description,omitempty" required:"false" minLenght:"1"`
    Status      string `json:"status,omitempty" required:"false" enum:"active,paused,archived"`
}
```

**Key Points**:

- `ProjectModel`: Full entity with all fields
- `ProjectSearchModel`: Query parameters with defaults
- `ProjectCreateModel`: Required fields only
- `ProjectUpdateModel`: All fields optional (`omitempty`, `required:"false"`)
- `ProjectPaginatedModel`: Wraps list with metadata

---

## 2. Repository Layer

**Location**: `internal/repositories/`

### Purpose

Handle all data persistence and retrieval operations:

- Execute SQL queries
- Manage database connections
- Transform database rows into models
- Handle database-specific errors

### Characteristics

- **Database Access Only**: No business logic
- **SQL Expert**: Knows how to query and manipulate data
- **Error Handling**: Converts database errors to application errors
- **Connection Pool**: Uses `*pgxpool.Pool` for connections

### Structure

```go
type <Entity>Repository struct {
    pgx *pgxpool.Pool
}

func New<Entity>Repository(pgx *pgxpool.Pool) <Entity>Repository {
    return <Entity>Repository{pgx}
}

// CRUD methods
func (repo <Entity>Repository) GetPaginated(ctx context.Context, query SearchModel) (PaginatedModel, error)
func (repo <Entity>Repository) GetDetail(ctx context.Context, id string) (EntityModel, error)
func (repo <Entity>Repository) Create(ctx context.Context, payload CreateModel) (EntityModel, error)
func (repo <Entity>Repository) Update(ctx context.Context, id string, payload UpdateModel) (EntityModel, error)
func (repo <Entity>Repository) Delete(ctx context.Context, id string) error
```

### Example: Project Repository (`project_repository.go`)

#### Constructor

```go
type ProjectRepository struct {
    pgx *pgxpool.Pool
}

func NewProjectRepository(pgx *pgxpool.Pool) ProjectRepository {
    return ProjectRepository{pgx}
}
```

#### Get Paginated (Complex Query)

```go
func (pr ProjectRepository) GetPaginated(ctx context.Context,
    query models.ProjectSearchModel) (models.ProjectPaginatedModel, error) {

    // Map enum values to SQL columns
    sortByMap := map[string]string{
        "createdAt": "created_at",
        "updatedAt": "updated_at",
        "status":    "status",
    }
    sortColumn, _ := sortByMap[query.SortBy]
    sortOrder, _ := map[string]string{"asc": "ASC", "desc": "DESC"}[query.SortOrder]

    // Calculate offset for pagination
    offset := (query.PageNumber - 1) * query.PageSize
    searchPattern := "%" + query.Query + "%"

    // SQL with CTE for filtering and counting
    sql := `
        WITH filtered AS (
            SELECT id, name, description, status, created_at, updated_at
            FROM projects
            WHERE deleted_at IS NULL
                AND ($1::uuid[] IS NULL OR CARDINALITY($1::uuid[]) = 0 OR id = ANY($1))
                AND ($2::text[] IS NULL OR CARDINALITY($2::text[]) = 0 OR status::text = ANY($2))
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

    // Execute query
    rows, err := pr.pgx.Query(ctx, sql,
        query.ID, query.Status, searchPattern, query.PageSize, offset)
    if err != nil {
        return models.ProjectPaginatedModel{},
            huma.Error400BadRequest("Unable to query projects", err)
    }
    defer rows.Close()

    // Scan results
    var items []models.ProjectModel
    var totalCount int

    for rows.Next() {
        var item models.ProjectModel
        err := rows.Scan(&item.ID, &item.Name, &item.Description,
            &item.Status, &item.CreatedAt, &item.UpdatedAt, &totalCount)
        if err != nil {
            return models.ProjectPaginatedModel{},
                huma.Error400BadRequest("Unable to scan project data", err)
        }
        items = append(items, item)
    }

    // Calculate pagination metadata
    totalPages := 0
    if totalCount > 0 {
        totalPages = (totalCount + query.PageSize - 1) / query.PageSize
    }

    // Return paginated result
    return models.ProjectPaginatedModel{
        Items:      items,
        PageNumber: query.PageNumber,
        PageSize:   query.PageSize,
        TotalPages: totalPages,
        TotalCount: totalCount,
    }, nil
}
```

#### Get Detail (Single Record)

```go
func (pr ProjectRepository) GetDetail(ctx context.Context, id string)
    (models.ProjectModel, error) {

    var data models.ProjectModel

    sql := `SELECT id, name, description, status, created_at, updated_at
            FROM projects
            WHERE id = $1::uuid AND deleted_at IS NULL`

    err := pr.pgx.QueryRow(ctx, sql, id).Scan(
        &data.ID, &data.Name, &data.Description,
        &data.Status, &data.CreatedAt, &data.UpdatedAt)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return models.ProjectModel{}, huma.Error404NotFound("No project found")
        }
        return models.ProjectModel{},
            huma.Error400BadRequest("Unable to query the project details", err)
    }

    return data, nil
}
```

#### Create (Insert)

```go
func (pr ProjectRepository) Create(ctx context.Context,
    payload models.ProjectCreateModel) (models.ProjectModel, error) {

    var data models.ProjectModel

    sql := `INSERT into projects (name, description, status)
            VALUES ($1, $2, $3)
            RETURNING id, name, description, status, created_at, updated_at`

    err := pr.pgx.QueryRow(ctx, sql,
        payload.Name, payload.Description, payload.Status).Scan(
        &data.ID, &data.Name, &data.Description,
        &data.Status, &data.CreatedAt, &data.UpdatedAt)

    if err != nil {
        return models.ProjectModel{},
            huma.Error400BadRequest("Unable to create project", err)
    }

    return data, nil
}
```

#### Update (Partial Update with COALESCE)

```go
func (pr ProjectRepository) Update(ctx context.Context, id string,
    payload models.ProjectUpdateModel) (models.ProjectModel, error) {

    var data models.ProjectModel

    sql := `UPDATE projects
            SET name = COALESCE($1, name),
                description = COALESCE($2, description),
                status = COALESCE($3, status),
                updated_at = CURRENT_TIMESTAMP
            WHERE id = $4::uuid AND deleted_at IS NULL
            RETURNING id, name, description, status, created_at, updated_at`

    err := pr.pgx.QueryRow(ctx, sql,
        payload.Name, payload.Description, payload.Status, id).Scan(
        &data.ID, &data.Name, &data.Description,
        &data.Status, &data.CreatedAt, &data.UpdatedAt)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return models.ProjectModel{}, huma.Error404NotFound("No project found")
        }
        return models.ProjectModel{},
            huma.Error400BadRequest("Unable to update the record", err)
    }

    return data, nil
}
```

**Note**: Uses `COALESCE()` to only update non-null fields (partial updates).

#### Delete (Soft Delete)

```go
func (pr ProjectRepository) Delete(ctx context.Context, id string) error {
    sql := `UPDATE projects
            SET deleted_at = CURRENT_TIMESTAMP
            WHERE id = $1::uuid AND deleted_at IS NULL`

    cmdTag, err := pr.pgx.Exec(ctx, sql, id)
    if err != nil {
        return huma.Error400BadRequest("Unable to delete the record", err)
    }
    if cmdTag.RowsAffected() == 0 {
        return huma.Error404NotFound("No project found")
    }

    return nil
}
```

**Note**: Soft delete - sets `deleted_at` timestamp instead of removing row.

### Example: Auth Repository (`auth_repository.go`)

The auth repository handles JWT token operations (non-database):

```go
type AuthRepository struct {
    pgx *pgxpool.Pool
}

func NewAuthRepository(pgx *pgxpool.Pool) AuthRepository {
    return AuthRepository{pgx}
}

const secretJWT = "some-random-things-that-soon-will-be-replaced"

// Generate both access and refresh tokens
func (ar AuthRepository) GenerateFreshTokens(m models.AuthLoginInputModel)
    (accessToken, refreshToken string, err error) {

    accessToken, err = generateToken(accessTokenType)
    if err != nil {
        return "", "", err
    }

    refreshToken, err = generateToken(refreshTokenType)
    if err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil
}

// Validate refresh token and generate new access token
func (ar AuthRepository) RegenerateAccessToken(refreshToken string) (string, error) {
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, AuthErrorInvalidSigningMethod
        }
        return []byte(secretJWT), nil
    })

    if err != nil || !token.Valid {
        return "", AuthErrorInvalidRefreshToken
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || claims["sub"] != refreshTokenType {
        return "", AuthErrorInvalidAccessToken
    }

    return generateToken(accessTokenType)
}

// Validate access token
func (ar AuthRepository) IsTokenValid(token string) bool {
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, AuthErrorInvalidSigningMethod
        }
        return []byte(secretJWT), nil
    })

    return err == nil && parsedToken.Valid
}
```

---

## 3. Service Layer

**Location**: `internal/services/`

### Purpose

Implement business logic and orchestration:

- Validate business rules
- Coordinate multiple repository calls
- Transform data between layers
- Make business decisions

### Characteristics

- **Business Logic Owner**: Contains domain-specific rules
- **No Database Knowledge**: Delegates persistence to repositories
- **Stateless**: No instance state between calls
- **Error Propagation**: Returns errors to resource layer

### Structure

```go
type <Entity>Service struct {
    <entity>Repo repositories.<Entity>Repository
}

func New<Entity>Service(<entity>Repo repositories.<Entity>Repository) <Entity>Service {
    return <Entity>Service{<entity>Repo}
}

// Business methods matching resource operations
func (svc *<Entity>Service) GetPaginated(ctx context.Context, q SearchModel) (PaginatedModel, error)
func (svc *<Entity>Service) GetDetail(ctx context.Context, id string) (EntityModel, error)
func (svc *<Entity>Service) Create(ctx context.Context, p CreateModel) (EntityModel, error)
func (svc *<Entity>Service) Update(ctx context.Context, id string, p UpdateModel) (EntityModel, error)
func (svc *<Entity>Service) Delete(ctx context.Context, id string) error
```

### Example: Project Service (`project_service.go`)

```go
type ProjectService struct {
    projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
    return ProjectService{projectRepo}
}

// Get paginated list with ID validation
func (ps *ProjectService) GetPaginated(ctx context.Context,
    q models.ProjectSearchModel) (models.ProjectPaginatedModel, error) {

    // Business rule: Validate all UUIDs in filter
    for _, id := range q.ID {
        if !common.ValidateUUID(id) {
            return models.ProjectPaginatedModel{},
                huma.Error400BadRequest("Must provide UUID format")
        }
    }

    // Delegate to repository
    return ps.projectRepo.GetPaginated(ctx, q)
}

// Get single project with ID validation
func (ps *ProjectService) GetDetail(ctx context.Context, id string)
    (models.ProjectModel, error) {

    // Business rule: Validate UUID format
    isValidID := common.ValidateUUID(id)
    if !isValidID {
        return models.ProjectModel{},
            huma.Error400BadRequest("Must provide UUID format")
    }

    // Delegate to repository
    return ps.projectRepo.GetDetail(ctx, id)
}

// Create new project (no additional validation needed)
func (ps *ProjectService) Create(ctx context.Context,
    p models.ProjectCreateModel) (models.ProjectModel, error) {

    // Could add business rules here (e.g., name uniqueness check)
    return ps.projectRepo.Create(ctx, p)
}

// Update project with ID validation
func (ps *ProjectService) Update(ctx context.Context, id string,
    p models.ProjectUpdateModel) (models.ProjectModel, error) {

    // Business rule: Validate UUID format
    isValidID := common.ValidateUUID(id)
    if !isValidID {
        return models.ProjectModel{},
            huma.Error400BadRequest("Must provide UUID format")
    }

    // Delegate to repository
    return ps.projectRepo.Update(ctx, id, p)
}

// Delete project with ID validation
func (ps *ProjectService) Delete(ctx context.Context, id string) error {
    // Business rule: Validate UUID format
    isValidID := common.ValidateUUID(id)
    if !isValidID {
        return huma.Error400BadRequest("Must provide UUID format")
    }

    // Delegate to repository
    return ps.projectRepo.Delete(ctx, id)
}
```

**Key Points**:

- All methods with ID parameters validate UUID format
- Simple delegation to repository after validation
- Could be extended with more complex business logic

### Example: Auth Service (`auth_service.go`)

```go
type AuthService struct {
    authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
    return AuthService{authRepo}
}

// Login with credential validation
func (as *AuthService) Login(data models.AuthLoginInputModel,
    env configs.Environment) (models.AuthLoginOutputModel, error) {

    // Business rule: Validate credentials against environment
    isValid := env.Admin.Username != "" &&
               env.Admin.Password != "" &&
               data.Username == env.Admin.Username &&
               data.Password == env.Admin.Password

    if !isValid {
        return models.AuthLoginOutputModel{},
            huma.Error401Unauthorized("Invalid credentials")
    }

    // Generate tokens via repository
    accessToken, refreshToken, err := as.authRepo.GenerateFreshTokens(data)
    if err != nil {
        return models.AuthLoginOutputModel{}, err
    }

    // Return success response
    return models.AuthLoginOutputModel{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        Username:     data.Username,
    }, nil
}

// Refresh access token
func (as *AuthService) Refresh(data models.AuthRefreshInputModel)
    (models.AuthRefreshOutputModel, error) {

    // Simple delegation - token validation in repository
    newAccessToken, err := as.authRepo.RegenerateAccessToken(data.RefreshToken)
    if err != nil {
        return models.AuthRefreshOutputModel{}, err
    }

    return models.AuthRefreshOutputModel{
        AccessToken: newAccessToken,
    }, nil
}
```

---

## 4. Resource Layer

**Location**: `internal/resources/`

### Purpose

Handle HTTP concerns and API endpoints:

- Define API routes and operations
- Parse HTTP requests
- Call service methods
- Format HTTP responses
- Provide OpenAPI metadata

### Characteristics

- **HTTP Expert**: Knows about requests/responses
- **Thin Layer**: Minimal logic, delegates to services
- **Route Registration**: Defines all endpoints
- **OpenAPI Documentation**: Provides API metadata

### Structure

```go
type <Entity>Resource struct {
    <entity>Srv services.<Entity>Service
}

func New<Entity>Resource(<entity>Srv services.<Entity>Service) <Entity>Resource {
    return <Entity>Resource{<entity>Srv}
}

// Register routes
func (r <Entity>Resource) Routes(api huma.API) {
    huma.Register(api, huma.Operation{...}, r.handlerMethod)
}

// Handler methods
func (r <Entity>Resource) handlerMethod(ctx context.Context,
    input *InputStruct) (*OutputStruct, error)
```

### Example: Project Resource (`project_resource.go`)

#### Constructor & Route Registration

```go
type ProjectResource struct {
    projectSrv services.ProjectService
}

func NewProjectResource(projectSrv services.ProjectService) ProjectResource {
    return ProjectResource{projectSrv}
}

func (pr ProjectResource) Routes(api huma.API) {
    // List projects
    huma.Register(api, huma.Operation{
        OperationID: "project-get-paginated",
        Method:      http.MethodGet,
        Path:        "/projects",
        Summary:     "Get Projects",
        Tags:        []string{"Project"},
        Security:    []map[string][]string{{"bearer": {}}},
    }, pr.getPaginated)

    // Get single project
    huma.Register(api, huma.Operation{
        OperationID: "project-get-detail",
        Method:      http.MethodGet,
        Path:        "/projects/{projectId}",
        Summary:     "Get Project detail",
        Tags:        []string{"Project"},
        Security:    []map[string][]string{{"bearer": {}}},
    }, pr.getDetail)

    // Create project
    huma.Register(api, huma.Operation{
        OperationID: "project-create",
        Method:      http.MethodPost,
        Path:        "/projects/{projectId}",
        Summary:     "Create single project",
        Tags:        []string{"Project"},
        Security:    []map[string][]string{{"bearer": {}}},
    }, pr.create)

    // Update project
    huma.Register(api, huma.Operation{
        OperationID: "project-update",
        Method:      http.MethodPatch,
        Path:        "/projects/{projectId}",
        Summary:     "Update single project",
        Tags:        []string{"Project"},
        Security:    []map[string][]string{{"bearer": {}}},
    }, pr.update)

    // Delete project
    huma.Register(api, huma.Operation{
        OperationID: "project-delete",
        Method:      http.MethodDelete,
        Path:        "/projects/{projectId}",
        Summary:     "Delete single project",
        Tags:        []string{"Project"},
        Security:    []map[string][]string{{"bearer": {}}},
    }, pr.delete)
}
```

#### Handler Methods

```go
// Get paginated list
func (pr ProjectResource) getPaginated(ctx context.Context,
    input *models.ProjectSearchModel)
    (*struct{ Body models.ProjectPaginatedModel }, error) {

    // Call service
    respSrv, err := pr.projectSrv.GetPaginated(ctx, *input)
    if err != nil {
        return nil, err
    }

    // Wrap and return response
    return &struct{ Body models.ProjectPaginatedModel }{
        Body: respSrv,
    }, nil
}

// Get single project
func (pr ProjectResource) getDetail(ctx context.Context,
    input *struct{ Path string `path:"projectId"` })
    (*struct{ Body models.ProjectModel }, error) {

    // Call service with path parameter
    respSrv, err := pr.projectSrv.GetDetail(ctx, input.Path)
    if err != nil {
        return nil, err
    }

    // Wrap and return response
    return &struct{ Body models.ProjectModel }{
        Body: respSrv,
    }, nil
}

// Create project
func (pr ProjectResource) create(ctx context.Context,
    input *struct{ Body models.ProjectCreateModel })
    (*struct{ Body models.ProjectModel }, error) {

    // Call service with request body
    respSrc, err := pr.projectSrv.Create(ctx, input.Body)
    if err != nil {
        return nil, err
    }

    // Wrap and return response
    return &struct{ Body models.ProjectModel }{
        Body: respSrc,
    }, nil
}

// Update project
func (pr ProjectResource) update(ctx context.Context,
    input *struct{
        Path string `path:"projectId"`
        Body models.ProjectUpdateModel
    }) (*struct{ Body models.ProjectModel }, error) {

    // Call service with path and body
    respSrc, err := pr.projectSrv.Update(ctx, input.Path, input.Body)
    if err != nil {
        return nil, err
    }

    // Wrap and return response
    return &struct{ Body models.ProjectModel }{
        Body: respSrc,
    }, nil
}

// Delete project
func (pr ProjectResource) delete(ctx context.Context,
    input *struct{ Path string `path:"projectId"` })
    (*struct{}, error) {

    // Call service
    err := pr.projectSrv.Delete(ctx, input.Path)
    if err != nil {
        return nil, err
    }

    // Return empty response
    return nil, nil
}
```

---

## Pattern Benefits

### 1. **Separation of Concerns**

- Each layer has a single, well-defined responsibility
- Changes in one layer don't cascade to others

### 2. **Testability**

- Each layer can be tested independently
- Easy to mock dependencies (repositories, services)

### 3. **Maintainability**

- Clear structure makes code easy to navigate
- New developers can understand the flow quickly

### 4. **Scalability**

- Easy to add new entities (follow the same pattern)
- Can extract services into microservices if needed

### 5. **Reusability**

- Services can be called from multiple resources
- Repositories can be shared across services

### 6. **Flexibility**

- Swap implementations (e.g., change database)
- Add caching at repository level
- Add authorization at service level

---

## Common Patterns

### Dependency Injection

```go
// Constructor functions inject dependencies
repo := repositories.NewProjectRepository(dbPool)
service := services.NewProjectService(repo)
resource := resources.NewProjectResource(service)
```

### Error Handling

```go
// Each layer returns errors up the chain
// Resource → Service → Repository → Database
// Error propagates back: Database → Repository → Service → Resource → HTTP
```

### Context Propagation

```go
// Context flows down through all layers
func (r Resource) handler(ctx context.Context, ...) {
    service.Method(ctx, ...)  // Pass context down
}
```

### Validation Strategy

- **Model Layer**: Struct tags (format, length, enum)
- **Service Layer**: Business rules (UUID validation, authorization)
- **Repository Layer**: Database constraints (foreign keys, uniqueness)

---

## Anti-Patterns to Avoid

### ❌ **Don't Skip Layers**

```go
// BAD: Resource calling Repository directly
func (r Resource) handler(ctx context.Context, ...) {
    return r.repository.Get(ctx, id)  // Skips service layer
}
```

### ❌ **Don't Put Business Logic in Resources**

```go
// BAD: Validation in resource
func (r Resource) handler(ctx context.Context, input *Input) {
    if !isValidUUID(input.ID) {  // Should be in service
        return nil, error
    }
}
```

### ❌ **Don't Put SQL in Services**

```go
// BAD: SQL query in service
func (s Service) GetData(ctx context.Context) {
    sql := "SELECT * FROM table"  // Should be in repository
    s.db.Query(sql)
}
```

### ❌ **Don't Put HTTP Logic in Services**

```go
// BAD: HTTP status codes in service
func (s Service) Process() error {
    return huma.Error404NotFound("...")  // OK for now, but could be improved
}
```

---

## Summary

The RSRM pattern provides:

- **Clear boundaries** between concerns
- **Consistent structure** across all entities
- **Maintainable code** that's easy to extend
- **Testable components** at every layer

Each layer has a specific job, and by respecting these boundaries, the codebase remains clean, organized, and professional.
