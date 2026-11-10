SELECT
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
    feature.variation_value_schema,
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
    feature.last_used_at
FROM
    feature
