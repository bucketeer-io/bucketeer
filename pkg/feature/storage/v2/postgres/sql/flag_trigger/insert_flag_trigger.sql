INSERT INTO flag_trigger(id,
                         feature_id,
                         environment_id,
                         type,
                         action,
                         description,
                         trigger_count,
                         last_triggered_at,
                         token,
                         disabled,
                         created_at,
                         updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
