# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is Bucketeer

Bucketeer is an open-source feature management and experimentation platform. It provides feature flags with targeted rollouts, Bayesian A/B testing, automated progressive rollouts, and audit logging. It is designed to scale from Docker Compose (small/medium) to Kubernetes (100M+ users).

## Repository Structure

```
cmd/           # Four binary entrypoints: api, web, batch, subscriber
pkg/           # All Go business logic, organized by domain
proto/         # Protobuf definitions; generated .pb.go files live here too
ui/dashboard/  # React 19 + Vite + TailwindCSS + TypeScript frontend (Admin console)
manifests/     # Helm charts for Kubernetes deployment
migration/     # Atlas-managed SQL migrations (mysql/ and postgres/)
docker-compose/ # Docker Compose local dev setup
tools/         # Dev tooling (certs, tokens, etc.)
evaluation/    # Go and other SDK evaluation logic
hack/          # One-off utility scripts
```

### Four Backend Services

| Binary | Role |
|--------|------|
| `cmd/api` | API Gateway — SDK-facing gRPC/HTTP endpoints |
| `cmd/web` | Web Gateway — Admin UI and internal service APIs |
| `cmd/batch` | Batch processor — scheduled jobs (experiment calculation, rollups) |
| `cmd/subscriber` | Event subscriber — processes Pub/Sub events asynchronously |

### pkg/ Domain Layout

Most domain packages follow this structure (not all layers are present in every domain — e.g. `auth` has `api/` and `client/` but no `domain/` or `storage/`):
- `domain/` — core business logic and domain types
- `api/` — gRPC service handlers
- `storage/v2/` — MySQL/Postgres persistence (BigQuery in `storage/bigquery/`)
- `client/` — gRPC client for inter-service communication
- `mock/` — generated mocks (do not edit manually)

Key domains: `feature`, `experiment`, `autoops`, `eventcounter`, `account`, `environment`, `auth`, `notification`, `push`, `auditlog`, `insights`, `coderef`, `tag`, `team`.

## Commands

### Install local dev tools (run once)
```bash
make local-deps
```
Installs: `goimports`, `golangci-lint`, `mockgen`, `protoc-gen-go`, `protoc-gen-grpc-gateway`, `protoc-gen-openapiv2`, `protolock`, `yq`.

Also required (not managed by `local-deps`):
- `protoc` v23.4 (`libprotoc 23.4`) — the repo expects exactly this version; generated `.pb.go` files must show `protoc v4.23.4`
- `clang-format` — required for `proto fmt` step (`brew install clang-format`)

### Build
```bash
make build-go          # Build all Go binaries into bin/
make build-<service>   # Build a single binary, e.g. make build-api
make build-web-console # Build the React frontend
```

### Test
```bash
make start-httpstan    # Required before running tests (Bayesian experiment package needs it)
make test-go           # Run all Go unit tests
make stop-httpstan     # Optional cleanup after tests

# Run a single test or package directly:
TZ=UTC CGO_ENABLED=0 go test -v ./pkg/feature/...
TZ=UTC CGO_ENABLED=0 go test -v -run TestFoo ./pkg/feature/domain/
```

Tests use a table-driven format (`[]struct{ ... }` test cases iterated with `t.Run`). Follow this pattern when adding new tests.

### Lint
```bash
make lint                # golangci-lint on cmd/, pkg/, evaluation/go/, hack/, test/
make gofmt               # Format with goimports (run after any Go changes)
yarn style:all --write   # Format after making changes to TypeScript files
```

### Code generation
```bash
make generate-all      # Regenerate proto Go files + mocks (runs proto-all + mockgen)
make proto-all         # Regenerate only proto Go files + OpenAPI/Swagger specs
make mockgen           # Regenerate only mocks (after changing interfaces that have mocks)
```
- Run `make proto-all` after any `.proto` file change — this regenerates Go bindings and OpenAPI/Swagger specs.
- Run `make mockgen` after changing any Go interface that has generated mocks in a `mock/` directory.
- The generated files are committed to the repo. `protoc` v23.4 must be on PATH ahead of any other version for the version header in `.pb.go` files to stay at `v4.23.4`.

### Local development

**Docker Compose (recommended for most development):**
```bash
make docker-compose-up      # Start all services
make docker-compose-status  # Check status
make docker-compose-logs    # View logs
make docker-compose-down    # Stop services
make docker-compose-clean   # Remove all containers, networks, volumes
```

Add to `/etc/hosts`:
```
127.0.0.1 web-gateway.bucketeer.io
127.0.0.1 api-gateway.bucketeer.io
```

**Minikube (Kubernetes-based):**
```bash
make start-minikube    # Always use this, not `minikube start` directly
make deploy-bucketeer  # Deploy all Helm charts
```

### Database migrations
```bash
make migration-validate     # Validate migration files with Atlas
make migration-hash-check   # Check Atlas migration hash is up to date
```
Migrations live in `migration/mysql/` and `migration/postgres/` and are managed with [Atlas](https://atlasgo.io/).

## Architecture Notes

### Proto-first API design
All service interfaces are defined in `proto/` and generated into `proto/<domain>/*.pb.go` and `*.pb.gw.go`. The `omitempty` JSON tag is intentionally stripped from all generated files (see `proto/Makefile`) because legacy projects rely on empty `environment_id` fields being serialized.

### Event-driven processing
The `subscriber` service consumes Google Cloud Pub/Sub events. Its runtime behavior is configured via JSON files in `docker-compose/config/subscriber-config/`: `subscribers.json`, `onDemandSubscribers.json`, `processors.json`, `onDemandProcessors.json`.

### Data warehouse flexibility
The codebase supports MySQL, PostgreSQL, and BigQuery as the event data warehouse. The active backend is controlled by `dataWarehouse.type` in `manifests/bucketeer/values.dev.yaml`. The Makefile auto-detects this to set `POSTGRES_ENABLED` / `BIGQUERY_ENABLED`. The `web` and `subscriber` services must use the same event store.

### Mock generation
Mocks are generated by `go generate -run="mockgen"` triggered from `//go:generate` directives in the source files. They are regenerated by `make mockgen` (or `make generate-all`). Never edit files in `mock/` directories manually.

### Module path
The Go module is `github.com/bucketeer-io/bucketeer/v2`. Use this prefix for all internal imports.
