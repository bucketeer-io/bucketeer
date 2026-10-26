SELECT
    sampling_seed
FROM
    holdout_group
WHERE
    environment_id = $1
    AND assignment_type = 0
    AND sampling_seed <> ''
ORDER BY
    created_at ASC
LIMIT 1
