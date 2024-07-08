SELECT
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
FROM
    audit_log
    %s %s %s