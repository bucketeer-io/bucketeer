WITH auto_ops_counts AS (
    SELECT
        aor.feature_id,
        aor.environment_id,
        COUNT(CASE WHEN aor.ops_type = 1 THEN 1 END) AS progressive_rollout_count,
        COUNT(CASE WHEN aor.ops_type = 2 THEN 1 END) AS schedule_count,
        COUNT(CASE WHEN aor.ops_type = 3 THEN 1 END) AS kill_switch_count
    FROM auto_ops_rule aor
    WHERE aor.deleted = FALSE
    GROUP BY aor.feature_id, aor.environment_id
)
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
    COALESCE(auto_ops_counts.progressive_rollout_count, 0) AS progressive_rollout_count,
    COALESCE(auto_ops_counts.schedule_count, 0) AS schedule_count,
    COALESCE(auto_ops_counts.kill_switch_count, 0) AS kill_switch_count,
    feature.last_used_at
FROM
    feature
LEFT JOIN auto_ops_counts ON
    feature.id = auto_ops_counts.feature_id AND
    feature.environment_id = auto_ops_counts.environment_id
LEFT JOIN LATERAL (
    SELECT e.id
    FROM experiment e
    WHERE e.feature_id = feature.id
    AND e.environment_id = feature.environment_id
    AND e.deleted = FALSE
    LIMIT 1
) experiment ON TRUE
