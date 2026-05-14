SELECT
    COUNT(DISTINCT feature.id) AS total_features,
    COUNT(DISTINCT CASE WHEN
        FROM_UNIXTIME(
            (SELECT flui.last_used_at
            FROM feature_last_used_info flui
            WHERE flui.feature_id = feature.id
            AND flui.environment_id = feature.environment_id
            ORDER BY flui.last_used_at DESC, flui.version DESC
            LIMIT 1)
        ) >= DATE_SUB(now(), INTERVAL 7 DAY)
        THEN feature.id
    END) AS active,
    COUNT(DISTINCT CASE WHEN
        FROM_UNIXTIME(
            (SELECT flui.last_used_at
            FROM feature_last_used_info flui
            WHERE flui.feature_id = feature.id
            AND flui.environment_id = feature.environment_id
            ORDER BY flui.last_used_at DESC, flui.version DESC
            LIMIT 1)
        ) < DATE_SUB(now(), INTERVAL 7 DAY)
        THEN feature.id
    END) AS inactive
FROM
    feature
WHERE
    feature.archived = 0 AND
    feature.deleted = 0 AND
    feature.environment_id = ?