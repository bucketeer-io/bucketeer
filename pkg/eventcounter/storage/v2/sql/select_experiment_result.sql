SELECT
    id,
    experiment_id,
    updated_at,
    data
FROM
    experiment_result
WHERE
    id = ? AND
    environment_id = ?