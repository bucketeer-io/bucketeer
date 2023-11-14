UPDATE flag_triggers
SET disabled = 1,
    updated_at = ?
WHERE id = ?
  AND environment_namespace = ?;