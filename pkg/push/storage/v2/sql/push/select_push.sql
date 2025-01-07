SELECT
    p.id,
    p.fcm_service_account,
    p.tags,
    p.deleted,
    p.name,
    p.created_at,
    p.updated_at,
    p.disabled,
    p.environment_id,
    env.name AS environment_name
FROM
    push as p
JOIN
    environment_v2 env ON p.environment_id = env.id
WHERE
    id = ? AND
    environment_id = ? 