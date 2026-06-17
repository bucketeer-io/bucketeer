SELECT
    COUNT(1)
FROM
    subscription sub
JOIN
    environment_v2 env ON sub.environment_id = env.id
