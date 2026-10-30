-- Fallback for keys created before api_key_hash existed. Only rows that have not
-- been backfilled yet are matched on the api_key column.
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
    api_key_hash IS NULL AND
    api_key = $1 AND
    environment_id = $2
