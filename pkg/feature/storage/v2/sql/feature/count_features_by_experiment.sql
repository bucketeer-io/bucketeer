SELECT
    COUNT(DISTINCT feature.id)
FROM
    feature
LEFT JOIN experiment ON
    feature.id = experiment.feature_id AND
    feature.environment_id = experiment.environment_id