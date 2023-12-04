UPDATE flag_triggers
SET trigger_times  = ?,
    last_triggered_at = ?,
    updated_at = ?
WHERE id = ?
  AND environment_namespace = ?
