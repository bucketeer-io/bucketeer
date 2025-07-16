SELECT DISTINCT
    feature.id,
    feature.name,
    feature.description,
    feature.enabled,
    feature.archived,
    feature.deleted,
    feature.evaluation_undelayable,
    feature.ttl,
    feature.version,
    feature.created_at,
    feature.updated_at,
    feature.variation_type,
    feature.variations,
    feature.targets,
    feature.rules,
    feature.default_strategy,
    feature.off_variation,
    feature.tags,
    feature.maintainer,
    feature.sampling_seed,
    feature.prerequisites,
    (
        SELECT COUNT(aor.id)
        FROM auto_ops_rule aor
        WHERE
            aor.feature_id = feature.id AND
            ops_type = 1 AND
            aor.deleted = 0
    ) AS progressive_rollout_count,
    (
        SELECT COUNT(aor.id)
        FROM auto_ops_rule aor
        WHERE
            aor.feature_id = feature.id AND
            ops_type = 2 AND
            aor.deleted = 0
    ) AS schedule_count,
    (
        SELECT COUNT(aor.id)
        FROM auto_ops_rule aor
        WHERE
            aor.feature_id = feature.id AND
            ops_type = 3 AND
            aor.deleted = 0
    ) AS kill_switch_count,
    COALESCE(
        (SELECT flui.feature_id
         FROM feature_last_used_info flui
         WHERE flui.feature_id = feature.id
         AND flui.environment_id = feature.environment_id
         ORDER BY flui.last_used_at DESC, flui.version DESC
         LIMIT 1),
        ''
    ) AS feature_id,
    COALESCE(
        (SELECT flui.version
        FROM feature_last_used_info flui
        WHERE flui.feature_id = feature.id
        AND flui.environment_id = feature.environment_id
        ORDER BY flui.last_used_at DESC, flui.version DESC
        LIMIT 1),
        0
    ) AS version,
    COALESCE(
        (SELECT flui.last_used_at
        FROM feature_last_used_info flui
        WHERE flui.feature_id = feature.id
        AND flui.environment_id = feature.environment_id
        ORDER BY flui.last_used_at DESC, flui.version DESC
        LIMIT 1),
        0
    ) AS last_used_at,
    COALESCE(
        (SELECT flui.created_at
        FROM feature_last_used_info flui
        WHERE flui.feature_id = feature.id
        AND flui.environment_id = feature.environment_id
        ORDER BY flui.last_used_at DESC, flui.version DESC
        LIMIT 1),
        0
    ) AS created_at,
    COALESCE(
        (SELECT flui.client_oldest_version
        FROM feature_last_used_info flui
        WHERE flui.feature_id = feature.id
        AND flui.environment_id = feature.environment_id
        ORDER BY flui.last_used_at DESC, flui.version DESC
        LIMIT 1),
        ''
    ) AS client_oldest_version,
    COALESCE(
        (SELECT flui.client_latest_version
        FROM feature_last_used_info flui
        WHERE flui.feature_id = feature.id
        AND flui.environment_id = feature.environment_id
        ORDER BY flui.last_used_at DESC, flui.version DESC
        LIMIT 1),
        ''
    ) AS client_latest_version
FROM
    feature
LEFT JOIN experiment ON
    feature.id = experiment.feature_id AND
    feature.environment_id = experiment.environment_id