SELECT
    COUNT(DISTINCT feature.id)
FROM
    feature
LEFT OUTER JOIN
    experiment
ON
    feature.id = experiment.feature_id AND
    feature.environment_id = experiment.environment_id
%s