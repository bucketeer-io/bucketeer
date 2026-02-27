-- Bucketeer PostgreSQL Schema Initialization

-- ============================================
-- Core Organization Tables
-- ============================================

-- Create "organization" table
CREATE TABLE organization (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    url_code VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    trial BOOLEAN NOT NULL DEFAULT FALSE,
    system_admin BOOLEAN NOT NULL DEFAULT FALSE,
    owner_email VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX unique_url_code ON organization (url_code);
CREATE INDEX idx_organization_disabled ON organization (disabled);
CREATE INDEX idx_organization_archived ON organization (archived);
CREATE INDEX idx_organization_disabled_archived_id ON organization (disabled, archived, id);

-- Create "project" table
CREATE TABLE project (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    url_code VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    trial BOOLEAN NOT NULL DEFAULT FALSE,
    creator_email VARCHAR(255) NOT NULL,
    organization_id VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX unique_organization_url_code ON project (organization_id, url_code);
CREATE INDEX idx_project_organization_id ON project (organization_id);

-- Create "environment_v2" table
CREATE TABLE environment_v2 (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    url_code VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    project_id VARCHAR(255) NOT NULL,
    organization_id VARCHAR(255) NOT NULL,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    require_comment BOOLEAN NOT NULL DEFAULT TRUE,
    auto_archive_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    auto_archive_unused_days INTEGER NOT NULL DEFAULT 90,
    auto_archive_check_code_refs BOOLEAN NOT NULL DEFAULT TRUE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT environment_v2_foreign_project_id FOREIGN KEY (project_id) REFERENCES project (id)
);
CREATE UNIQUE INDEX unique_project_id_url_code ON environment_v2 (project_id, url_code);
CREATE INDEX idx_environment_v2_organization_id ON environment_v2 (organization_id);
CREATE INDEX idx_environment_auto_archive_enabled ON environment_v2 (auto_archive_enabled);

-- Create "team" table
CREATE TABLE team (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    organization_id VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX unique_team_name_org ON team (name, organization_id);
CREATE INDEX idx_team_organization_id ON team (organization_id);

-- ============================================
-- Account Tables
-- ============================================

-- Create "account_v2" table
CREATE TABLE account_v2 (
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL DEFAULT '',
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) NOT NULL DEFAULT '',
    language VARCHAR(10) NOT NULL DEFAULT '',
    avatar_image_url VARCHAR(255) NOT NULL,
    avatar_file_type VARCHAR(50) NOT NULL DEFAULT '',
    avatar_image BYTEA,
    tags JSONB NOT NULL DEFAULT '[]',
    teams JSONB,
    organization_id VARCHAR(255) NOT NULL,
    organization_role INTEGER NOT NULL,
    environment_roles JSONB NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    search_filters JSONB,
    last_seen BIGINT NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (email, organization_id),
    CONSTRAINT account_v2_foreign_organization_id FOREIGN KEY (organization_id) REFERENCES organization (id)
);
CREATE INDEX idx_account_v2_organization_id ON account_v2 (organization_id);
CREATE INDEX idx_account_v2_email ON account_v2 (email);

-- Create "admin_account" table
CREATE TABLE admin_account (
    id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role INTEGER NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX unique_admin_email ON admin_account (email);

-- ============================================
-- Feature Flag Tables
-- ============================================

-- Create "feature" table
CREATE TABLE feature (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(511) NOT NULL,
    description TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    evaluation_undelayable BOOLEAN NOT NULL DEFAULT FALSE,
    ttl INTEGER NOT NULL,
    version INTEGER NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    variations JSONB NOT NULL,
    targets JSONB NOT NULL,
    rules JSONB NOT NULL,
    default_strategy JSONB NOT NULL,
    off_variation VARCHAR(255) NOT NULL,
    tags JSONB NOT NULL,
    maintainer VARCHAR(255) NOT NULL,
    variation_type INTEGER NOT NULL,
    sampling_seed VARCHAR(255) NOT NULL DEFAULT '',
    prerequisites JSONB,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id)
);

-- Create "feature_last_used_info" table
CREATE TABLE feature_last_used_info (
    id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL,
    last_used_at BIGINT NOT NULL,
    client_oldest_version VARCHAR(255) NOT NULL,
    client_latest_version VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id)
);
CREATE INDEX idx_flui ON feature_last_used_info (feature_id, environment_id, version);

-- Create "flag_trigger" table
CREATE TABLE flag_trigger (
    id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    type INTEGER NOT NULL,
    action BOOLEAN NOT NULL,
    description TEXT NOT NULL,
    trigger_count INTEGER NOT NULL,
    last_triggered_at BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id),
    CONSTRAINT foreign_flag_trigger_feature FOREIGN KEY (feature_id, environment_id) REFERENCES feature (id, environment_id)
);
CREATE INDEX idx_flag_trigger_feature ON flag_trigger (feature_id, environment_id);

