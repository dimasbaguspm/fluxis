# Migrations Workflow Documentation

## Overview

Database migrations in the Fluxis backend provide a systematic way to manage database schema changes through version-controlled SQL files. The application uses the **golang-migrate** library to handle migrations automatically on startup, ensuring the database schema is always up-to-date with the application code.

---

## Migration System Architecture

### Components

```
┌─────────────────────────────────────────┐
│     Application Startup (main.go)      │
└────────────────┬────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────┐
│   Migration Config (migration.go)      │
│   • Initialize migrator                 │
│   • Configure source (file://)          │
│   • Configure database (postgres://)    │
└────────────────┬────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────┐
│    golang-migrate/migrate/v4            │
│   • Read migration files                │
│   • Track migration versions            │
│   • Execute SQL statements              │
│   • Maintain schema_migrations table    │
└────────────────┬────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────┐
│       Migration Files (migrations/)     │
│   • 000001_setup.up.sql                 │
│   • 000001_setup.down.sql               │
│   • 000002_*.up.sql (future)            │
│   • 000002_*.down.sql (future)          │
└────────────────┬────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────┐
│          PostgreSQL Database            │
│   • Execute DDL statements              │
│   • Store migration version             │
│   • Apply schema changes                │
└─────────────────────────────────────────┘
```

---

## Migration Configuration

**Location**: `internal/configs/migration.go`

### Implementation

```go
package configs

import (
    "errors"
    "fmt"
    "path/filepath"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

type migration struct {
    env Environment
}

func Migration(env Environment) migration {
    return migration{env}
}

func (m migration) Up() error {
    return proceedMigration(m, "up")
}

func (m migration) Down() error {
    return proceedMigration(m, "down")
}

func proceedMigration(m migration, option string) error {
    // 1. Initialize migrator with source and database
    migrator, err := migrate.New(
        fmt.Sprintf("file://%s", filepath.Clean("migrations")),
        m.env.Database.Url)

    if err != nil {
        return fmt.Errorf("Failed to setup migration: %w", err)
    }
    defer migrator.Close()

    var postErr error

    // 2. Execute migration based on direction
    switch option {
    case "up":
        if err := migrator.Up(); err != nil {
            if errors.Is(err, migrate.ErrNoChange) {
                return nil  // No new migrations to apply
            }
            postErr = fmt.Errorf("failed to run migrations: %w", err)
        }
    case "down":
        if err := migrator.Down(); err != nil {
            if errors.Is(err, migrate.ErrNoChange) {
                return nil  // No migrations to rollback
            }
            postErr = fmt.Errorf("failed to run migrations: %w", err)
        }
    }

    if postErr != nil {
        return postErr
    }

    return nil
}
```

### Configuration Details

#### Source Configuration

```go
fmt.Sprintf("file://%s", filepath.Clean("migrations"))
```

- **Protocol**: `file://` - Read migrations from filesystem
- **Path**: `migrations/` - Relative to application root
- **Behavior**: Reads all `.sql` files in the directory

#### Database Configuration

```go
m.env.Database.Url
// Example: "postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

- **Protocol**: `postgres://` - PostgreSQL database driver
- **Connection**: From environment variables
- **SSL Mode**: Configurable via connection string

---

## Migration Execution on Startup

**Location**: `cmd/app/main.go`

```go
func main() {
    // ... (earlier initialization)

    // Initialize migration system
    migration := configs.Migration(env)

    // Run migrations
    slog.Info("Performing migration")
    if err := migration.Up(); err != nil {
        slog.Error("Failed to run migration", "error", err.Error())
        panic(err)
    }
    slog.Info("DB migration completed")

    // ... (continue startup)
}
```

### Execution Flow

1. **Initialize Migrator**: Create migration instance with environment config
2. **Execute Up Migrations**: Apply all pending migrations
3. **Handle Results**:
   - **Success**: Log completion and continue
   - **No Changes**: Silently continue (no new migrations)
   - **Failure**: Log error and panic (prevent app startup)

### Why Execute on Startup?

- **Automatic Schema Updates**: Developers don't need to manually run migrations
- **Environment Parity**: Dev, staging, and production always have correct schema
- **Fail Fast**: Database schema issues detected immediately
- **Simplified Deployment**: Single deployment step (no separate migration command)

---

