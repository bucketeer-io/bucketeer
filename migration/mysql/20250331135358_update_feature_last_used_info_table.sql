-- Create a composite index on feature_id, environment_id, and version
CREATE INDEX idx_flui ON feature_last_used_info (feature_id, environment_id, version);
