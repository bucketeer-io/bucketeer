SELECT
    environment_v2.id,
    environment_v2.name,
    environment_v2.url_code,
    environment_v2.description,
    environment_v2.project_id,
    environment_v2.organization_id,
    environment_v2.archived,
    environment_v2.require_comment,
    environment_v2.created_at,
    environment_v2.updated_at,
    environment_v2.auto_archive_enabled,
    environment_v2.auto_archive_unused_days,
    environment_v2.auto_archive_check_code_refs,
    COALESCE(COUNT(DISTINCT feature.id), 0) AS feature_count
FROM
    environment_v2
        LEFT JOIN
    feature ON environment_v2.id = feature.environment_id
%s
GROUP BY
    environment_v2.id
%s %s