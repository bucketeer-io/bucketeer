# RFC: Polymorphic Database Storage Layer

## Summary

This RFC proposes implementing PostgreSQL as an alternative primary storage backend alongside MySQL, selected at runtime via configuration. A unified `database.Client` interface abstracts transaction handling, allowing the API layer to call `dbClient.RunInTransactionV2()` without database-specific conditionals. The database selection logic is consolidated in the service constructor.

## Background

Currently, Bucketeer services are tightly coupled to MySQL for primary storage:

1. **API Layer Dependencies**: Services like `PushService` directly hold `mysql.Client` and use MySQL-specific types:

   ```go
   type PushService struct {
       mysqlClient      mysql.Client
       pushStorage      v2ps.PushStorage
       // ...
   }
   ```

2. **Storage Layer Dependencies**: Storage implementations depend on MySQL-specific interfaces and types:

   ```go
   type pushStorage struct {
       qe mysql.QueryExecer  // MySQL-specific interface
   }
   ```

3. **MySQL-Specific Types Throughout**:
   - `mysql.JSONObject` for JSON serialization
   - `mysql.ErrDuplicateEntry`, `mysql.ErrNoRows` for error handling
   - `mysql.ListOptions` for query construction
   - MySQL placeholder syntax (`?` vs `$1`)

We've successfully implemented PostgreSQL as an alternative for the data warehouse (eventcounter), but extending this to primary storage requires a more comprehensive abstraction.

## Goals

1. Enable PostgreSQL as an alternative primary storage backend alongside MySQL
2. Achieve database selection through configuration only - no if-else in API layer business logic
3. Maintain backward compatibility with existing MySQL deployments
4. Provide a clean abstraction that can be extended to other databases in the future
5. Minimize code duplication between MySQL and PostgreSQL implementations

## Design Overview

### Architecture Layers

```
┌─────────────────────────────────────────────────────────────────┐
│                        API Layer                                │
│  (Holds database.Client for transactions)                       │
│  (Holds Storage interface for CRUD operations)                  │
└─────────────────────────────────────────────────────────────────┘
                              │
            ┌─────────────────┴─────────────────┐
            ▼                                   ▼
┌───────────────────────┐       ┌─────────────────────────────────┐
│   database.Client     │       │      Storage Interface          │
│   (for transactions)  │       │  (PushStorage, TagStorage, ...) │
└───────────────────────┘       └─────────────────────────────────┘
            │                                   │
            │                   ┌───────────────┴───────────────┐
            │                   ▼                               ▼
            │       ┌─────────────────────┐       ┌──────────────────────────┐
            │       │ MySQL Implementation│       │ PostgreSQL Implementation│
            │       │ (mysql_push.go)     │       │ (postgres_push.go)       │
            │       └─────────────────────┘       └──────────────────────────┘
            │                   │                               │
            ▼                   ▼                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                    pkg/storage/v2/database                      │
│  (database.Client interface + MySQL/Postgres adapters)          │
└─────────────────────────────────────────────────────────────────┘
            │                                   │
            ▼                                   ▼
┌─────────────────────────┐       ┌─────────────────────────┐
│   pkg/storage/v2/mysql  │       │ pkg/storage/v2/postgres │
└─────────────────────────┘       └─────────────────────────┘
```

### Key Components

#### 1. Unified Database Client Interface

Create a `database.Client` interface that both MySQL and PostgreSQL clients implement:

```go
// pkg/storage/v2/database/client.go
package database

import "context"

// Client is the unified database client interface
type Client interface {
    RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error
    Close() error
}
```

```go
// pkg/storage/v2/database/mysql_client.go
package database

import (
    "context"

    "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
)

type mysqlClientAdapter struct {
    mc mysql.Client
}

func NewMySQLStorageClient(mc mysql.Client) Client {
    return &mysqlClientAdapter{mc: mc}
}

func (c *mysqlClientAdapter) RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error {
    return c.mc.RunInTransactionV2(ctx, func(ctx context.Context, _ mysql.Transaction) error {
        return f(ctx)
    })
}

func (c *mysqlClientAdapter) Close() error {
    return c.mc.Close()
}
```

```go
// pkg/storage/v2/database/postgres_client.go
package database

import (
    "context"

    "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
)

type postgresClientAdapter struct {
    pc postgres.Client
}

func NewPostgresStorageClient(pc postgres.Client) Client {
    return &postgresClientAdapter{pc: pc}
}

func (c *postgresClientAdapter) RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error {
    return c.pc.RunInTransactionV2(ctx, func(ctx context.Context, _ postgres.Transaction) error {
        return f(ctx)
    })
}

func (c *postgresClientAdapter) Close() error {
    return c.pc.Close()
}
```

