SELECT
    sub.id,
    sub.created_at,
    sub.updated_at,
    sub.disabled,
    sub.source_types,
    sub.recipient,
    sub.name,
    sub.feature_flag_tags,
    sub.environment_id,
    env.name as environment_name
FROM
    subscription sub
JOIN
    environment_v2 env ON sub.environment_id = env.id
