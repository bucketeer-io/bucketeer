INSERT INTO experiment_result
    (id, experiment_id, updated_at, data, environment_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, environment_id) DO UPDATE SET
    experiment_id = EXCLUDED.experiment_id,
    updated_at = EXCLUDED.updated_at,
    data = EXCLUDED.data