#### 2. Storage Interface Pattern

Each storage package defines an interface with separate MySQL and PostgreSQL implementations:

```go
// pkg/push/storage/v2/push.go - Interface
type PushStorage interface {
    CreatePush(ctx context.Context, e *domain.Push, environmentId string) error
    UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error
    GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error)
    ListPushes(ctx context.Context, option *ListOptions) ([]*proto.Push, int, int64, error)
    DeletePush(ctx context.Context, id, environmentId string) error
}

// pkg/push/storage/v2/mysql_push.go - MySQL implementation (existing)
func NewMySQLPushStorage(qe mysql.QueryExecer) PushStorage

// pkg/push/storage/v2/postgres_push.go - PostgreSQL implementation (new)
func NewPostgresPushStorage(qe postgres.QueryExecer) PushStorage
```

#### 3. Shared Filter Types and Query Builders

##### Option A: Wrapper Pattern with Separate Packages

**Shared types in `pkg/storage/v2/database/filter.go`** (pure data structs, no `SQLString()` methods):

- `ListOptions`, `FilterV2`, `InFilter`, `NullFilter`, `JSONFilter`, `SearchQuery`, `OrFilter`, `Order`
- Constants: `Operator`, `JSONFilterFunc`, `OrderDirection`

**Database-specific wrappers** use composition pattern:

- `pkg/storage/v2/mysql/filter.go`: Embeds filter structs, adds `SQLString()` with `?` placeholders and MySQL JSON functions
- `pkg/storage/v2/postgres/filter.go`: Embeds filter structs, adds `SQLString()` with `$N` placeholders and PostgreSQL JSONB operators

**Example:**

```go
// pkg/storage/v2/database/filter.go - pure data
type FilterV2 struct {
    Column   string
    Operator Operator
    Value    interface{}
}

// pkg/storage/v2/mysql/filter.go - wraps and adds SQLString()
type FilterV2 struct {
    filter.FilterV2
}
func (f *FilterV2) SQLString() (string, []interface{}) {
    return fmt.Sprintf("%s %s ?", f.Column, f.Operator), []interface{}{f.Value}
}

// pkg/storage/v2/postgres/filter.go - wraps and adds SQLString() with index
type FilterV2 struct {
    filter.FilterV2
    index *int
}
func (f *FilterV2) SQLString() (string, []interface{}) {
    *f.index++
    return fmt.Sprintf("%s %s $%d", f.Column, f.Operator, *f.index), []interface{}{f.Value}
}
```

##### Option B: Unified Package with Placeholder Replacement

A single `database` package with one client implementation that internally handles database differences:

1. **Single client struct** with `dbType` field to track MySQL vs PostgreSQL
2. **Automatic placeholder replacement** - `?` → `$1, $2, ...` for PostgreSQL at query execution time
3. **Unified JSON filter** - Single `JSONFilter` type with `DBType` field that generates MySQL or PostgreSQL syntax
4. **Unified error handling** - Converts both MySQL and PostgreSQL errors to common types

**Structure:**

```
pkg/storage/v2/database/
├── client.go           # Single client with dbType field, NewClient/NewMySQLClient/NewPostgresClient
├── transaction.go      # Transaction with dbType for error conversion
├── query.go            # Filters, query builder, ReplacePlaceholders(), JSONFilter
├── result.go           # Result, Row, Rows interfaces
├── error.go            # Unified errors + MySQL/PostgreSQL error converters
├── metrics.go          # Database metrics
└── mock/               # Generated mocks
```

**Single client with dbType:**

```go
// pkg/storage/v2/database/client.go
package database

type DBType string

const (
    DBTypeMySQL    DBType = "mysql"
    DBTypePostgres DBType = "postgres"
)

type client struct {
    db     *sql.DB
    dbType DBType  // Determines placeholder replacement and error conversion
    opts   *options
    logger *zap.Logger
}

// NewClient creates a database client based on dbType
func NewClient(
    ctx context.Context,
    dbType DBType,
    dbUser, dbPass, dbHost string,
    dbPort int,
    dbName string,
    opts ...Option,
) (Client, error) {
    // Creates MySQL or PostgreSQL connection based on dbType
}

// Convenience functions
func NewMySQLClient(ctx, user, pass, host string, port int, dbName string, opts ...Option) (Client, error) {
    return NewClient(ctx, DBTypeMySQL, user, pass, host, port, dbName, opts...)
}

func NewPostgresClient(ctx, user, pass, host string, port int, dbName string, opts ...Option) (Client, error) {
    return NewClient(ctx, DBTypePostgres, user, pass, host, port, dbName, opts...)
}
```

