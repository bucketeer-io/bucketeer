SELECT
    id,
    name,
    description,
    enabled,
    archived,
    deleted,
    evaluation_undelayable,
    ttl,
    version,
    created_at,
    updated_at,
    variation_type,
    variations,
    targets,
    rules,
    default_strategy,
    off_variation,
    tags,
    maintainer,
    sampling_seed,
    prerequisites,
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
FROM
    feature
    %s %s %s