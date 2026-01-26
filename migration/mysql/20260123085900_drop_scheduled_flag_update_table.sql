-- Drop the old scheduled_flag_update table
-- This table is being replaced by the new scheduled_feature_change table
-- which provides a more comprehensive schema for scheduled flag changes.

DROP TABLE IF EXISTS scheduled_flag_update;
