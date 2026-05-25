SELECT id,
       feature_id,
       environment_id,
       type,
       `action`,
       description,
       trigger_count,
       last_triggered_at,
       token,
       disabled,
       created_at,
       updated_at
from flag_trigger
where token = ?