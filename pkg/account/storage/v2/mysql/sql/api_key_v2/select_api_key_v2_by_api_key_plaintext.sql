-- Fallback for keys created before api_key_hash existed. The api_key column is
-- still accurate for every row, so a row that was backfilled between the digest
-- lookup and this one is matched as well. The comparison is binary: the table
-- collation is case-insensitive, while a key must match byte for byte, exactly
-- as it does through its digest.
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
    api_key COLLATE utf8mb4_bin = ? AND
    environment_id = ?
