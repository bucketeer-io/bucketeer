INSERT INTO environment_v2 (
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
    auto_archive_check_code_refs
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)