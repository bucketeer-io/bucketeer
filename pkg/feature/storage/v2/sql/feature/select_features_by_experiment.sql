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
    COALESCE(flui.feature_id, '') AS feature_id,
    COALESCE(flui.version, 0) AS version,
    COALESCE(flui.last_used_at, 0) AS last_used_at,
    COALESCE(flui.created_at, 0) AS created_at,
    COALESCE(flui.client_oldest_version, '') AS client_oldest_version,
    COALESCE(flui.client_latest_version, '') AS client_latest_version
FROM
    feature
LEFT OUTER JOIN feature_last_used_info AS flui ON
    feature.id = flui.feature_id AND
    feature.environment_id = flui.environment_id
LEFT OUTER JOIN experiment ON
    feature.id = experiment.feature_id AND
    feature.environment_id = experiment.environment_id
%s %s %s