DELETE FROM
    account_v2
WHERE
    email = $1
    AND organization_id = $2
