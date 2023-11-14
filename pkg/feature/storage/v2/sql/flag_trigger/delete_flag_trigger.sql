UPDATE flag_triggers
SET deleted   = 1,
    updated_at = ?
WHERE id = ?
  AND environment_namespace = ?