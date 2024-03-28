SELECT
    id,
    created_at,
    updated_at,
    disabled,
    source_types,
    recipient,
    name
FROM
    subscription
WHERE
    id = ? AND
    environment_namespace = ?