-- Create "scheduled_feature_change" table
CREATE TABLE scheduled_feature_change (
    id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    scheduled_at BIGINT NOT NULL,
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    payload JSONB NOT NULL,
    comment TEXT,
    status SMALLINT NOT NULL DEFAULT 1,
    failure_reason TEXT,
    flag_version_at_creation INTEGER NOT NULL,
    conflicts JSONB,
    locked_at BIGINT,
    locked_by VARCHAR(255),
    created_by VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_by VARCHAR(255),
    updated_at BIGINT NOT NULL,
    executed_at BIGINT,
    PRIMARY KEY (id),
    CONSTRAINT fk_scheduled_feature_change_feature FOREIGN KEY (feature_id, environment_id) REFERENCES feature (id, environment_id) ON DELETE RESTRICT
);
CREATE INDEX idx_scheduled_at_status ON scheduled_feature_change (scheduled_at, status);
CREATE INDEX idx_sfc_feature_env ON scheduled_feature_change (feature_id, environment_id);
CREATE INDEX idx_sfc_environment_status ON scheduled_feature_change (environment_id, status);

-- ============================================
-- Auto Operations Tables
-- ============================================

-- Create "auto_ops_rule" table
CREATE TABLE auto_ops_rule (
    id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    ops_type INTEGER NOT NULL,
    clauses JSONB NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    status INTEGER NOT NULL DEFAULT 0,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id),
    CONSTRAINT foreign_auto_ops_rule_feature FOREIGN KEY (feature_id, environment_id) REFERENCES feature (id, environment_id)
);
CREATE INDEX idx_auto_ops_rule_feature ON auto_ops_rule (feature_id, environment_id);

-- Create "ops_count" table
CREATE TABLE ops_count (
    id VARCHAR(255) NOT NULL,
    auto_ops_rule_id VARCHAR(255) NOT NULL,
    clause_id VARCHAR(255) NOT NULL,
    updated_at BIGINT NOT NULL,
    ops_event_count BIGINT NOT NULL,
    evaluation_count BIGINT NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id),
    CONSTRAINT foreign_ops_count_auto_ops_rule FOREIGN KEY (auto_ops_rule_id, environment_id) REFERENCES auto_ops_rule (id, environment_id),
    CONSTRAINT foreign_ops_count_feature FOREIGN KEY (feature_id, environment_id) REFERENCES feature (id, environment_id)
);
CREATE INDEX idx_ops_count_auto_ops_rule ON ops_count (auto_ops_rule_id, environment_id);
CREATE INDEX idx_ops_count_feature ON ops_count (feature_id, environment_id);

-- Create "ops_progressive_rollout" table
CREATE TABLE ops_progressive_rollout (
    id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    clause JSONB NOT NULL,
    status INTEGER NOT NULL,
    stopped_by INTEGER NOT NULL DEFAULT 0,
    type INTEGER NOT NULL,
    stopped_at BIGINT NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id),
    CONSTRAINT foreign_progressive_rollout_feature FOREIGN KEY (feature_id, environment_id) REFERENCES feature (id, environment_id)
);
CREATE INDEX idx_progressive_rollout_feature ON ops_progressive_rollout (feature_id, environment_id);

-- ============================================
-- Experiment Tables
-- ============================================

-- Create "goal" table
CREATE TABLE goal (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(511) NOT NULL,
    description TEXT NOT NULL,
    connection_type INTEGER NOT NULL DEFAULT 0,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id)
);

-- Create "experiment" table
CREATE TABLE experiment (
    id VARCHAR(255) NOT NULL,
    goal_id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    feature_version INTEGER NOT NULL,
    variations JSONB NOT NULL,
    start_at BIGINT NOT NULL,
    stop_at BIGINT NOT NULL,
    stopped BOOLEAN NOT NULL DEFAULT FALSE,
    stopped_at BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    goal_ids JSONB NOT NULL,
    name VARCHAR(511) NOT NULL,
    description TEXT NOT NULL,
    base_variation_id VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL,
    maintainer VARCHAR(255) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id),
    CONSTRAINT foreign_experiment_feature FOREIGN KEY (feature_id, environment_id) REFERENCES feature (id, environment_id)
);
CREATE INDEX idx_experiment_feature ON experiment (feature_id, environment_id);

-- Create "experiment_result" table
CREATE TABLE experiment_result (
    id VARCHAR(255) NOT NULL,
    experiment_id VARCHAR(255) NOT NULL,
    updated_at BIGINT NOT NULL,
    data JSONB NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id)
);

-- ============================================
-- Segment Tables
-- ============================================

-- Create "segment" table
CREATE TABLE segment (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    rules JSONB NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    version BIGINT NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    included_user_count BIGINT NOT NULL,
    excluded_user_count BIGINT NOT NULL,
    status INTEGER NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id)
);

