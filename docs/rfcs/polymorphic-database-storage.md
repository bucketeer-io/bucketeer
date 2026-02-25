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

**Shared types in `pkg/storage/v2/filter.go`** (pure data structs, no `SQLString()` methods):

- `ListOptions`, `FilterV2`, `InFilter`, `NullFilter`, `JSONFilter`, `SearchQuery`, `OrFilter`, `Order`
- Constants: `Operator`, `JSONFilterFunc`, `OrderDirection`

**Database-specific wrappers** use composition pattern:

- `pkg/storage/v2/mysql/filter.go`: Embeds v2 structs, adds `SQLString()` with `?` placeholders and MySQL JSON functions
- `pkg/storage/v2/postgres/filter.go`: Embeds v2 structs, adds `SQLString()` with `$N` placeholders and PostgreSQL JSONB operators

**Example:**

```go
// pkg/storage/v2/filter.go - pure data
type FilterV2 struct {
    Column   string
    Operator Operator
    Value    interface{}
}

// pkg/storage/v2/mysql/filter.go - wraps and adds SQLString()
type FilterV2 struct {
    v2.FilterV2
}
func (f *FilterV2) SQLString() (string, []interface{}) {
    return fmt.Sprintf("%s %s ?", f.Column, f.Operator), []interface{}{f.Value}
}

// pkg/storage/v2/postgres/filter.go - wraps and adds SQLString() with index
type FilterV2 struct {
    v2.FilterV2
    index *int
}
func (f *FilterV2) SQLString() (string, []interface{}) {
    *f.index++
    return fmt.Sprintf("%s %s $%d", f.Column, f.Operator, *f.index), []interface{}{f.Value}
}
```

**Output:**

```sql
-- MySQL
SELECT * FROM push WHERE name = ? AND JSON_CONTAINS(tags, ?)

-- PostgreSQL
SELECT * FROM push WHERE name = $1 AND tags @> $2::jsonb
```

### Implementation Strategy

#### Phase 1: Implement Storage Package and Query Builder

1. Create shared filter types in `pkg/storage/v2/filter.go` (pure data structs)
2. Create MySQL filter wrappers in `pkg/storage/v2/mysql/filter.go`
3. Create PostgreSQL filter wrappers in `pkg/storage/v2/postgres/filter.go`
4. Implement PostgreSQL query builder in `pkg/storage/v2/postgres/query.go`
5. Unified Client interface

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
    "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database"
    "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
    "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
    v2ps "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2"
)

type PushService struct {
    dbClient     database.Client    // Unified database client for transactions
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

    var dbClient database.Client
    var pushStorage v2ps.PushStorage

    switch dopts.storageConfig.Type {
    case "mysql":
        dbClient = database.NewMySQLStorageClient(mysqlClient)
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
        dbClient = database.NewPostgresStorageClient(pgClient)
        pushStorage = v2ps.NewPostgresPushStorage(pgClient)
    default:
        dbClient = database.NewMySQLStorageClient(mysqlClient)
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

| Phase | Description                                                     | Estimated Effort |
| ----- | --------------------------------------------------------------- | ---------------- |
| 1     | Implement storage package postgres and query builder            | 1-2 weeks        |
| 2     | Create PostgreSQL schema migration                              | 1 week           |
| 3     | Implement storage layer for each package and refactor API layer | 4-5 weeks        |
