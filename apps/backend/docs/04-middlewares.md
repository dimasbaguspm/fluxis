# Middlewares Documentation

## Overview

Middlewares in the Fluxis backend are functions that intercept HTTP requests before they reach the final handler. They provide cross-cutting concerns like authentication, logging, and request validation. The application uses the **Huma middleware system** which integrates seamlessly with the API framework.

## Middleware Architecture

### Middleware Flow

```
HTTP Request
    ↓
HTTP Router
    ↓
┌──────────────────────────┐
│   Middleware Chain       │
│                          │
│  1. Middleware 1         │
│  2. Middleware 2         │
│  3. Middleware N         │
└──────────────────────────┘
    ↓
Route Handler (Resource)
    ↓
Service Layer
    ↓
Repository Layer
    ↓
HTTP Response
```

### Middleware Function Signature

```go
func MiddlewareName(api huma.API) func(huma.Context, func(huma.Context)) {
    return func(ctx huma.Context, next func(huma.Context)) {
        // Pre-processing logic (before handler)

        // Validation/checks
        if !isValid {
            // Short-circuit: return error without calling next
            huma.WriteErr(api, ctx, statusCode, errorMessage)
            return
        }

        // Call next middleware or handler
        next(ctx)

        // Post-processing logic (after handler)
    }
}
```

**Key Components**:

- `api huma.API` - API instance for error handling
- `ctx huma.Context` - Request context with headers, body, etc.
- `next func(huma.Context)` - Function to call the next middleware/handler

### Middleware Behavior

1. **Pre-processing**: Execute logic before the request reaches the handler
2. **Short-circuit**: Return early without calling `next()` if validation fails
3. **Pass-through**: Call `next(ctx)` to continue to the next middleware/handler
4. **Post-processing**: Execute logic after the handler completes

---

## Current Middleware: Session Middleware

**Location**: `internal/middlewares/session_middleware.go`

### Purpose

Validate JWT access tokens for protected routes, ensuring only authenticated users can access private endpoints.

### Implementation

```go
package middlewares

import (
    "github.com/danielgtaylor/huma/v2"
    "github.com/dimasbaguspm/fluxis/internal/repositories"
)

var authRepo = repositories.AuthRepository{}

func SessionMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
    return func(ctx huma.Context, next func(huma.Context)) {
        // 1. Extract Authorization header
        cat := ctx.Header("Authorization")

        // 2. Strip "Bearer " prefix (if present)
        if len(cat) > 7 && cat[:7] == "Bearer " {
            cat = cat[7:]
        }

        // 3. Validate JWT token
        isValid := authRepo.IsTokenValid(cat)

        // 4. If invalid, return error and stop processing
        if !isValid {
            huma.WriteErr(api, ctx,
                repositories.AuthErrorInvalidAccessToken.GetStatus(),
                repositories.AuthErrorInvalidAccessToken.Error())
            return
        }

        // 5. Token is valid, proceed to next handler
        next(ctx)
    }
}
```

### Step-by-Step Breakdown

#### Step 1: Extract Authorization Header

```go
cat := ctx.Header("Authorization")
```

Retrieves the `Authorization` header from the HTTP request:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Step 2: Strip Bearer Prefix

```go
if len(cat) > 7 && cat[:7] == "Bearer " {
    cat = cat[7:]
}
```

Removes the `Bearer ` prefix to get the raw JWT token:

```
Before: "Bearer eyJhbGciOiJI..."
After:  "eyJhbGciOiJI..."
```

**Why?** JWT validation expects only the token, not the authentication scheme.

#### Step 3: Validate Token

```go
isValid := authRepo.IsTokenValid(cat)
```

