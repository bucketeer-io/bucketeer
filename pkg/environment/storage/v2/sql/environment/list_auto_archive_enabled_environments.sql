SELECT
    id,
    name,
    url_code,
    description,
    project_id,
    organization_id,
    archived,
    require_comment,
    created_at,
    updated_at,
    auto_archive_enabled,
    auto_archive_unused_days,
    auto_archive_require_no_code_refs
FROM
    environment_v2
WHERE
    auto_archive_enabled = TRUE
ORDER BY
    id ASC
