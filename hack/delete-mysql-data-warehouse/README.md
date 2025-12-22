# Delete MySQL Data Warehouse

This tool truncates the MySQL data warehouse tables (`evaluation_event` and `goal_event`).

## Usage

```bash
# Build the tool
make build

# Run the tool
./delete-mysql-data-warehouse truncate \
    --mysql-user=bucketeer \
    --mysql-pass=bucketeer \
    --mysql-host=localhost \
    --mysql-port=3306 \
    --mysql-db-name=bucketeer \
    --no-profile \
    --no-gcp-trace-enabled
```

## From the root Makefile

```bash
# For minikube
make delete-mysql-data-warehouse-data

# For docker-compose
make docker-compose-delete-mysql-data-warehouse-data
```
