UPDATE admin_subscription SET
    updated_at = $1,
    disabled = $2,
    source_types = $3,
    recipient = $4,
    name = $5
WHERE
    id = $6