## Migration Files Structure

**Location**: `migrations/`

### File Naming Convention

```
{version}_{name}.{direction}.sql
```

**Components**:

- `{version}`: Sequential number with leading zeros (e.g., `000001`, `000002`)
- `{name}`: Descriptive name (e.g., `setup`, `add_users_table`)
- `{direction}`: Either `up` or `down`

**Examples**:

```
000001_setup.up.sql           # First migration - apply schema
000001_setup.down.sql         # First migration - rollback schema
000002_add_users_table.up.sql
000002_add_users_table.down.sql
```

### Migration Pairs

Every migration must have both:

- **`.up.sql`**: Forward migration (apply changes)
- **`.down.sql`**: Reverse migration (rollback changes)

This allows:

- **Rolling forward**: Apply new schema changes
- **Rolling back**: Revert to previous schema state

---

## Current Migrations

### Migration 000001: Initial Setup

#### Up Migration (`000001_setup.up.sql`)

**Purpose**: Create initial database schema for project management system

**Schema Components**:

##### 1. Project Status Enum

```sql
CREATE TYPE project_status AS ENUM ('active', 'paused', 'archived');
```

Defines allowed project status values.

##### 2. Projects Table

```sql
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    status project_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
```

**Features**:

- UUID primary key (auto-generated)
- Enum type for status
- Soft delete support (`deleted_at`)
- Automatic timestamps

##### 3. Statuses Table (Kanban Columns)

```sql
CREATE TABLE statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Features**:

- Per-project status columns (for Kanban boards)
- Cascade delete when project is deleted
- Position-based ordering
- Default status flag

##### 4. Default Statuses Trigger

```sql
CREATE OR REPLACE FUNCTION create_default_statuses_for_project()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO statuses (project_id, name, slug, position, is_default)
    VALUES (NEW.id, 'Todo', 'todo', 0, true),
           (NEW.id, 'In Progress', 'in_progress', 1, false),
           (NEW.id, 'Done', 'done', 2, false);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_default_statuses
AFTER INSERT ON projects
FOR EACH ROW
EXECUTE FUNCTION create_default_statuses_for_project();
```

**Behavior**: Automatically creates three default statuses when a new project is created:

- "Todo" (default, position 0)
- "In Progress" (position 1)
- "Done" (position 2)

##### 5. Tasks Table

```sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    details TEXT,
    status_id UUID REFERENCES statuses(id) ON DELETE SET NULL,
    priority INT NOT NULL DEFAULT 1 CHECK (priority >= 0),
    due_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
```

**Features**:

- Foreign key to projects (cascade delete)
- Foreign key to statuses (set null on delete)
- Priority with validation (>= 0)
- Optional due date
- Soft delete support

##### 6. Activity Logs Table

```sql
CREATE TABLE logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
    status_id UUID REFERENCES statuses(id) ON DELETE SET NULL,
    entry TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Features**:

- Audit trail for project/task activities
- Optional task and status references
- Immutable records (no update/delete timestamps)

#### Down Migration (`000001_setup.down.sql`)

**Purpose**: Completely reverse the initial schema setup

```sql
BEGIN;

-- Remove trigger
DROP TRIGGER IF EXISTS trg_create_default_statuses ON projects;

-- Remove function
DROP FUNCTION IF EXISTS create_default_statuses_for_project() CASCADE;

-- Drop tables in dependency order
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS statuses;
DROP TABLE IF EXISTS projects;

-- Drop enum type
DROP TYPE IF EXISTS project_status;

COMMIT;
```

**Characteristics**:

- Wrapped in transaction (atomic rollback)
- Drops in reverse dependency order
- Uses `IF EXISTS` for idempotency
- Removes triggers, functions, tables, and types

---

## Migration Tracking

### Schema Migrations Table

golang-migrate automatically creates a tracking table:

```sql
CREATE TABLE schema_migrations (
    version bigint NOT NULL PRIMARY KEY,
    dirty boolean NOT NULL
);
```

**Fields**:

- `version`: Current migration version (e.g., `1` for `000001_setup`)
- `dirty`: Indicates incomplete migration (used for recovery)

**Behavior**:

- Updated automatically by golang-migrate
- Prevents re-applying completed migrations
- Tracks partial failures

### Example State

