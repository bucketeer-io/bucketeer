SELECT
    id,
    name,
    role,
    disabled,
    created_at,
    updated_at
FROM
    api_key
WHERE
    id = ? AND
    (
        environment_namespace = ? or
        environment_id = ?
    )