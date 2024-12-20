SELECT
    sub.id,
    sub.created_at,
    sub.updated_at,
    sub.disabled,
    sub.source_types,
    sub.recipient,
    sub.name,
    env.name as environment_name
FROM
    subscription sub
JOIN environment_v2 env
     ON sub.environment_id = env.id
WHERE
    sub.id = ? AND
    sub.environment_id = ?
