SELECT
    updated_at
FROM
    holdout_group
WHERE
    environment_id = $1
    AND status IN ($2, $3)
ORDER BY
    updated_at DESC
LIMIT 1
