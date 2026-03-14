-- TimescaleDB Event Tables

-- Enable TimescaleDB extension (run as superuser if not already enabled)
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Evaluation Event Table (hypertable)
CREATE TABLE IF NOT EXISTS evaluation_event (
    id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    "timestamp" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    feature_version INT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    user_data JSONB,
    variation_id VARCHAR(255) NOT NULL,
    reason TEXT,
    tag VARCHAR(255),
    source_id VARCHAR(255),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("timestamp", id)
);

-- Convert to hypertable with 1 month chunks (time-based partitioning)
SELECT create_hypertable(
    'evaluation_event',
    'timestamp',
    chunk_time_interval => INTERVAL '1 month',
    if_not_exists => TRUE
);

-- Compression: segment by common filter columns, order by time for range queries
ALTER TABLE evaluation_event SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'environment_id, feature_id',
    timescaledb.compress_orderby = '"timestamp" DESC'
);

-- Compress chunks older than 7 days (idempotent: skip if policy exists)
SELECT add_compression_policy('evaluation_event', INTERVAL '7 days', if_not_exists => TRUE);

-- Optional: retain raw data for 1 year (adjust or remove as needed)
-- SELECT add_retention_policy('evaluation_event', INTERVAL '1 year');

-- Indexes for common filters (created after hypertable)
CREATE INDEX IF NOT EXISTS idx_evaluation_environment_id ON evaluation_event (environment_id);
CREATE INDEX IF NOT EXISTS idx_evaluation_timestamp ON evaluation_event ("timestamp");
CREATE INDEX IF NOT EXISTS idx_evaluation_feature_id ON evaluation_event (feature_id);
CREATE INDEX IF NOT EXISTS idx_evaluation_user_id ON evaluation_event (user_id);
CREATE INDEX IF NOT EXISTS idx_evaluation_variation_id ON evaluation_event (variation_id);

-- Goal Event Table (hypertable)
CREATE TABLE IF NOT EXISTS goal_event (
    id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    "timestamp" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    goal_id VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION,
    user_id VARCHAR(255) NOT NULL,
    user_data JSONB,
    tag VARCHAR(255),
    source_id VARCHAR(255),
    feature_id VARCHAR(255),
    feature_version INT,
    variation_id VARCHAR(255),
    reason TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("timestamp", id)
);

SELECT create_hypertable(
    'goal_event',
    'timestamp',
    chunk_time_interval => INTERVAL '1 month',
    if_not_exists => TRUE
);

ALTER TABLE goal_event SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'environment_id, goal_id, feature_id',
    timescaledb.compress_orderby = '"timestamp" DESC'
);

SELECT add_compression_policy('goal_event', INTERVAL '7 days', if_not_exists => TRUE);

-- Optional: SELECT add_retention_policy('goal_event', INTERVAL '1 year');

CREATE INDEX IF NOT EXISTS idx_goal_environment_id ON goal_event (environment_id);
CREATE INDEX IF NOT EXISTS idx_goal_timestamp ON goal_event ("timestamp");
CREATE INDEX IF NOT EXISTS idx_goal_goal_id ON goal_event (goal_id);
CREATE INDEX IF NOT EXISTS idx_goal_user_id ON goal_event (user_id);
CREATE INDEX IF NOT EXISTS idx_goal_feature_id ON goal_event (feature_id);
CREATE INDEX IF NOT EXISTS idx_goal_variation_id ON goal_event (variation_id);
