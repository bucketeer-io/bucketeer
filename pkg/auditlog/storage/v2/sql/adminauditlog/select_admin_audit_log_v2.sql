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
    admin_audit_log
    %s %s %s