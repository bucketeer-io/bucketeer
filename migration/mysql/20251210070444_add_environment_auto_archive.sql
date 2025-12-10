-- Add auto-archive configuration columns to environment_v2 table
ALTER TABLE environment_v2
  ADD COLUMN auto_archive_enabled BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE environment_v2
  ADD COLUMN auto_archive_unused_days INT NOT NULL DEFAULT 90;

ALTER TABLE environment_v2
  ADD COLUMN auto_archive_check_code_refs BOOLEAN NOT NULL DEFAULT TRUE;

-- Create index for efficient lookup of environments with auto-archive enabled
CREATE INDEX idx_environment_auto_archive_enabled
  ON environment_v2 (auto_archive_enabled);
