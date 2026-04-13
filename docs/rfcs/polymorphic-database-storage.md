# RFC: Polymorphic Database Storage Layer

## I. Summary

This RFC proposes PostgreSQL as an alternative primary storage backend alongside MySQL, selected at runtime via configuration.

Key design principles:

1. **Unified transaction interface** — A narrow `database.Client` (`RunInTransactionV2`, `Close`) abstracts transaction handling so orchestration code does not depend on `mysql.Client` / `postgres.Client` for txs. Adapters live in `pkg/storage/v2/database` (e.g. `NewMySQLStorageClient`, `NewPostgresStorageClient`).
2. **Semantic storage interfaces** — The API passes domain structs and enums (not `mysql.ListOptions` or SQL fragments). Each storage implementation maps those inputs to dialect-specific list options and SQL.
3. **Split query construction by dialect** — `pkg/storage/v2/mysql` and `pkg/storage/v2/postgres` each own placeholders (`?` vs `$n`), JSON predicates, and helpers. Domain storage ships per-dialect files (e.g. `*_mysql.go`, `*_postgres.go`) or a factory chosen at server startup.

## II. Background and motivation

Bucketeer services are tightly coupled to MySQL for primary storage:

1. **API layer holds MySQL clients and types**  
   Services such as `PushService` use `mysql.Client` directly, and list endpoints often build `mysql.ListOptions`, `mysql.FilterV2`, `mysql.JSONFilter`, and `mysql.Order` in the API package.

2. **Storage implementations depend on MySQL**  
   Storage structs use `mysql.QueryExecer`; interfaces expose methods that take `*mysql.ListOptions`, which pins every caller and mock to MySQL.

3. **MySQL-specific concerns leak upward**  
   JSON serialization, duplicate-key and no-rows errors, placeholder style (`?` vs `$1`), and query-builder types are visible outside storage.

PostgreSQL is already used for the data warehouse (eventcounter). Primary OLTP storage needs a clearer boundary: **transactions** can be abstracted once; **SQL shape** cannot be fully hidden without either fragile shortcuts or repeated `if postgres` branches.

### Why this change is needed

A major blocker for polymorphic storage is not only `mysql.Client` on services but **list and filter logic living in the API layer** using types that are really **query builder fragments**.

**Symptoms**

- The API imports `pkg/storage/v2/mysql` to construct `[]*mysql.FilterV2`, `[]*mysql.JSONFilter`, `*mysql.SearchQuery`, `[]*mysql.Order`, and `*mysql.ListOptions`.
- **`FilterV2.Column` is sometimes a full SQL expression**, not a column name. For example, listing features with “has feature flag as rule” uses MySQL-specific JSON functions in the column slot:

  ```go
  filters = append(filters, &mysql.FilterV2{
      Column:   "JSON_CONTAINS(JSON_EXTRACT(rules, '$[*].clauses[*].operator'), '11')",
      Operator: mysql.OperatorEqual,
      Value:    true,
  })
  ```

  PostgreSQL would need a different predicate (`jsonb_path_exists`, `@>`, etc.). No amount of `?` → `$1` rewriting fixes that.

- **Ordering** is mapped in the API with names like `newListFeaturesOrdersMySQL`, producing SQL-oriented values (qualified columns and expressions such as `(progressive_rollout_count + schedule_count + kill_switch_count)`).
- **Storage interfaces** expose `ListFeatures(ctx, *mysql.ListOptions)`, so the API must speak MySQL to call storage—this contradicts the goal of database selection by configuration only.

**Why this matters**

- Polymorphic backends require **dialect-specific predicates** for JSON, upsert, and some aggregates. If the API continues to assemble `mysql.ListOptions`, every new endpoint duplicates that coupling.
- The goal “no database-specific conditionals in the API layer” is violated as long as the API builds MySQL filters and embeds MySQL SQL in struct fields meant for columns.

## III. Goals

1. Enable PostgreSQL as an alternative primary storage backend alongside MySQL.
2. Select the database via configuration; **API business logic does not branch on dialect** and does not import mysql/postgres query types for lists.
3. Maintain backward compatibility for existing MySQL deployments.
4. Keep a small unified surface for **transactions** (`database.Client`) where it reduces boilerplate.
5. Accept **duplication across MySQL and PostgreSQL storage paths** where it keeps SQL explicit and avoids a single mega-implementation full of `if dbType`.

## IV. Design overview

### Approach

Treat list/filter **intent** as part of the **storage** contract: the API validates auth and request shape, then calls storage with **semantic parameters** (environment ID, tag list, enums for order and status, optional booleans, pagination). Each storage implementation maps those parameters to its own list options and SQL. This is how we address the coupling described in §II without pushing dialect-specific query building back into handlers.

