SELECT
    COUNT(1)
FROM
    api_key
LEFT JOIN environment_v2 ON api_key.environment_id = environment_v2.id
%s