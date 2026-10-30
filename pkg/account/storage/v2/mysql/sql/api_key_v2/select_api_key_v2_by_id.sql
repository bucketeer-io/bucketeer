SELECT
    id,
    name,
    role,
    disabled,
    created_at,
    updated_at,
    description,
    api_key,
    COALESCE(api_key_prefix, '') AS api_key_prefix,
    maintainer,
    last_used_at
FROM
    api_key
WHERE
    id = ? AND
    environment_id = ?
