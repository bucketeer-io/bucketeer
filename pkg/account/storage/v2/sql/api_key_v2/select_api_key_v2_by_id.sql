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
    description,
    api_key,
    maintainer
FROM
    api_key
WHERE
    id = ? AND
    environment_id = ?