-- Create "segment_user" table
CREATE TABLE segment_user (
    id VARCHAR(511) NOT NULL,
    segment_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    state INTEGER NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id),
    CONSTRAINT foreign_segment_user_segment FOREIGN KEY (segment_id, environment_id) REFERENCES segment (id, environment_id)
);
CREATE INDEX idx_segment_user_segment ON segment_user (segment_id, environment_id);

-- ============================================
-- API & Audit Tables
-- ============================================

-- Create "api_key" table
CREATE TABLE api_key (
    id VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL DEFAULT '',
    name VARCHAR(255) NOT NULL,
    role INTEGER NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    maintainer VARCHAR(255) NOT NULL DEFAULT '',
    description VARCHAR(255) NOT NULL DEFAULT '',
    last_used_at BIGINT DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id)
);

-- Create "audit_log" table
CREATE TABLE audit_log (
    id VARCHAR(255) NOT NULL,
    timestamp BIGINT NOT NULL,
    entity_type INTEGER NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    type INTEGER NOT NULL,
    event JSONB NOT NULL,
    editor JSONB NOT NULL,
    options JSONB NOT NULL,
    entity_data TEXT NOT NULL,
    previous_entity_data TEXT NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id)
);
CREATE INDEX idx_audit_log_entity ON audit_log (entity_type, entity_id);
CREATE INDEX idx_audit_log_timestamp_desc ON audit_log (timestamp DESC);
CREATE INDEX idx_audit_log_environment_timestamp ON audit_log (environment_id, timestamp DESC);

-- Create "admin_audit_log" table
CREATE TABLE admin_audit_log (
    id VARCHAR(255) NOT NULL,
    timestamp BIGINT NOT NULL,
    entity_type INTEGER NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    type INTEGER NOT NULL,
    event JSONB NOT NULL,
    editor JSONB NOT NULL,
    options JSONB NOT NULL,
    entity_data TEXT NOT NULL,
    previous_entity_data TEXT NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX idx_admin_audit_log_timestamp_desc ON admin_audit_log (timestamp DESC);

-- ============================================
-- Subscription & Push Tables
-- ============================================

-- Create "subscription" table
CREATE TABLE subscription (
    id VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    source_types JSONB NOT NULL,
    recipient JSONB NOT NULL,
    name VARCHAR(255) NOT NULL,
    feature_flag_tags JSONB NOT NULL DEFAULT '[]',
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id)
);

-- Create "admin_subscription" table
CREATE TABLE admin_subscription (
    id VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    source_types JSONB NOT NULL,
    recipient JSONB NOT NULL,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

-- Create "push" table
CREATE TABLE push (
    id VARCHAR(255) NOT NULL,
    fcm_api_key VARCHAR(511),
    fcm_service_account JSONB NOT NULL,
    tags JSONB NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, environment_id)
);

-- ============================================
-- Tag Table
-- ============================================

-- Create "tag" table
CREATE TABLE tag (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    entity_type INTEGER NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX unique_tag_name_env_entity ON tag (name, environment_id, entity_type);

-- ============================================
-- Code Reference Table
-- ============================================

-- Create "code_reference" table
CREATE TABLE code_reference (
    id VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    file_extension VARCHAR(32) NOT NULL DEFAULT '',
    line_number INTEGER NOT NULL,
    code_snippet TEXT NOT NULL,
    content_hash VARCHAR(64) NOT NULL,
    aliases JSONB,
    repository_name VARCHAR(255) NOT NULL,
    repository_owner VARCHAR(255) NOT NULL,
    repository_type SMALLINT NOT NULL,
    repository_branch VARCHAR(255) NOT NULL,
    commit_hash VARCHAR(40) NOT NULL,
    environment_id VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT foreign_code_references_feature FOREIGN KEY (feature_id, environment_id) REFERENCES feature (id, environment_id)
);
CREATE INDEX idx_code_reference_file_path ON code_reference (file_path);

-- ============================================
-- Schema Migration Tables
-- ============================================

-- Create "schema_migrations" table
CREATE TABLE schema_migrations (
    version BIGINT NOT NULL,
    dirty BOOLEAN NOT NULL,
    PRIMARY KEY (version)
);

-- Create "atlas_schema_revisions" table
CREATE TABLE IF NOT EXISTS atlas_schema_revisions (
    version VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    type BIGINT NOT NULL DEFAULT 2,
    applied BIGINT NOT NULL DEFAULT 0,
    total BIGINT NOT NULL DEFAULT 0,
    executed_at TIMESTAMP NOT NULL,
    execution_time BIGINT NOT NULL,
    error TEXT,
    error_stmt TEXT,
    hash VARCHAR(255) NOT NULL,
    partial_hashes JSONB,
    operator_version VARCHAR(255) NOT NULL,
    PRIMARY KEY (version)
);
