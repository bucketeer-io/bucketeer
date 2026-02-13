-- Create scheduled_feature_change table
-- This replaces the old scheduled_flag_update table with a more comprehensive schema
-- that supports structured change payloads, conflict detection, and executor locking.

CREATE TABLE IF NOT EXISTS scheduled_feature_change (
    -- Identity
    id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    
    -- Scheduling
    scheduled_at BIGINT NOT NULL,
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    
    -- Content
    payload JSON NOT NULL,                    -- ScheduledChangePayload as JSON (supports JSON queries)
    comment TEXT,
    
    -- Status tracking
    -- 1=PENDING, 2=EXECUTED, 3=FAILED, 4=CANCELLED, 5=CONFLICT
    status TINYINT NOT NULL DEFAULT 1,
    failure_reason TEXT,
    
    -- Conflict detection
    flag_version_at_creation INT NOT NULL,
    conflicts JSON,                           -- Array of ScheduledChangeConflict as JSON
    
    -- Concurrency control (for executor); DB-only, not exposed via ScheduledFlagChange proto
    locked_at BIGINT,                         -- When executor locked the row (executor coordination only)
    locked_by VARCHAR(255),                   -- Which executor instance locked it (executor coordination only)
    
    -- Audit
    created_by VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_by VARCHAR(255),
    updated_at BIGINT NOT NULL,
    executed_at BIGINT,
    
    -- Keys & Constraints
    PRIMARY KEY (id),
    CONSTRAINT fk_scheduled_feature_change_feature
        FOREIGN KEY (feature_id, environment_id)
        REFERENCES feature(id, environment_id)
        ON DELETE RESTRICT,
    
    -- Indexes for performance
    INDEX idx_scheduled_at_status (scheduled_at, status),
    INDEX idx_feature_env (feature_id, environment_id),
    INDEX idx_environment_status (environment_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
