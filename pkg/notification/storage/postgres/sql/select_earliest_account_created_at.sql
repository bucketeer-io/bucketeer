SELECT
    COALESCE(MIN(created_at), 0)
FROM
    account_v2
WHERE
    email = $1
