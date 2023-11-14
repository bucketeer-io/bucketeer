SELECT id,
       feature_id,
       environment_namespace,
       type,
       `action`,
       description,
       trigger_times,
       last_triggered_at,
       uuid,
       disabled,
       deleted,
       created_at,
       updated_at
from flag_triggers
where id = ?
  AND environment_namespace = ?