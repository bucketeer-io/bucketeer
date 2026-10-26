-- Universal holdout groups (RFC 0048, phase 1).
-- holdout_group: environment-scoped control group.
-- holdout_group_feature: target flags, soft-removed so re-adding keeps history.

CREATE TABLE IF NOT EXISTS holdout_group (
    id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    name VARCHAR(511) NOT NULL,
    description TEXT NOT NULL,
    assignment_type INT NOT NULL,
    sampling_seed VARCHAR(255) NOT NULL DEFAULT '',
    bucket_start DECIMAL(6,3) NOT NULL DEFAULT 0,
    percentage DECIMAL(6,3) NOT NULL DEFAULT 0,
    segment_id VARCHAR(255) NOT NULL DEFAULT '',
    status INT NOT NULL,
    start_at BIGINT NOT NULL,
    stop_at BIGINT NOT NULL,
    start_analysis_at BIGINT NOT NULL DEFAULT 0,
    analysis_started_at BIGINT NOT NULL DEFAULT 0,
    stopped_at BIGINT NOT NULL DEFAULT 0,
    goal_ids JSON NOT NULL,
    experiment_id VARCHAR(255) NOT NULL DEFAULT '',
    analysis_experiment_id VARCHAR(255) NOT NULL DEFAULT '',
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    archived BOOL NOT NULL DEFAULT 0,
    deleted BOOL NOT NULL DEFAULT 0,
    PRIMARY KEY (id, environment_id),
    INDEX idx_holdout_group_environment_id_status (environment_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS holdout_group_feature (
    id VARCHAR(255) NOT NULL,
    holdout_group_id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    holdout_variation_id VARCHAR(255) NOT NULL,
    added_at BIGINT NOT NULL,
    removed_at BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id, environment_id),
    INDEX idx_holdout_group_feature_group_environment (holdout_group_id, environment_id),
    INDEX idx_holdout_group_feature_feature_environment_removed (feature_id, environment_id, removed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
