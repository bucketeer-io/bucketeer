-- The last evaluation time of a flag is stored on the feature row itself.
-- The latest value across every version in feature_last_used_info is copied over
-- so existing flags keep their status; that table is dropped in a later migration.
ALTER TABLE `feature` ADD COLUMN `last_used_at` bigint NOT NULL DEFAULT 0, ALGORITHM=INSTANT;
UPDATE `feature` f
JOIN (
    SELECT feature_id, environment_id, MAX(last_used_at) AS last_used_at
    FROM feature_last_used_info
    GROUP BY feature_id, environment_id
) flui ON f.id = flui.feature_id AND f.environment_id = flui.environment_id
SET f.last_used_at = flui.last_used_at;
