SELECT
    COUNT(DISTINCT feature.id) AS total_features,
    COUNT(DISTINCT CASE WHEN
        (SELECT flui.last_used_at
            FROM feature_last_used_info flui
            WHERE flui.feature_id = feature.id
            AND flui.environment_id = feature.environment_id
            ORDER BY flui.last_used_at DESC, flui.version DESC
            LIMIT 1)
        >= EXTRACT(EPOCH FROM NOW())::BIGINT - 604800
        THEN feature.id
    END) AS active,
    COUNT(DISTINCT CASE WHEN
        (SELECT flui.last_used_at
            FROM feature_last_used_info flui
            WHERE flui.feature_id = feature.id
            AND flui.environment_id = feature.environment_id
            ORDER BY flui.last_used_at DESC, flui.version DESC
            LIMIT 1)
        < EXTRACT(EPOCH FROM NOW())::BIGINT - 604800
        THEN feature.id
    END) AS inactive
FROM
    feature
WHERE
    feature.archived = FALSE AND
    feature.deleted = FALSE AND
    feature.environment_id = $1