### Architecture layers

```mermaid
flowchart TB
  subgraph API_Layer["API layer"]
    direction TB
    T1["Transactions: database.Client (dialect-agnostic)"]
    T2["CRUD / lists: storage interfaces, semantic params only"]
    T3["No mysql.ListOptions or SQL fragments in handlers"]
  end

  API_Layer --> Client["database.Client + adapters"]
  API_Layer --> Impl["Storage implementation (chosen at server startup)"]

  Impl --> M["mysql_*.go (uses mysql package)"]
  Impl --> P["postgres_*.go (uses postgres package)"]

  M --> MP["pkg/storage/v2/mysql: query builders, filters, ? placeholders"]
  P --> PP["pkg/storage/v2/postgres: query builders, filters, $n placeholders"]
```

### 1. Unified database client (transactions)

`database.Client` wraps the existing MySQL and PostgreSQL clients for **`RunInTransactionV2` and `Close` only**. Storage continues to execute queries through a **`QueryExecer`** that respects the transactional context (same pattern as today: transaction attached to `context`, client methods dispatch to the active tx).

```go
// pkg/storage/v2/database/client.go
package database

import "context"

type Client interface {
    RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error
    Close() error
}
```

Adapters delegate to `mysql.Client` / `postgres.Client` and strip the unused `Transaction` argument from the callback while preserving context-based tx propagation inside the underlying client.

### 2. Split query packages and storage implementations

**MySQL** and **PostgreSQL** each keep their own:

- Placeholder style (`?` vs numbered `$n`).
- Filter / list helpers that implement a shared **conceptual** model (equality, IN, NULL, JSON contains, search, ORDER BY) with **different emitted SQL** (e.g. mysql `?` vs postgres numbered `$n` via `BindSQL` and composition helpers).
- Any dialect-specific upsert, JSON, or locking syntax.

Domain storage (e.g. `pkg/feature/storage/v2`) exposes:

- A **storage interface** that uses **semantic list parameters**, not `*mysql.ListOptions`.
- **`NewMySQLFeatureStorage(qe mysql.QueryExecer) FeatureStorage`** and **`NewPostgresFeatureStorage(qe postgres.QueryExecer) FeatureStorage`** (names illustrative), or a factory in `server` that picks one.

Simple CRUD can remain one file per dialect or share scanning helpers; **list methods** translate semantic params → `mysql.ListOptions` or `postgres.ListOptions` **inside** the implementation file for that backend.

### 3. List query semantics live in storage

**Example direction for feature listing:**

```go
// pkg/feature/storage/v2 — semantic, no mysql import required for callers

type ListFeaturesParams struct {
    EnvironmentID         string
    Tags                  []string
    Maintainer            string
    Enabled               *bool
    Archived              *bool
    HasPrerequisites      *bool
    HasFeatureFlagAsRule  *bool
    SearchKeyword         string
    Status                featureproto.FeatureLastUsedInfo_Status
    OrderBy               featureproto.ListFeaturesRequest_OrderBy
    OrderDirection        featureproto.ListFeaturesRequest_OrderDirection
    Limit                 int
    Offset                int
}

type FeatureStorage interface {
    ListFeatures(ctx context.Context, p ListFeaturesParams) ([]*featureproto.Feature, int, int64, error)
    ListFeaturesFilteredByExperiment(ctx context.Context, p ListFeaturesParams) ([]*featureproto.Feature, int, int64, error)
    // ... other methods without *mysql.ListOptions in the public interface
}
```

- **`mysql` implementation**: maps `p` to `*mysql.ListOptions`, including translating `HasFeatureFlagAsRule` into `JSON_CONTAINS(JSON_EXTRACT(...))` (or a dedicated helper in `pkg/storage/v2/mysql`).
- **`postgres` implementation**: maps the same `p` into postgres list options with JSONB predicates that express the same intent.

The API calls `ListFeatures(ctx, ListFeaturesParams{...})` after unmarshalling and permission checks; it does not reference `mysql.JSONFilter` or SQL strings.

### 4. Shared neutral types (optional, limited)

Where several domains need the same **pure data** shape (e.g. operator enum, order direction as int), those can live in `pkg/storage/v2/database` **without** query-builder methods. Dialect packages embed or convert them. This is optional; avoiding a heavy shared query DSL is fine as long as **SQL emission stays in mysql/postgres or in `*_mysql.go` / `*_postgres.go`**.

## V. Alternative considered: unified client and `?` → `$N` replacement (Option B)

