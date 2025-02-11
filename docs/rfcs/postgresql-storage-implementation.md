# Summary

This RFC proposes implementing PostgreSQL as an alternative storage solution for Bucketeer, targeting open-source deployments and smaller companies. The focus is on SQL compatibility with MySQL and implementing analytics capabilities to replace BigQuery functionality.

# Background

Currently, Bucketeer uses MySQL and BigQuery for data storage and analytics. While this architecture serves well for large-scale deployments, it introduces complexity in setup and maintenance for smaller deployments. PostgreSQL offers a powerful alternative that can handle both transactional and analytical workloads in a single database system.

# Goals

- Provide SQL migration guidelines from MySQL to PostgreSQL
- Implement BigQuery-like analytics capabilities using PostgreSQL features
- Support efficient querying for both OLTP and OLAP workloads
- Maintain performance for small to medium-scale deployments

# Implementation Details

## SQL Migration Considerations

### MySQL to PostgreSQL Syntax Differences

Key differences to handle in SQL migration:

1. **Data Types**
   ```sql
   -- MySQL
   UNSIGNED INTEGER -> INTEGER CHECK (column_name >= 0)
   DATETIME -> TIMESTAMP WITH TIME ZONE
   BOOL -> BOOLEAN
   ```

2. **Auto-increment**
   ```sql
   -- MySQL
   id INT AUTO_INCREMENT
   
   -- PostgreSQL
   id SERIAL
   ```

3. **String Functions**
   ```sql
   -- MySQL
   CONCAT(str1, str2) -> same in PostgreSQL
   IFNULL() -> COALESCE()
   NOW() -> CURRENT_TIMESTAMP
   ```

4. **Group By Handling**
   ```sql
   -- MySQL allows columns not in GROUP BY
   SELECT id, name, COUNT(*) FROM features GROUP BY id;
   
   -- PostgreSQL requires all non-aggregated columns
   SELECT id, name, COUNT(*) FROM features GROUP BY id, name;
   ```

## Analytics Implementation

### TimescaleDB for Time-Series Data

1. **Setup TimescaleDB**
   ```sql
   -- Enable TimescaleDB extension
   CREATE EXTENSION IF NOT EXISTS timescaledb;
   ```

2. **Evaluation Events Table**
   ```sql
   -- Create regular table first
   CREATE TABLE evaluation_events (
       id VARCHAR(255),
       environment_id VARCHAR(255),
       timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
       feature_id VARCHAR(255),
       feature_version INTEGER,
       user_id VARCHAR(255),
       user_data JSONB,
       variation_id VARCHAR(255),
       reason VARCHAR(255),
       tag VARCHAR(255),
       source_id VARCHAR(255)
   );

   -- Convert to hypertable with 1 month chunks
   SELECT create_hypertable('evaluation_events', 'timestamp',
       chunk_time_interval => INTERVAL '1 month');

   -- Create compression policy
   ALTER TABLE evaluation_events SET (
       timescaledb.compress,
       timescaledb.compress_segmentby = 'environment_id,feature_id',
       timescaledb.compress_orderby = 'timestamp DESC'
   );

   -- Automatically compress chunks older than 7 days
   SELECT add_compression_policy('evaluation_events', 
       INTERVAL '7 days');
   ```

3. **Goal Events Table**
   ```sql
   CREATE TABLE goal_events (
       id VARCHAR(255),
       environment_id VARCHAR(255),
       timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
       goal_id VARCHAR(255),
       value DOUBLE PRECISION,
       user_id VARCHAR(255),
       user_data JSONB,
       tag VARCHAR(255),
       source_id VARCHAR(255),
       feature_id VARCHAR(255),
       feature_version INTEGER,
       variation_id VARCHAR(255),
       reason VARCHAR(255)
   );

   -- Convert to hypertable with 1 month chunks
   SELECT create_hypertable('goal_events', 'timestamp',
       chunk_time_interval => INTERVAL '1 month');

   -- Create compression policy
   ALTER TABLE goal_events SET (
       timescaledb.compress,
       timescaledb.compress_segmentby = 'environment_id,goal_id,feature_id',
       timescaledb.compress_orderby = 'timestamp DESC'
   );

   -- Automatically compress chunks older than 7 days
   SELECT add_compression_policy('goal_events', 
       INTERVAL '7 days');
   ```

4. **TimescaleDB-Specific Optimizations**
   ```sql
   -- Enable parallel query for hypertables
   ALTER DATABASE bucketeer SET timescaledb.max_parallel_chunk_scan = 4;

   -- Set retention policy (e.g., keep 1 year of raw data)
   SELECT add_retention_policy('evaluation_events', INTERVAL '1 year');
   SELECT add_retention_policy('goal_events', INTERVAL '1 year');

   -- Create cagg refresh background job
   SELECT add_job('refresh_continuous_aggregate', '1h');
   ```

### Performance Benefits of TimescaleDB

1. **Automatic Chunk Management**
   - Automatic creation of time-based chunks
   - Efficient query planning based on time ranges
   - Automatic removal of old data with retention policies

2. **Optimized Time-Series Operations**
   - Better compression for time-series data
   - Efficient time-based queries
   - Automatic parallel query execution

3. **Query Examples**
   ```sql
   -- Efficient time-range queries
   SELECT feature_id, COUNT(*)
   FROM evaluation_events
   WHERE timestamp >= NOW() - INTERVAL '7 days'
   GROUP BY feature_id;

   -- Using continuous aggregates
   SELECT day, feature_id, total_evaluations
   FROM daily_evaluation_metrics
   WHERE day >= NOW() - INTERVAL '30 days'
   ORDER BY day DESC;
   ```

# Trade-offs

## Advantages

1. **Single System Analytics**
   - No need for separate analytics infrastructure
   - Real-time analytics capability
   - Simplified operational complexity

2. **SQL Compatibility**
   - Minor syntax differences from MySQL
   - Standard SQL compliance
   - Rich set of analytical functions

3. **Cost Effectiveness**
   - No data transfer between systems
   - Reduced cloud service costs
   - Simplified licensing model

## Disadvantages

1. **Scale Limitations**
   - Limited by single instance capacity
   - Not suitable for petabyte-scale analytics
   - Resource intensive for very large datasets

2. **Query Performance**
   - May require more optimization than BigQuery
   - Limited parallel processing capabilities
   - Memory constraints for large aggregations
