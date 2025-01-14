SELECT 
    id,
    name,
    created_at,
    updated_at,
    entity_type,
    environment_id
FROM
    tag
WHERE
    id = ? AND
    environment_id = ?
