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
        select count(aor.id)
        from auto_ops_rule aor
        where
            aor.feature_id = feature.id and
            ops_type = 1 and
            aor.deleted = 0
    ) as progressive_rollout_count,
    (
        select count(aor.id)
        from auto_ops_rule aor
        where
            aor.feature_id = feature.id and
            ops_type = 2 and
            aor.deleted = 0
    ) as schedule_count,
    (
        select count(aor.id)
        from auto_ops_rule aor
        where
            aor.feature_id = feature.id and
            ops_type = 3 and
            aor.deleted = 0
    ) as kill_switch_count
FROM feature
LEFT OUTER JOIN experiment ON
    feature.id = experiment.feature_id AND
    feature.environment_id = experiment.environment_id
        %s %s %s