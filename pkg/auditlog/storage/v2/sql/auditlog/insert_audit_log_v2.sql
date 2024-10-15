INSERT INTO audit_log (
    id,
    timestamp,
    entity_type,
    entity_id,
    type,
    event,
    editor,
    options,
    environment_namespace,
    entity_data,
    previous_entity_data
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
