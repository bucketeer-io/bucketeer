SELECT
    project.*,
    COALESCE(COUNT(DISTINCT environment_v2.id), 0) AS environment_count,
    COALESCE(COUNT(DISTINCT feature.id), 0) AS feature_count
FROM
    project
        LEFT JOIN
    environment_v2 ON project.id = environment_v2.project_id
        LEFT JOIN
    feature ON environment_v2.id = feature.environment_id
%s 
GROUP BY
    project.id
%s %s