**Automatic placeholder replacement in query methods:**

```go
func (c *client) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
    // Replace ? with $1, $2, ... for PostgreSQL
    if c.dbType == DBTypePostgres {
        query = ReplacePlaceholders(query)
    }

    // Execute query...
    err = c.convertError(err)  // Convert DB-specific errors
    return &result{sret}, err
}

func (c *client) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
    if c.dbType == DBTypePostgres {
        query = ReplacePlaceholders(query)
    }
    // ...
}

func (c *client) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
    if c.dbType == DBTypePostgres {
        query = ReplacePlaceholders(query)
    }
    // ...
}
```

**ReplacePlaceholders function:**

```go
// pkg/storage/v2/database/query.go

// ReplacePlaceholders converts ? placeholders to $1, $2, $3... for PostgreSQL
func ReplacePlaceholders(query string) string {
    var result strings.Builder
    paramIndex := 1
    for i := 0; i < len(query); i++ {
        if query[i] == '?' {
            result.WriteString(fmt.Sprintf("$%d", paramIndex))
            paramIndex++
        } else {
            result.WriteByte(query[i])
        }
    }
    return result.String()
}
```

**Unified JSON filter with DBType:**

```go
// JSONFilter - Single type that handles both MySQL and PostgreSQL
type JSONFilter struct {
    Column string
    Func   JSONFilterFunc
    Values []interface{}
    DBType DBType
}

func NewJSONFilter(column string, f JSONFilterFunc, values []interface{}, dbType DBType) WherePart {
    return &JSONFilter{Column: column, Func: f, Values: values, DBType: dbType}
}

func (f *JSONFilter) SQLString() (string, []interface{}) {
    switch f.Func {
    case JSONContainsString:
        if f.DBType == DBTypePostgres {
            return fmt.Sprintf("%s @> ?::jsonb", f.Column), []interface{}{formatJSONArray(f.Values, true)}
        }
        return fmt.Sprintf("JSON_CONTAINS(%s, ?)", f.Column), []interface{}{formatJSONArray(f.Values, true)}
    case JSONLengthGreaterThan:
        if f.DBType == DBTypePostgres {
            return fmt.Sprintf("jsonb_array_length(%s) > %v", f.Column, f.Values[0]), nil
        }
        return fmt.Sprintf("JSON_LENGTH(%s) > %v", f.Column, f.Values[0]), nil
    // ... other cases
    }
}
```

**Unified error conversion:**

```go
// pkg/storage/v2/database/error.go
var (
    ErrNoRows         = errors.New("database: no rows")
    ErrTxDone         = errors.New("database: tx done")
    ErrDuplicateEntry = errors.New("database: duplicate entry")
)

func (c *client) convertError(err error) error {
    if c.dbType == DBTypePostgres {
        return convertPostgresError(err)
    }
    return convertMySQLError(err)
}
```

##### Comparison

| Aspect                  | Option A (Wrapper Pattern)                | Option B (Unified Package)                        |
| ----------------------- | ----------------------------------------- | ------------------------------------------------- |
| Package structure       | Separate mysql/postgres/filter packages   | Single `database` package                         |
| Client type             | Adapter wrapping native clients           | Single client with `dbType` field                 |
| Filter types            | Wrapper types with embedding              | Shared types + Unified `JSONFilter` with `DBType` |
| Placeholder handling    | Generated at filter level with `$N`       | Auto-replaced at query execution                  |
| Storage implementations | Separate MySQL/PostgreSQL implementations | Single implementation for both                    |
| Code duplication        | More (separate implementations)           | Less (shared code)                                |
| Complexity              | Higher (more abstraction)                 | Lower (simpler)                                   |

**Output (both options):**

```sql
-- MySQL (uses ? placeholders natively)
SELECT * FROM push WHERE name = ? AND JSON_CONTAINS(tags, ?)

-- PostgreSQL (? auto-converted to $1, $2, ... at execution time)
SELECT * FROM push WHERE name = $1 AND tags @> $2::jsonb
```

### Implementation Strategy

#### Phase 1: Implement Storage Package and Query Builder

