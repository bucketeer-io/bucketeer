INSERT INTO code_reference (
    id,
    feature_id,
    file_path,
    line_number,
    code_snippet,
    content_hash,
    aliases,
    repository_name,
    repository_owner,
    repository_type,
    repository_branch,
    commit_hash,
    environment_id,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 