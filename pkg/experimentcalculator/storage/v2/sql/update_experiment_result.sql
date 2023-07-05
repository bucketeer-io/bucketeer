INSERT INTO experiment_result
    (id, experiment_id, updated_at, data, environment_namespace)
VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY
UPDATE
    experiment_id = VALUES(experiment_id),
    updated_at = VALUES(updated_at),
    data = VALUES(data)