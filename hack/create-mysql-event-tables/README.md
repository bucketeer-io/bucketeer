# Create MySQL Event Tables

This utility creates MySQL event tables (`evaluation_event` and `goal_event`) for the Bucketeer data warehouse on a separate MySQL database.

## Purpose

When using MySQL as a data warehouse backend, users may want to configure a separate MySQL database instance specifically for the data warehouse, different from the main console database. This utility allows you to create the required event tables on that separate MySQL database.

## Usage

### Using Make Commands (Recommended)

For convenience, you can use the provided make commands:

```bash
# For general usage
MYSQL_HOST=your-host MYSQL_PORT=3306 MYSQL_USER=your-user MYSQL_PASS=your-pass MYSQL_DB_NAME=your-db make create-mysql-event-tables

# For Docker Compose environment (uses default localhost settings)
make docker-compose-create-mysql-event-tables
```

### Direct Go Usage

```bash
go run ./hack/create-mysql-event-tables create \
  --mysql-host=your-datawarehouse-mysql-host \
  --mysql-user=your-mysql-user \
  --mysql-pass=your-mysql-password \
  --mysql-db-name=your-datawarehouse-database
```

### Full Options

```bash
go run ./hack/create-mysql-event-tables create \
  --mysql-host=your-datawarehouse-mysql-host \
  --mysql-port=3306 \
  --mysql-user=your-mysql-user \
  --mysql-pass=your-mysql-password \
  --mysql-db-name=your-datawarehouse-database
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `--mysql-host` | Yes | - | MySQL server hostname or IP address |
| `--mysql-port` | No | 3306 | MySQL server port |
| `--mysql-user` | Yes | - | MySQL username |
| `--mysql-pass` | Yes | - | MySQL password |
| `--mysql-db-name` | Yes | - | MySQL database name for the data warehouse |

## Tables Created

The utility creates the following tables in your specified MySQL database:

1. **evaluation_event** - Stores feature flag evaluation events
2. **goal_event** - Stores goal conversion events

Both tables include appropriate indexes for optimal query performance.

## Example Scenarios

### Scenario 1: Docker Compose Environment
You're using Docker Compose and want to create event tables in the local MySQL database:

```bash
# Uses default Docker Compose MySQL settings (localhost:3306, user:bucketeer, pass:bucketeer, db:bucketeer)
make docker-compose-create-mysql-event-tables
```

### Scenario 2: Separate Data Warehouse Database
Your main Bucketeer console uses one MySQL database, but you want to use a separate MySQL database for analytics/data warehouse purposes.

```bash
# Using make command (recommended)
MYSQL_HOST=analytics-mysql.example.com \
MYSQL_USER=analytics_user \
MYSQL_PASS=analytics_password \
MYSQL_DB_NAME=bucketeer-analytics \
make create-mysql-event-tables

# Or using direct Go command
go run ./hack/create-mysql-event-tables create \
  --mysql-host=analytics-mysql.example.com \
  --mysql-user=analytics_user \
  --mysql-pass=analytics_password \
  --mysql-db-name=bucketeer-analytics
```

### Scenario 3: Different MySQL Server
You want to use a completely different MySQL server instance for the data warehouse.

```bash
# Using make command (recommended)
MYSQL_HOST=data-warehouse-mysql.example.com \
MYSQL_PORT=3307 \
MYSQL_USER=dw_user \
MYSQL_PASS=dw_password \
MYSQL_DB_NAME=bucketeer_dw \
make create-mysql-event-tables

# Or using direct Go command
go run ./hack/create-mysql-event-tables create \
  --mysql-host=data-warehouse-mysql.example.com \
  --mysql-port=3307 \
  --mysql-user=dw_user \
  --mysql-pass=dw_password \
  --mysql-db-name=bucketeer_dw
```

## Prerequisites

- Go installed and configured
- Network access to the target MySQL server
- MySQL user with CREATE TABLE permissions on the specified database
- The target database must already exist

## Notes

- The utility will automatically check if tables already exist before attempting to create them
- If tables already exist, the utility will skip their creation and provide informative messages
- If all required tables already exist, the utility will exit successfully without making any changes
- The utility will automatically split and execute the SQL statements from the migration file
- The utility includes connection timeout handling and detailed logging
- All SQL statements are executed within the same database transaction context 