UPDATE flag_triggers
SET disabled = 0,
    updated_at = ?
WHERE id = ?
  AND environment_namespace = ?;

