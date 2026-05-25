INSERT INTO admin_audit_log (
    id,
    timestamp,
    entity_type,
    entity_id,
    type,
    event,
    editor,
    options,
    entity_data,
    previous_entity_data
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
