SELECT
    id,
    name,
    role,
    disabled,
    api_key,
    maintainer,
    description,
    created_at,
    updated_at,
    description
FROM
    api_key
WHERE
    id = ? AND
    environment_id = ?