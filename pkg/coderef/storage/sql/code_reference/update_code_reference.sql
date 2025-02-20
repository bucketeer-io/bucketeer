UPDATE code_reference
SET
    file_path = ?,
    file_extension = ?,
    line_number = ?,
    code_snippet = ?,
    content_hash = ?,
    aliases = ?,
    repository_name = ?,
    repository_owner = ?,
    repository_type = ?,
    repository_branch = ?,
    commit_hash = ?,
    updated_at = ?
WHERE
    id = ?
    AND environment_id = ? 