```sql
SELECT * FROM schema_migrations;

 version | dirty
---------+-------
       1 | f
```

This indicates:

- Migration version `1` (000001_setup) has been applied
- Migration completed successfully (`dirty = false`)

---

## Migration Workflows

### 1. First Application Startup (Fresh Database)

```
1. Application starts
2. Migration system initializes
3. Detects no schema_migrations table
4. Applies all .up.sql files in order:
   • 000001_setup.up.sql → Creates schema
5. Creates schema_migrations table
6. Inserts version 1 with dirty=false
7. Application continues startup
```

**Database State**:

- All tables created: `projects`, `statuses`, `tasks`, `logs`
- Trigger and function installed
- Migration version: 1

### 2. Application Startup (Existing Database, No New Migrations)

```
1. Application starts
2. Migration system initializes
3. Reads schema_migrations table → version = 1
4. Scans migrations/ directory
5. Finds only 000001_* (version 1)
6. No new migrations detected
7. Returns migrate.ErrNoChange (silently continues)
8. Application continues startup
```

**Log Output**:

```
INFO Performing migration
INFO DB migration completed
```

### 3. Application Startup (New Migration Available)

```
Scenario: Developer adds 000002_add_users_table.up.sql

1. Application starts
2. Migration system initializes
3. Reads schema_migrations table → version = 1
4. Scans migrations/ directory
5. Finds 000002_add_users_table.up.sql (version 2)
6. Executes 000002_add_users_table.up.sql
7. Updates schema_migrations → version = 2
8. Application continues startup
```

**Database State**:

- Version 1 schema remains intact
- Version 2 changes applied
- Migration version: 2

### 4. Migration Failure

```
Scenario: 000002_add_users_table.up.sql has SQL error

1. Application starts
2. Migration system initializes
3. Reads schema_migrations table → version = 1
4. Attempts to execute 000002_add_users_table.up.sql
5. SQL error encountered
6. Sets dirty=true in schema_migrations
7. Returns error to application
8. Application panics and stops
```

**Recovery Steps**:

1. Fix SQL error in migration file
2. Manually rollback partial changes (if any)
3. Reset dirty flag: `UPDATE schema_migrations SET dirty=false`
4. Restart application

---

## Creating New Migrations

### Step-by-Step Guide

#### 1. Determine Next Version Number

```bash
ls migrations/
# Output: 000001_setup.down.sql  000001_setup.up.sql
# Next version: 000002
```

#### 2. Create Migration Files

```bash
touch migrations/000002_add_users_table.up.sql
touch migrations/000002_add_users_table.down.sql
```

#### 3. Write Up Migration

```sql
-- migrations/000002_add_users_table.up.sql

-- Add users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add user_id to projects table
ALTER TABLE projects ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE SET NULL;

-- Create index for faster lookups
CREATE INDEX idx_projects_user_id ON projects(user_id);
```

#### 4. Write Down Migration

```sql
-- migrations/000002_add_users_table.down.sql

-- Remove index
DROP INDEX IF EXISTS idx_projects_user_id;

-- Remove column from projects
ALTER TABLE projects DROP COLUMN IF EXISTS user_id;

-- Remove users table
DROP TABLE IF EXISTS users;
```

#### 5. Test Migration

**Test Up**:

```bash
# Start application
go run cmd/app/main.go

# Verify schema
psql -d fluxis -c "\dt"
# Should show: projects, statuses, tasks, logs, users
```

**Test Down** (if needed):

```go
// Temporary test in main.go
migration := configs.Migration(env)
if err := migration.Down(); err != nil {
    panic(err)
}
```

---

## Best Practices

### 1. **Atomic Migrations**

Each migration should be self-contained and atomic:

```sql
-- Good: Wrapped in transaction
BEGIN;
CREATE TABLE users (...);
ALTER TABLE projects ADD COLUMN user_id UUID;
COMMIT;

-- Bad: No transaction (partial failures possible)
CREATE TABLE users (...);
ALTER TABLE projects ADD COLUMN user_id UUID;
```

### 2. **Idempotent Down Migrations**

Always use `IF EXISTS` in down migrations:

```sql
-- Good: Safe to run multiple times
DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_users_email;

-- Bad: Fails if already dropped
DROP TABLE users;
```

### 3. **Test Both Directions**

Always test:

