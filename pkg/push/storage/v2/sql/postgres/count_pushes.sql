SELECT
    COUNT(1)
FROM
    push
JOIN
    environment_v2 env ON push.environment_id = env.id