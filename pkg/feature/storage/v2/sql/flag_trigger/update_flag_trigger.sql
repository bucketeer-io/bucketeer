UPDATE flag_triggers
SET feature_id = ?,
    type = ?,
    action = ?,
    description = ?,
    trigger_times = ?,
    last_triggered_at = ?,
    uuid = ?,
    disabled = ?,
    created_at = ?,
    updated_at = ?
WHERE id = ?
  AND environment_namespace = ?
