UPDATE flag_triggers
SET description  = ?,
    updated_at = ?
WHERE id = ?
  AND environment_namespace = ?