Calls the repository method to validate the token (see [Token Validation Logic](#token-validation-logic) below).

#### Step 4: Handle Invalid Token

```go
if !isValid {
    huma.WriteErr(api, ctx,
        repositories.AuthErrorInvalidAccessToken.GetStatus(),
        repositories.AuthErrorInvalidAccessToken.Error())
    return
}
```

If validation fails:

- Write an error response (400 Bad Request)
- **Return immediately** without calling `next(ctx)`
- Request processing stops here

**Response**:

```json
{
  "status": 400,
  "title": "Bad Request",
  "detail": "Invalid access token"
}
```

#### Step 5: Continue to Handler

```go
next(ctx)
```

If validation succeeds, call the next middleware or route handler.

---

### Token Validation Logic

**Location**: `internal/repositories/auth_repository.go`

```go
func (ar AuthRepository) IsTokenValid(token string) bool {
    // 1. Parse JWT token with validation callback
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        // 2. Verify signing method is HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, AuthErrorInvalidSigningMethod
        }
        // 3. Return secret key for signature verification
        return []byte(secretJWT), nil
    })

    // 4. Check for parsing errors
    if err != nil {
        return false
    }

    // 5. Verify token is valid (not expired, signature matches)
    if !parsedToken.Valid {
        return false
    }

    return true
}
```

**Validation Steps**:

1. **Parse Token**: Decode JWT structure
2. **Verify Algorithm**: Ensure HMAC-SHA256 signing method
3. **Verify Signature**: Check token hasn't been tampered with
4. **Check Expiration**: Ensure token hasn't expired
5. **Return Result**: `true` if all checks pass

**Security Notes**:

- Validates signature using secret key (`secretJWT`)
- Checks expiration time (set during token generation)
- Rejects tokens with invalid signing methods (prevents algorithm substitution attacks)

---

## Middleware Registration

### Public Routes (No Middleware)

```go
func RegisterPublicRoutes(api huma.API, pgx *pgxpool.Pool) {
    // No middleware applied
    authRepo := repositories.NewAuthRepository(pgx)
    authSrv := services.NewAuthService(authRepo)
    resources.NewAuthResource(authSrv).Routes(api)
}
```

**Accessible without authentication**:

- `POST /auth/login`
- `POST /auth/refresh`

### Private Routes (With Session Middleware)

```go
func RegisterPrivateRoutes(api huma.API, pgx *pgxpool.Pool) {
    // Apply middleware to ALL routes registered after this line
    api.UseMiddleware(middlewares.SessionMiddleware(api))

    projectRepo := repositories.NewProjectRepository(pgx)
    projectSrv := services.NewProjectService(projectRepo)
    resources.NewProjectResource(projectSrv).Routes(api)
}
```

**Requires authentication**:

- `GET /projects`
- `GET /projects/{projectId}`
- `POST /projects/{projectId}`
- `PATCH /projects/{projectId}`
- `DELETE /projects/{projectId}`

---

## Request Flow with Middleware

### Example: Protected Request Flow

```
1. Client Request
   GET /projects
   Authorization: Bearer <valid_token>

2. HTTP Router
   Routes to /projects handler

3. SessionMiddleware
   • Extract Authorization header
   • Strip "Bearer " prefix
   • Validate JWT token
   • Token valid? → Call next()

4. ProjectResource.getPaginated()
   • Parse query parameters
   • Call service layer

5. ProjectService.GetPaginated()
   • Validate business rules
   • Call repository

6. ProjectRepository.GetPaginated()
   • Execute SQL query
   • Return results

7. Response Flow Back
   Repository → Service → Resource → HTTP Response
```

### Example: Blocked Request Flow

```
1. Client Request
   GET /projects
   Authorization: Bearer <invalid_token>

2. HTTP Router
   Routes to /projects handler

3. SessionMiddleware
   • Extract Authorization header
   • Strip "Bearer " prefix
   • Validate JWT token
   • Token invalid? → WriteErr() and return

4. Request Stopped
   Handler never called

5. Error Response
   HTTP 400 Bad Request
   { "status": 400, "detail": "Invalid access token" }
```

---

## Middleware Best Practices

### 1. **Order Matters**

Middlewares are executed in the order they're registered:

```go
api.UseMiddleware(loggingMiddleware)   // Runs first
api.UseMiddleware(authMiddleware)      // Runs second
api.UseMiddleware(rateLimitMiddleware) // Runs third
```

### 2. **Short-Circuit on Failure**

Always return immediately after writing an error:

```go
if !isValid {
    huma.WriteErr(api, ctx, statusCode, message)
    return  // Don't call next()
}
```

### 3. **Stateless Validation**

Middlewares should be stateless (JWT tokens are self-contained):

```go
// Good: Stateless validation
isValid := validateToken(token)

// Bad: Stateful validation (requires session store)
session := getSessionFromStore(sessionId)
```

### 4. **Performance Considerations**

- Minimize expensive operations (database queries, external API calls)
- Cache validation results if possible
- Use early returns to avoid unnecessary processing

### 5. **Error Handling**

Use consistent error responses:

```go
// Define error constants
var (
    AuthErrorInvalidToken = huma.Error400BadRequest("Invalid access token")
    AuthErrorExpiredToken = huma.Error401Unauthorized("Token expired")
    AuthErrorMissingToken = huma.Error401Unauthorized("Missing authorization header")
)

// Use in middleware
if tokenMissing {
    huma.WriteErr(api, ctx, AuthErrorMissingToken.GetStatus(), AuthErrorMissingToken.Error())
    return
}
```

---

## Future Middleware Enhancements

### 1. **Logging Middleware**

```go
func LoggingMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
    return func(ctx huma.Context, next func(huma.Context)) {
        start := time.Now()

        // Log request
        slog.Info("Request started",
            "method", ctx.Method(),
            "path", ctx.URL().Path,
        )

        next(ctx)

        // Log response
        slog.Info("Request completed",
            "duration", time.Since(start),
            "status", ctx.Status(),
        )
    }
}
```

### 2. **Rate Limiting Middleware**

```go
func RateLimitMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
    return func(ctx huma.Context, next func(huma.Context)) {
        clientIP := ctx.RemoteAddr()

        if isRateLimited(clientIP) {
            huma.WriteErr(api, ctx, 429, "Too many requests")
            return
        }

        next(ctx)
    }
}
```

### 3. **CORS Middleware**

```go
func CORSMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
    return func(ctx huma.Context, next func(huma.Context)) {
        ctx.SetHeader("Access-Control-Allow-Origin", "*")
        ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
        ctx.SetHeader("Access-Control-Allow-Headers", "Authorization, Content-Type")

        if ctx.Method() == "OPTIONS" {
            ctx.SetStatus(204)
            return
        }

        next(ctx)
    }
}
```

### 4. **Request ID Middleware**

```go
func RequestIDMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
    return func(ctx huma.Context, next func(huma.Context)) {
        requestID := uuid.New().String()
        ctx.SetHeader("X-Request-ID", requestID)

        // Add to context for logging
        ctx = context.WithValue(ctx, "request_id", requestID)

        next(ctx)
    }
}
```

### 5. **Role-Based Access Control (RBAC) Middleware**

```go
func RBACMiddleware(requiredRole string) func(huma.API) func(huma.Context, func(huma.Context)) {
    return func(api huma.API) func(huma.Context, func(huma.Context)) {
        return func(ctx huma.Context, next func(huma.Context)) {
            token := extractToken(ctx)
            claims := parseTokenClaims(token)

            userRole := claims["role"].(string)
            if userRole != requiredRole {
                huma.WriteErr(api, ctx, 403, "Insufficient permissions")
                return
            }

            next(ctx)
        }
    }
}
```

---

## Testing Middlewares

### Unit Testing

```go
func TestSessionMiddleware(t *testing.T) {
    api := mockHumaAPI()
    middleware := SessionMiddleware(api)

    // Test valid token
    t.Run("ValidToken", func(t *testing.T) {
        ctx := createMockContext("Bearer valid_token")
        nextCalled := false

        middleware(ctx, func(ctx huma.Context) {
            nextCalled = true
        })

        assert.True(t, nextCalled)
    })

    // Test invalid token
    t.Run("InvalidToken", func(t *testing.T) {
        ctx := createMockContext("Bearer invalid_token")
        nextCalled := false

        middleware(ctx, func(ctx huma.Context) {
            nextCalled = true
        })

        assert.False(t, nextCalled)
        assert.Equal(t, 400, ctx.Status())
    })
}
```

### Integration Testing

```go
func TestProtectedEndpoint(t *testing.T) {
    // Setup server with middleware
    server := setupTestServer()

    // Test without token
    resp := httpGet(server, "/projects", "")
    assert.Equal(t, 400, resp.StatusCode)

    // Test with valid token
    token := generateValidToken()
    resp = httpGet(server, "/projects", "Bearer "+token)
    assert.Equal(t, 200, resp.StatusCode)
}
```

---

## Common Issues and Solutions

### Issue 1: Middleware Not Applied

**Problem**: Routes not protected despite middleware registration

**Solution**: Ensure middleware is registered BEFORE route registration

```go
// Correct order
api.UseMiddleware(authMiddleware)
resource.Routes(api)  // Routes registered after middleware

// Wrong order
resource.Routes(api)  // Routes registered first
api.UseMiddleware(authMiddleware)  // Middleware applied after
```

### Issue 2: Token Prefix Handling

**Problem**: Validation fails due to "Bearer " prefix

**Solution**: Always strip prefix before validation

```go
token := ctx.Header("Authorization")
if len(token) > 7 && token[:7] == "Bearer " {
    token = token[7:]
}
```

### Issue 3: Missing Error Handling

**Problem**: Middleware crashes on nil values

**Solution**: Add defensive checks

```go
token := ctx.Header("Authorization")
if token == "" {
    huma.WriteErr(api, ctx, 401, "Missing authorization header")
    return
}
```

---

## Summary

**Current Middlewares**:

- ✅ `SessionMiddleware` - JWT authentication for protected routes

**Middleware Characteristics**:

- Intercept requests before handlers
- Perform validation and authentication
- Short-circuit on failure (don't call `next()`)
- Stateless operation (JWT tokens)

**Best Practices**:

- Register middleware before routes
- Use consistent error responses
- Keep middleware stateless and fast
- Test middleware independently

**Future Enhancements**:

- Logging middleware for request tracking
- Rate limiting for API protection
- CORS configuration for browser clients
- RBAC for fine-grained access control
