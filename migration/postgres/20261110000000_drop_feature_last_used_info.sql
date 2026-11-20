-- The last evaluation time now lives in feature.last_used_at.
-- Evaluations the previous subscriber wrote to feature_last_used_info after the
-- backfill are merged first so no flag moves backwards, then the table is dropped.
UPDATE feature f
SET last_used_at = GREATEST(f.last_used_at, flui.last_used_at)
FROM (
    SELECT feature_id, environment_id, MAX(last_used_at) AS last_used_at
    FROM feature_last_used_info
    GROUP BY feature_id, environment_id
) flui
WHERE f.id = flui.feature_id AND f.environment_id = flui.environment_id;
DROP TABLE IF EXISTS feature_last_used_info;
