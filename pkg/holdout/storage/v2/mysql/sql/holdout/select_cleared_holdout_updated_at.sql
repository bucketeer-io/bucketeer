SELECT
    updated_at
FROM
    holdout_group
WHERE
    environment_id = ?
    AND status IN (?, ?)
ORDER BY
    updated_at DESC
LIMIT 1
