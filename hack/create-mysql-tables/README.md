# MySQL Storage for Bucketeer Events

This directory contains tools for setting up MySQL as a lightweight alternative to BigQuery for storing Bucketeer event data. 

## Overview

Bucketeer events (evaluation and goal events) can be stored in MySQL instead of BigQuery. This provides a simpler, more lightweight setup that's easier to deploy in environments where BigQuery might not be available or practical.

## Features

- MySQL-based storage for evaluation and goal events
- Compatible with existing Bucketeer interfaces
- Simple setup and configuration
- Lower operational overhead compared to BigQuery

## Setup

### Prerequisites

- MySQL 5.7+ or MariaDB 10.2+ server
- Go 1.16+

### Environment Variables

Configure the MySQL connection using these environment variables:

- `MYSQL_HOST` - MySQL host (default: "localhost")
- `MYSQL_PORT` - MySQL port (default: "3306")
- `MYSQL_USER` - MySQL username (default: "root")
- `MYSQL_PASSWORD` - MySQL password (default: "")
- `SCHEMA_FILE` - Path to SQL schema file (default: "./create_tables.sql")

### Creating Tables

Run the following command to create the required MySQL tables:

```bash
cd hack/create-mysql-tables
go run main.go
```

## Usage

To use MySQL for event storage in your Bucketeer services, you'll need to configure your service to use the MySQL implementation instead of BigQuery.

Example configuration:

```go
import (
    "github.com/bucketeer-io/bucketeer/pkg/subscriber/storage/v2"
)

// Create a MySQL evaluation event writer
evalEventWriter, err := v2.NewEvalEventAdapter(
    ctx,
    "root:password@tcp(localhost:3306)/bucketeer",
    100, // batch size
    logger,
)
if err != nil {
    // handle error
}

// Create a MySQL goal event writer
goalEventWriter, err := v2.NewGoalEventAdapter(
    ctx,
    "root:password@tcp(localhost:3306)/bucketeer",
    100, // batch size
    logger,
)
if err != nil {
    // handle error
}
```

## Table Schemas

The MySQL implementation creates two main tables:

### evaluation_event

Stores feature flag evaluation events with the following schema:
- `id` - Unique event ID
- `environment_id` - Environment ID
- `timestamp` - When the event occurred
- `feature_id` - Feature ID
- `feature_version` - Feature version
- `user_id` - User ID
- `user_data` - User data (JSON)
- `variation_id` - Variation ID
- `reason` - Reason for the evaluation
- `tag` - Tag
- `source_id` - Source ID
- `created_at` - When the record was created

### goal_event

Stores goal events with the following schema:
- `id` - Unique event ID
- `environment_id` - Environment ID
- `timestamp` - When the event occurred
- `goal_id` - Goal ID
- `value` - Goal value
- `user_id` - User ID
- `user_data` - User data (JSON)
- `tag` - Tag
- `source_id` - Source ID
- `feature_id` - Feature ID
- `feature_version` - Feature version
- `variation_id` - Variation ID
- `reason` - Reason
- `created_at` - When the record was created 