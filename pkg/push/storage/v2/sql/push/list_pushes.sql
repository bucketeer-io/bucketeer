SELECT
    push.id,
    push.fcm_service_account,
    push.tags,
    push.deleted,
    push.name,
    push.created_at,
    push.updated_at,
    push.disabled,
    push.environment_id,
    env.name AS environment_name
FROM
    push
JOIN
    environment_v2 env ON push.environment_id = env.id