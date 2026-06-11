SELECT
    id,
    experiment_id,
    updated_at,
    data
FROM
    experiment_result
WHERE
    id = $1 AND
    environment_id = $2