Option B - Unified `database` package:

1. Single `client` struct with `dbType` field (`DBTypeMySQL`, `DBTypePostgres`)
2. `NewClient(ctx, dbType, ...)`, `NewMySQLClient(...)`, `NewPostgresClient(...)`
3. Automatic placeholder replacement (`?` → `$1, $2, ...`) for PostgreSQL
4. Shared filter types (`FilterV2`, `InFilter`, `NullFilter`, `SearchQuery`, etc.)
5. Unified error handling (`ErrNoRows`, `ErrDuplicateEntry`, `ErrTxDone`)
6. Query builder functions (`ConstructQueryAndWhereArgs`, `ConstructCountQuery`)

#### Phase 2: Create PostgreSQL Schema Migration

Create PostgreSQL schema migrations for all tables to match the existing MySQL schema.

#### Phase 3: Implement Storage Layer and Refactor API Layer

**Storage Layer (Option A - Separate Implementations):**

```go
// pkg/push/storage/v2/push.go - Interface definition
type PushStorage interface {
    CreatePush(ctx context.Context, e *domain.Push, environmentId string) error
    UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error
    GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error)
    ListPushes(ctx context.Context, option *ListOptions) ([]*proto.Push, int, int64, error)
    DeletePush(ctx context.Context, id, environmentId string) error
}

// pkg/push/storage/v2/mysql_push.go - MySQL implementation
func NewMySQLPushStorage(qe mysql.QueryExecer) PushStorage

// pkg/push/storage/v2/postgres_push.go - PostgreSQL implementation
func NewPostgresPushStorage(qe postgres.QueryExecer) PushStorage
```

**Storage Layer (Option B - Single Implementation with Placeholder Replacement):**

With placeholder replacement at the client level, we can have a **single storage implementation**:

```go
// pkg/push/storage/v2/push.go - Single implementation for both MySQL and PostgreSQL
type PushStorage interface {
    CreatePush(ctx context.Context, e *domain.Push, environmentId string) error
    UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error
    GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error)
    ListPushes(ctx context.Context, option *database.ListOptions) ([]*proto.Push, int, int64, error)
    DeletePush(ctx context.Context, id, environmentId string) error
}

// Single constructor - works with both MySQL and PostgreSQL
func NewPushStorage(qe database.QueryExecer) PushStorage

// Usage:
// - All SQL uses ? placeholders
// - PostgreSQL client replaces ? → $1, $2, ... at execution time
// - For JSON filters, use JSONFilter with DBType parameter
```

**Example - Single implementation using ? placeholders:**

```go
func (s *pushStorage) CreatePush(ctx context.Context, push *domain.Push, envID string) error {
    // Same SQL works for both MySQL and PostgreSQL
    query := `INSERT INTO push (id, name, environment_id, created_at) VALUES (?, ?, ?, ?)`
    _, err := s.qe.ExecContext(ctx, query, push.Id, push.Name, envID, push.CreatedAt)
    return err
}

func (s *pushStorage) ListPushes(ctx context.Context, opts *database.ListOptions) ([]*proto.Push, error) {
    query := `SELECT * FROM push`
    // Query builder handles placeholder replacement internally
    query, args := database.ConstructQueryAndWhereArgs(query, opts)
    rows, err := s.qe.QueryContext(ctx, query, args...)
    // ...
}
```

**For JSON filters, use JSONFilter with DBType:**

```go
// JSONFilter generates appropriate SQL based on DBType
jsonFilter := database.NewJSONFilter("tags", database.JSONContainsString, []interface{}{"web"}, dbType)
// MySQL output: JSON_CONTAINS(tags, ?)
// PostgreSQL output: tags @> ?::jsonb
```

**API Layer:**
API services receive `database.Client` interface (not `mysql.Client`). Database client initialization happens in `server.go` only:

```go
// pkg/push/api/api.go
package api

import (
    "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database"
    v2ps "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2"
)

type PushService struct {
    dbClient     database.Client    // Unified database client (MySQL or PostgreSQL)
    pushStorage  v2ps.PushStorage   // Storage interface (MySQL or PostgreSQL implementation)
    // ... other fields
}

// NewPushService receives database.Client, not mysql.Client
// The caller (server.go) decides which database to use
func NewPushService(
    dbClient database.Client,        // Changed from mysql.Client
    pushStorage v2ps.PushStorage,    // Storage implementation injected
    // ... other params
    opts ...Option,
) *PushService {
    return &PushService{
        dbClient:    dbClient,
        pushStorage: pushStorage,
        // ...
    }
}

func (s *PushService) CreatePush(ctx context.Context, req *pushproto.CreatePushRequest) (*pushproto.CreatePushResponse, error) {
    // ... validation ...

    // No database-specific code - just use dbClient
    err = s.dbClient.RunInTransactionV2(ctx, func(ctx context.Context) error {
        if err := s.pushStorage.CreatePush(ctx, push, req.EnvironmentId); err != nil {
            return err
        }
        return nil
    })
    // ...
}
```