We considered a **single** `database` package with one `QueryExecer`, a `dbType` field, **automatic replacement** of `?` with `$1`, `$2`, … at execution time (similar in spirit to [sqlx `Rebind`](https://github.com/jmoiron/sqlx/blob/master/bind.go)), unified error translation, and optionally a single storage implementation for both databases.

### Problems with that approach

1. **Placeholder rewriting is not SQL-aware**  
   Typical implementations scan the string for `?` and replace in order. That breaks if `?` appears inside a string literal or comment; sqlx documents this limitation explicitly. Safer approaches either restrict all SQL to fixed templates or emit `$n` at build time (as GORM does per dialect via `BindVarTo`).

2. **Dialect differences are not only placeholders**  
   JSON (`JSON_CONTAINS` vs `@>` / `jsonb_path_exists`), upsert (`ON DUPLICATE KEY UPDATE` vs `ON CONFLICT ... DO UPDATE`), some aggregates and full-text APIs differ. A unified implementation either **branches on `dbType` everywhere** (collapsing into one file of conditionals) or still needs separate SQL strings per dialect—so the “single implementation” win is small.

3. **API-embedded SQL fragments make unified SQL untenable**  
   Even with perfect placeholder replacement, expressions like `JSON_CONTAINS(JSON_EXTRACT(rules, ...))` in filter “columns” must become postgres-specific text. That logic cannot live in a generic `ReplacePlaceholders` pass; it belongs next to the backend that understands the schema.

4. **Transaction path must rewrite consistently**  
   Every execution path (`ExecContext`, `QueryContext`, on pooled `DB` and on `Tx`) must apply the same rules; a single bug duplicates subtle production errors.

### Why we chose split query + split storage paths

- **Explicit SQL per dialect** is easier to review and test than a large `if dbType` matrix.
- **Placeholders** are correct by construction in each package (`?` vs `$n` indexing), matching how ORMs like GORM generate SQL.
- **List intent** in the storage API prevents the feature (and other) API packages from becoming query builders.

We keep the **unified client only where the abstraction is honest**: transactions and startup wiring, not “one string fits all databases.”

### Trade-offs (summary)

**Advantages**

- Clear ownership: **API** = auth, validation, orchestration; **storage** = persistence and SQL shape.
- PostgreSQL support without smuggling MySQL expressions through generic filter fields.
- Easier code review than a single dialect-agnostic query interpreter.

**Disadvantages**

- More files and some duplicated mapping logic between `*_mysql.go` and `*_postgres.go`.
- Broader test matrix (both databases for behavior that differs).

## VI. SQL compatibility (reference)

### Placeholders

- MySQL: `?`
- PostgreSQL: `$1`, `$2`, …

### JSON

- MySQL: `JSON_CONTAINS`, `JSON_EXTRACT`, etc.
- PostgreSQL: JSONB operators and functions (`@>`, `jsonb_path_exists`, `jsonb_array_length`, …)

### Upsert

| Aspect          | MySQL                          | PostgreSQL                             |
| --------------- | ------------------------------ | -------------------------------------- |
| Clause          | `ON DUPLICATE KEY UPDATE`      | `ON CONFLICT (...) DO UPDATE SET`      |
| Conflict target | Implicit from unique/PK        | Must match a unique index / constraint |
| Updated values  | `VALUES(col)` (legacy) / alias | `EXCLUDED.col`                         |

Per-dialect helpers (in `mysql` / `postgres` packages or next to the storage impl) should generate the appropriate clause; PostgreSQL requires schema-level unique constraints that align with the chosen conflict target.

### Auto-increment vs serial

Handled in migrations (MySQL `AUTO_INCREMENT`, PostgreSQL `SERIAL` / `IDENTITY`). Application code may need `RETURNING` for last-insert id on PostgreSQL where MySQL used `LastInsertId`.

### Transaction propagation

Today, transactional execution uses **context** (transaction value on context) and a **client** that implements `QueryExecer`. Any unified `database.Client` used only for `RunInTransactionV2` must not break that pattern: storage must still run queries through a `QueryExecer` that participates in the same transaction.

## VII. Implementation details

This section combines **scope** (which packages move), a **vertical walkthrough** (how code is wired), and **testing** expectations.

### Affected packages

Refactoring touches storage and APIs that currently import mysql list types or construct filters in the API layer, including but not limited to:

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

### Examples (client, filters, API, storage)

The walkthrough uses feature listing as a reference, in wiring order: **client → postgres list helpers → API → storage**, with dialect-specific query construction.

#### Unified database client and transactions

The API depends on `database.Client` only to run work inside a transaction. Adapters wrap the existing `mysql.Client` / `postgres.Client` and forward `RunInTransactionV2` while keeping **transaction propagation on `context`** (the underlying client attaches the active `Transaction` to context; `ExecContext` / `QueryContext` on that client use the tx when present).

```go
// pkg/storage/v2/database — narrow interface
type Client interface {
    RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error
    Close() error
}

// Adapter (see mysqlStorageClientAdapter in pkg/storage/v2/database): delegate to mysql.Client, drop unused tx param in callback
func (a *mysqlStorageClientAdapter) RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error {
    return a.inner.RunInTransactionV2(ctx, func(ctx context.Context, _ mysql.Transaction) error {
        return f(ctx)
    })
}
```

**Storage** still receives a **`QueryExecer`** (`mysql.Client` or `postgres.Client` implementing `ExecContext` / `QueryContext` / `QueryRowContext`) so queries run on the pool or on the transactional connection via context. The unified `database.Client` does not replace `QueryExecer` for CRUD unless you explicitly design a facade that implements both.

#### PostgreSQL filter / list package

`pkg/storage/v2/postgres` grows types analogous to `mysql` (e.g. `FilterV2`, `ListOptions`, `JSONFilter`, `Order`). WHERE fragments implement **`BindSQL(next int) (sql string, args []interface{}, nextAfter int)`** so placeholders are **`$1`, `$2`, …** by construction; `ConstructWhereSQLString` chains parts with a running index. JSON helpers use **JSONB** (`@>`, `jsonb_array_length`, …). Helpers such as `WritePlaceHolder` can format grouped placeholders (e.g. `IN ($1,$2)`).

Illustrative shape (not the full API):

```go
// pkg/storage/v2/postgres — each part knows how many bind slots it consumes
func (f *FilterV2) BindSQL(next int) (string, []interface{}, int) {
    sql := fmt.Sprintf("%s = $%d", f.Column, next)
    return sql, []interface{}{f.Value}, next + 1
}

// JSON contains — contrast with mysql.JSON_CONTAINS
func (j *JSONFilter) BindSQL(next int) (string, []interface{}, int) {
    return fmt.Sprintf("(%s::jsonb @> $%d::jsonb)", j.Column, next), []interface{}{formatJSONBArr(j.Values)}, next + 1
}
```

List assembly mirrors `mysql.ConstructQueryAndWhereArgs` / count helpers but never uses `?`.

#### API layer changes

**Transactions only:** the API layer replaces **`mysql.Client`** with **`database.Client`** for `RunInTransactionV2` / `Close`. It does **not** open connections or choose a dialect—that stays in `server`. Storage implementations still receive a concrete **`mysql.Client`** or **`postgres.Client`** as `QueryExecer` where needed, but **service structs** depend on the narrow `database.Client` interface plus **`PushStorage` / `FeatureStorage`** (etc.).

**Before:**

```go
type PushService struct {
    mysqlClient mysql.Client
    pushStorage v2ps.PushStorage
}

func NewPushService(mysqlClient mysql.Client, pushStorage v2ps.PushStorage, ...) *PushService {
    return &PushService{mysqlClient: mysqlClient, pushStorage: pushStorage, ...}
}

func (s *PushService) CreatePush(...) error {
    return s.mysqlClient.RunInTransactionV2(ctx, func(ctx context.Context, _ mysql.Transaction) error {
        return s.pushStorage.CreatePush(ctx, ...)
    })
}
```

**After:**

```go
type PushService struct {
    dbClient    database.Client
    pushStorage v2ps.PushStorage
}

func NewPushService(dbClient database.Client, pushStorage v2ps.PushStorage, ...) *PushService {
    return &PushService{dbClient: dbClient, pushStorage: pushStorage, ...}
}

func (s *PushService) CreatePush(...) error {
    return s.dbClient.RunInTransactionV2(ctx, func(ctx context.Context) error {
        return s.pushStorage.CreatePush(ctx, ...)
    })
}
```

**List / filter paths** stay behind storage: handlers pass **semantic params** (e.g. `ListFeaturesParams`), not `mysql.ListOptions`.

```go
params := v2fs.ListFeaturesParams{
    EnvironmentID: req.EnvironmentId,
    Tags:          req.Tags,
    // ... enums, optional bools, pagination from req
}
features, cursor, total, err := s.featureStorage.ListFeatures(ctx, params)
```

##### Wiring in `pkg/web/cmd/server` (config-driven clients)

Only the server (or equivalent composition root) reads **`StorageType`** (or env / flags), constructs **`mysql.Client`** or **`postgres.Client`**, wraps the transactional surface as **`database.Client`**, and picks **`NewMySQLPushStorage`** vs **`NewPostgresPushStorage`**. Adapters such as **`database.NewMySQLStorageClient` / `database.NewPostgresStorageClient`** live in `pkg/storage/v2/database` (names illustrative).

```go
func (s *server) createServices() {
    ctx := context.Background() // or server-held root context
    var dbClient database.Client
    var pushStorage v2ps.PushStorage

    switch s.config.StorageType {
    case "mysql":
        mysqlClient, _ := mysql.NewClient(ctx, user, pass, host, port, dbName)
        dbClient = database.NewMySQLStorageClient(mysqlClient)
        pushStorage = v2ps.NewMySQLPushStorage(mysqlClient)
    case "postgres":
        pgClient, _ := postgres.NewClient(ctx, user, pass, host, port, dbName)
        dbClient = database.NewPostgresStorageClient(pgClient)
        pushStorage = v2ps.NewPostgresPushStorage(pgClient)
    }

    pushService := pushapi.NewPushService(dbClient, pushStorage, ...)
    _ = pushService
}
```

`PushService` never imports `mysql` or `postgres`; it only sees `database.Client` and `PushStorage`. The same pattern applies to `FeatureService`, `TagService`, and other APIs after their storage constructors are split by dialect.

#### Storage layer: queries in mysql and postgres packages

Each dialect implementation maps **`ListFeaturesParams` → native list options → SQL**.

**MySQL** (`feature_mysql.go` or equivalent) builds `*mysql.ListOptions`, then uses existing helpers against embedded SQL templates:

```go
func (s *featureStorageMySQL) ListFeatures(ctx context.Context, p ListFeaturesParams) ([]*featureproto.Feature, int, int64, error) {
    opts := buildMySQLListOptions(p) // FilterV2, JSONFilter, Order — may include JSON_CONTAINS / JSON_EXTRACT for rule filters
    query, args := mysql.ConstructQueryAndWhereArgs(selectFeaturesSQLQuery, opts)
    rows, err := s.qe.QueryContext(ctx, query, args...)
    // scan...
}
```

**PostgreSQL** (`feature_postgres.go`) builds `*postgres.ListOptions` with the same _intent_ but different `BindSQL` / composed SQL output, then uses postgres `ConstructQueryAndWhereArgs` (or equivalent):

```go
func (s *featureStoragePostgres) ListFeatures(ctx context.Context, p ListFeaturesParams) ([]*featureproto.Feature, int, int64, error) {
    opts := buildPostgresListOptions(p) // same p, different predicates e.g. jsonb_path_exists / @> for rules
    query, args := postgres.ConstructQueryAndWhereArgs(selectFeaturesSQLQuery, opts)
    rows, err := s.qe.QueryContext(ctx, query, args...)
    // scan identical proto rows
}
```

The **embedded base SELECT** (`select_features.sql`) can stay one file per dialect if fragments differ, or share a static `SELECT ... FROM feature ...` string with dialect-specific `WHERE` composition only—either way, **the strings that differ (JSON, upsert)** live next to `mysql` or `postgres` helpers, not in `pkg/feature/api`.

### Testing strategy

1. **Storage unit tests** per dialect for list-parameter mapping and query construction (including JSON and edge filters).
2. **API unit tests** use the `FeatureStorage` (etc.) interface with mocks that accept **semantic params**, not `*mysql.ListOptions`.
3. **Integration / E2E** run against MySQL and, as coverage grows, PostgreSQL.

## VIII. Implementation timeline (indicative)

| Phase | Description                                                                                                                                                                                                                                                      |
| ----- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1     | **Helm + Atlas prep:** chart values / config for PostgreSQL (DSN, secrets, `StorageType`); add an **Atlas** migration **folder** (or dialect target) for PostgreSQL next to the existing MySQL migration layout (e.g. under `migration/`, migration job charts). |
| 2     | Harden `pkg/storage/v2/postgres` query/list helpers to parity with mysql where needed; unified `database.Client` adapters if not already done.                                                                                                                   |
| 3     | Introduce semantic list params and refactor **feature** storage + API as the reference pattern; remove mysql filter construction from `pkg/feature/api`.                                                                                                         |
| 4     | Repeat storage split + API decoupling for remaining packages in the affected list.                                                                                                                                                                               |

---

## IX. Appendix: Example SQL outputs (illustrative)

```sql
-- MySQL list fragment
SELECT * FROM feature WHERE environment_id = ? AND JSON_CONTAINS(tags, ?)

-- PostgreSQL equivalent intent
SELECT * FROM feature WHERE environment_id = $1 AND tags @> $2::jsonb
```

These strings are built inside dialect-specific code paths, not in the API layer.
