UPDATE environment_v2
SET
    name = ?,
    description = ?,
    archived = ?,
    require_comment = ?,
    created_at = ?,
    updated_at = ?,
    auto_archive_enabled = ?,
    auto_archive_unused_days = ?,
    auto_archive_require_no_code_refs = ?
WHERE
    id = ?