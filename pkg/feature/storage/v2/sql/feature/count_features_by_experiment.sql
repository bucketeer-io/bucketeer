SELECT
    COUNT(DISTINCT feature.id)
FROM
    feature
LEFT OUTER JOIN feature_last_used_info ON
    feature.id = feature_last_used_info.feature_id AND
    feature.environment_id = feature_last_used_info.environment_id AND
    NOT EXISTS (
        SELECT 1 FROM feature_last_used_info flui2
        WHERE flui2.feature_id = feature_last_used_info.feature_id
        AND flui2.environment_id = feature_last_used_info.environment_id
        AND (flui2.last_used_at > feature_last_used_info.last_used_at
        OR (flui2.last_used_at = feature_last_used_info.last_used_at AND flui2.version > feature_last_used_info.version))
    )
LEFT JOIN experiment ON
    feature.id = experiment.feature_id AND
    feature.environment_id = experiment.environment_id