- ✅ Up migration applies successfully
- ✅ Down migration reverts correctly
- ✅ Re-applying up migration after down works

### 4. **Backward Compatibility**

When possible, make additive changes:

```sql
-- Good: Add nullable column (backward compatible)
ALTER TABLE projects ADD COLUMN archived_at TIMESTAMPTZ;

-- Risky: Add required column (breaks existing queries)
ALTER TABLE projects ADD COLUMN owner_id UUID NOT NULL;
```

### 5. **Data Migrations**

For data transformations, include both schema and data changes:

```sql
-- Add new column
ALTER TABLE projects ADD COLUMN status_text TEXT;

-- Migrate existing data
UPDATE projects SET status_text = status::text;

-- Make required after migration
ALTER TABLE projects ALTER COLUMN status_text SET NOT NULL;
```

### 6. **Sequential Versioning**

Always use sequential version numbers:

```
000001_setup.sql           ✅
000002_add_users.sql       ✅
000003_add_tasks_index.sql ✅

000001_setup.sql           ❌
000003_add_users.sql       ❌ (skips 2)
000002_add_tasks_index.sql ❌ (out of order)
```

---

## Troubleshooting

### Issue 1: Dirty Migration State

**Symptom**:

```
Error: Dirty database version 2. Fix and force version.
```

**Cause**: Migration failed partway through, leaving database in inconsistent state

**Solution**:

```sql
-- 1. Inspect database state
SELECT * FROM schema_migrations;

-- 2. Manually fix any partial changes from failed migration

-- 3. Reset dirty flag
UPDATE schema_migrations SET dirty = false;

-- 4. Fix migration file and restart application
```

### Issue 2: Migration Out of Order

**Symptom**:

```
Error: Migration version 3 is less than latest version 5
```

**Cause**: Attempting to apply an older migration after newer ones

**Solution**: Renumber migration files to maintain sequential order

### Issue 3: File Not Found

**Symptom**:

```
Error: Failed to setup migration: source: file does not exist
```

**Cause**: Migrations folder not found (wrong working directory)

**Solution**:

- Ensure application runs from project root
- Verify `migrations/` directory exists
- Check file paths are relative to working directory

### Issue 4: Connection Failure

**Symptom**:

```
Error: Failed to run migration: pq: connection refused
```

**Cause**: Database not accessible

**Solution**:

- Verify database is running
- Check environment variables (`DB_HOST`, `DB_PORT`, etc.)
- Test connection: `psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME`

---

## Manual Migration Commands

While migrations run automatically on startup, you can also run them manually:

### Using Go Code

```go
// Create migration instance
env := configs.NewEnvironment()
migration := configs.Migration(env)

// Apply migrations
if err := migration.Up(); err != nil {
    log.Fatal(err)
}

// Rollback migrations
if err := migration.Down(); err != nil {
    log.Fatal(err)
}
```

### Using migrate CLI (Optional)

Install:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Commands:

```bash
# Apply migrations
migrate -path migrations/ -database "$DATABASE_URL" up

# Rollback one migration
migrate -path migrations/ -database "$DATABASE_URL" down 1

# Force version (recovery)
migrate -path migrations/ -database "$DATABASE_URL" force 1
```

---

## Summary

### Migration System Features

- ✅ **Automatic Execution**: Runs on every application startup
- ✅ **Version Tracking**: Prevents duplicate applications
- ✅ **Atomic Operations**: Transactional safety
- ✅ **Bidirectional**: Support for rollbacks (up/down)
- ✅ **Fail-Safe**: Prevents startup with migration errors

### Current Schema (Version 1)

- ✅ Projects table with soft deletes
- ✅ Dynamic statuses (Kanban columns) per project
- ✅ Tasks with status tracking
- ✅ Activity logs for audit trail
- ✅ Automatic default statuses via trigger

### Best Practices

1. ✅ Always create both `.up.sql` and `.down.sql`
2. ✅ Use transactions for atomicity
3. ✅ Make down migrations idempotent
4. ✅ Test both directions
5. ✅ Keep migrations small and focused
6. ✅ Use sequential versioning

### Future Enhancements

- Add migration testing framework
- Implement migration rollback strategies
- Create migration templates/generators
- Add pre/post migration hooks
- Support multiple environments (dev/staging/prod)