**Server Layer (initialization):**
Database selection and client creation happens in `server.go`:

```go
// pkg/web/cmd/server/server.go (Option A - Separate storage implementations)
func (s *server) createServices() {
    var dbClient database.Client
    var pushStorage v2ps.PushStorage

    switch s.config.StorageType {
    case "mysql":
        dbClient, _ = database.NewMySQLClient(ctx, user, pass, host, port, dbName)
        pushStorage = v2ps.NewMySQLPushStorage(dbClient)
    case "postgres":
        dbClient, _ = database.NewPostgresClient(ctx, user, pass, host, port, dbName)
        pushStorage = v2ps.NewPostgresPushStorage(dbClient)
    }

    pushService := pushapi.NewPushService(dbClient, pushStorage, ...)
}
```

```go
// pkg/web/cmd/server/server.go (Option B - Single storage implementation)
func (s *server) createServices() {
    var dbClient database.Client

    switch s.config.StorageType {
    case "mysql":
        dbClient, _ = database.NewMySQLClient(ctx, user, pass, host, port, dbName)
    case "postgres":
        dbClient, _ = database.NewPostgresClient(ctx, user, pass, host, port, dbName)
    }

    // Single storage implementation works with both databases
    // All SQL uses ? placeholders - PostgreSQL client auto-converts to $1, $2, ...
    pushStorage := v2ps.NewPushStorage(dbClient)

    pushService := pushapi.NewPushService(dbClient, pushStorage, ...)
}
```

### SQL Compatibility Considerations

#### Placeholder Syntax

- MySQL: `?` for all parameters
- PostgreSQL: `$1, $2, $3, ...`

#### JSON Operations

- MySQL: `JSON_CONTAINS()`
- PostgreSQL: `@>` JSONB operator

#### Auto-increment vs Serial

- MySQL: `AUTO_INCREMENT`
- PostgreSQL: `SERIAL` or `GENERATED ALWAYS AS IDENTITY`

This is handled at the schema level, not in application code.

### Affected Packages

Based on analysis, the following packages need refactoring:

- `pkg/account/storage/v2`
- `pkg/feature/storage/v2`
- `pkg/experiment/storage/v2`
- `pkg/environment/storage/v2`
- `pkg/push/storage/v2`
- `pkg/notification/storage/v2`
- `pkg/autoops/storage/v2`
- `pkg/auditlog/storage/v2`
- `pkg/tag/storage`
- `pkg/team/storage`
- `pkg/mau/storage`
- `pkg/opsevent/storage/v2`
- `pkg/coderef/storage`
- `pkg/subscriber/storage/v2`
- `pkg/experimentcalculator/storage/v2`
- `pkg/eventcounter/storage/v2`

### Testing Strategy

1. **Unit Tests API layer**: refactor API layer to use unit test for both MySQL and PostgreSQL implementations
2. **Unit Tests postgresQL builder**: Create test cases for PostgreSQL query builder functions to ensure correct placeholder generation and SQL syntax
3. **E2E Tests**: Existing E2E tests should pass with either database backend

## Trade-offs

### Advantages

1. **Flexibility**: Users can choose their preferred database
2. **Open Source Friendly**: PostgreSQL is fully open source
3. **Cost Reduction**: No BigQuery dependency for analytics when using PostgreSQL with TimescaleDB
4. **Simplified Operations**: Single database system for both OLTP and OLAP

### Disadvantages

1. **Increased Complexity**: More code to maintain (two implementations)
2. **Testing Overhead**: Need to test against both databases
3. **SQL Differences**: Some queries may need database-specific versions

## Implementation Timeline

| Phase | Description                                                     | Estimated Effort |
| ----- | --------------------------------------------------------------- | ---------------- |
| 1     | Implement storage package postgres and query builder            | 1-2 weeks        |
| 2     | Create PostgreSQL schema migration                              | 1 week           |
| 3     | Implement storage layer for each package and refactor API layer | 4-5 weeks        |
