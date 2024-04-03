SELECT
    id,
    created_at,
    updated_at,
    disabled,
    source_types,
    recipient,
    name
FROM
    subscription
    %s %s %s
