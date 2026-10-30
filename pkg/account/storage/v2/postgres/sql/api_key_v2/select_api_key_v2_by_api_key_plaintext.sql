-- Fallback for keys created before api_key_hash existed. The api_key column is
-- still accurate for every row, so a row that was backfilled between the digest
-- lookup and this one is matched as well.
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
    api_key = $1 AND
    environment_id = $2
