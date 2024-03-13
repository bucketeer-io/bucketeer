UPDATE
    api_key
SET
    name = %s,
    role = %d,
    disabled = %t,
    created_at = %d,
    updated_at = %d
WHERE
    id = %s AND
    environment_namespace = %s