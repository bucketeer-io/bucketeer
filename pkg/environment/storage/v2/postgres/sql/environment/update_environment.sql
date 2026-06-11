UPDATE environment_v2
SET
    name = $1,
    description = $2,
    archived = $3,
    require_comment = $4,
    created_at = $5,
    updated_at = $6,
    auto_archive_enabled = $7,
    auto_archive_unused_days = $8,
    auto_archive_check_code_refs = $9
WHERE
    id = $10
