SELECT
    tag.id,
    tag.name,
    tag.created_at,
    tag.updated_at,
    tag.entity_type,
    tag.environment_id,
    env.name as environment_name
FROM
    tag
JOIN
    environment_v2 env ON tag.environment_id = env.id
WHERE
    tag.name = $1 AND
    tag.environment_id = $2 AND
    tag.entity_type = $3
