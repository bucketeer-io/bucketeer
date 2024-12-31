SELECT
    id,
    name,
    role,
    disabled,
    created_at,
    updated_at,
    description,
    api_key,
    maintainer
FROM
    api_key
WHERE
    api_key = ? AND
    environment_id = ?