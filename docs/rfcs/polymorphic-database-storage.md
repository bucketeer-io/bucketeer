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

func NewMySQLClient(mc mysql.Client) Client {
    return &mysqlClientAdapter{mc: mc}
}

func (c *mysqlClientAdapter) RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error {
    // mysql transaction implementation
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

func NewPostgresClient(pc postgres.Client) Client {
    return &postgresClientAdapter{pc: pc}
}

func (c *postgresClientAdapter) RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error {
   // postgres transaction implementation
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

#### 2. SQL Query Helper Functions

Each database package provides its own query helper functions (no shared interface needed):

**MySQL** (existing in `pkg/storage/v2/mysql/query.go`):
```go
// Uses ? placeholders
func ConstructQueryAndWhereArgs(baseQuery string, options *ListOptions) (string, []interface{})
func ConstructCountQuery(baseQuery string, options *ListOptions) (string, []interface{})
func ConstructWhereSQLString(wps []WherePart) (string, []interface{})
```

**PostgreSQL** (to be implemented in `pkg/storage/v2/postgres/query.go`):
```go
// Uses $1, $2, $3... placeholders
func ConstructQueryAndWhereArgs(baseQuery string, options *ListOptions) (string, []interface{})
func ConstructCountQuery(baseQuery string, options *ListOptions) (string, []interface{})
func ConstructWhereSQLString(wps []WherePart) (string, []interface{})

// Already exists:
func WritePlaceHolder(template string, start, count int) string
```

**WritePlaceHolder Examples:**

| Input | Output |
|-------|--------|
| `WritePlaceHolder("($%d, $%d)", 1, 2)` | `"($1, $2)"` |
| `WritePlaceHolder("($%d, $%d, $%d)", 1, 3)` | `"($1, $2, $3)"` |
| `WritePlaceHolder("($%d, $%d, $%d)", 4, 3)` | `"($4, $5, $6)"` |

**Placeholder Syntax Difference:**
```sql
-- MySQL uses ? for all parameters
INSERT INTO push (id, name, tags) VALUES (?, ?, ?)

-- PostgreSQL uses $1, $2, $3...
INSERT INTO push (id, name, tags) VALUES ($1, $2, $3)
```

### Implementation Strategy

#### Phase 1: Implement Storage Package and Query Builder

Implement PostgreSQL query builder functions in `pkg/storage/v2/postgres/`:

```go
// pkg/storage/v2/postgres/query.go
func ConstructQueryAndWhereArgs(baseQuery string, options *ListOptions) (string, []interface{})
func ConstructCountQuery(baseQuery string, options *ListOptions) (string, []interface{})
func ConstructWhereSQLString(wps []WherePart) (string, []interface{})

// Already exists:
func WritePlaceHolder(template string, start, count int) string
```

Also implement the unified `Client` interface in `pkg/storage/v2/`:

```go
// pkg/storage/v2/client.go - Already implemented
type Client interface {
    RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error
    Close() error
}

// pkg/storage/v2/mysql_client.go - Already implemented
func NewMySQLClient(mc mysql.Client) Client

// pkg/storage/v2/postgres_client.go - Already implemented
func NewPostgresClient(pc postgres.Client) Client
```

#### Phase 2: Create PostgreSQL Schema Migration

Create PostgreSQL schema migrations for all tables to match the existing MySQL schema.

#### Phase 3: Implement Storage Layer and Refactor API Layer

For each storage package, create PostgreSQL implementations and update the API layer:

**Storage Layer:**
```go
// pkg/push/storage/v2/push.go - Interface definition
type PushStorage interface {
    CreatePush(ctx context.Context, e *domain.Push, environmentId string) error
    UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error
    GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error)
    ListPushes(ctx context.Context, option *ListOptions) ([]*proto.Push, int, int64, error)
    DeletePush(ctx context.Context, id, environmentId string) error
}

// pkg/push/storage/v2/mysql_push.go - MySQL implementation (existing, rename from push.go)
func NewMySQLPushStorage(qe mysql.QueryExecer) PushStorage

// pkg/push/storage/v2/postgres_push.go - PostgreSQL implementation (new)
func NewPostgresPushStorage(qe postgres.QueryExecer) PushStorage
```

**API Layer:**
Update API services to use `v2.Client` interface:

```go
// pkg/push/api/api.go
package api

import (
    storage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2"
    "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
    "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
    v2ps "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2"
)

type PushService struct {
    dbClient     storage.Client    // Unified database client for transactions
    pushStorage  v2ps.PushStorage
    // ... other fields
}

func NewPushService(
    mysqlClient mysql.Client,
    // ... other params
    opts ...Option,
) *PushService {
    dopts := &options{
        logger: zap.NewNop(),
        storageConfig: &StorageConfig{
            Type: "mysql", // default
        },
    }
    for _, opt := range opts {
        opt(dopts)
    }

    var dbClient storage.Client
    var pushStorage v2ps.PushStorage

    switch dopts.storageConfig.Type {
    case "mysql":
        dbClient = storage.NewMySQLClient(mysqlClient)
        pushStorage = v2ps.NewMySQLPushStorage(mysqlClient)
    case "postgres":
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        pgClient, err := createPostgresClient(ctx, dopts.storageConfig.Postgres, dopts.logger)
        if err != nil {
            dopts.logger.Error("Failed to create Postgres client",
                zap.Error(err),
                zap.String("host", dopts.storageConfig.Postgres.Host),
            )
            return nil
        }
        dbClient = storage.NewPostgresClient(pgClient)
        pushStorage = v2ps.NewPostgresPushStorage(pgClient)
    default:
        dbClient = storage.NewMySQLClient(mysqlClient)
        pushStorage = v2ps.NewMySQLPushStorage(mysqlClient)
    }

    return &PushService{
        dbClient:    dbClient,
        pushStorage: pushStorage,
        // ...
    }
}

func (s *PushService) CreatePush(ctx context.Context, req *pushproto.CreatePushRequest) (*pushproto.CreatePushResponse, error) {
    // ... validation ...
    
    var event *eventproto.Event
    
    // No if-else needed - just use dbClient
    err = s.dbClient.RunInTransactionV2(ctx, func(ctx context.Context) error {
        if err := s.pushStorage.CreatePush(ctx, push, req.EnvironmentId); err != nil {
            return err
        }
        // ... domain event handling ...
        return nil
    })
    // ...
}
```

### SQL Compatibility Considerations

#### Placeholder Syntax
- MySQL: `?` for all parameters
- PostgreSQL: `$1, $2, $3, ...` (use `WritePlaceHolder` function)

#### JSON Operations
Both MySQL and PostgreSQL support JSON, but with slightly different syntax:
- MySQL: `JSON_EXTRACT()`, `JSON_SET()`
- PostgreSQL: `->`, `->>`, `jsonb` operators

For most cases, storing JSON as text and handling serialization in Go is sufficient.

#### Auto-increment vs Serial
- MySQL: `AUTO_INCREMENT`
- PostgreSQL: `SERIAL` or `GENERATED ALWAYS AS IDENTITY`

This is handled at the schema level, not in application code.

### Affected Packages

Based on analysis, the following packages need refactoring:

- [ ] `pkg/account/storage/v2`
- [ ] `pkg/feature/storage/v2`
- [ ] `pkg/experiment/storage/v2`
- [ ] `pkg/environment/storage/v2`
- [ ] `pkg/push/storage/v2`
- [ ] `pkg/notification/storage/v2`
- [ ] `pkg/autoops/storage/v2`
- [ ] `pkg/auditlog/storage/v2`
- [ ] `pkg/tag/storage`
- [ ] `pkg/team/storage`
- [ ] `pkg/mau/storage`
- [ ] `pkg/opsevent/storage/v2`
- [ ] `pkg/coderef/storage`
- [ ] `pkg/subscriber/storage/v2`
- [ ] `pkg/experimentcalculator/storage/v2`
- [ ] `pkg/eventcounter/storage/v2`

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

| Phase | Description | Estimated Effort |
|-------|-------------|------------------|
| 1 | Implement storage package postgres and query builder | 1-2 weeks |
| 2 | Create PostgreSQL schema migration | 1 week |
| 3 | Implement storage layer for each package and refactor API layer | 4-5 weeks |
