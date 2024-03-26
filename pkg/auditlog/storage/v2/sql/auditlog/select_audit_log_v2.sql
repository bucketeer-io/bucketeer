SELECT
    id,
    timestamp,
    entity_type,
    entity_id,
    type,
    event,
    editor,
    options
FROM
    audit_log
    %s %s %s