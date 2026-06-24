# Delete PostgreSQL Data Warehouse

This tool truncates the PostgreSQL data warehouse tables (`evaluation_event` and `goal_event`).

## Usage

```bash
# Build the tool
make build

# Run the tool
./delete-postgres-data-warehouse truncate \
    --postgres-user=bucketeer \
    --postgres-pass=bucketeer \
    --postgres-host=localhost \
    --postgres-port=5432 \
    --postgres-db-name=bucketeer \
    --no-profile \
    --no-gcp-trace-enabled
```

## From the root Makefile

```bash
# For minikube
make delete-postgres-data-warehouse-data

# For docker-compose
make docker-compose-delete-postgres-data-warehouse-data
```

## Create docker image

```bash
make deps

export PAT=<PERSONAL_ACCESS_TOKEN>
export GITHUB_USER_NAME=<GITHUB_USER_NAME>
export TAG=<TAG>

make docker-build
make docker-push
```

Personal Access Token needs to have `write:packages` permission.
