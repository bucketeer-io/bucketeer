UPDATE flag_triggers
SET uuid      = ?,
    updated_at = ?
WHERE id = ?
  AND environment_namespace = ?
