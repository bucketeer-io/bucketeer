-- Modify "audit_log" table
CREATE INDEX idx_environment_id_timestamp_desc ON audit_log (environment_id, timestamp DESC); 