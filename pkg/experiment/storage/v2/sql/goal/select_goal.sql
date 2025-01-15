SELECT
    id,
    name,
    description,
    connection_type,
    archived,
    deleted,
    created_at,
    updated_at,
    CASE
        WHEN (
            SELECT
                COUNT(1)
            FROM
                experiment ex1
            WHERE
                ex1.environment_id = ? AND
                ex1.goal_ids LIKE concat("%", goal.id, "%")
            ) > 0 THEN TRUE
        ELSE FALSE
    END AS is_in_use_status,
    (
        select
            CONCAT('[', GROUP_CONCAT(JSON_OBJECT('id', ex2.id, 'name', ex2.name)), ']')
        from experiment ex2
        where json_contains(ex2.goal_ids, concat('"', goal.id, '"'), '$')
    ) as experiments
FROM
    goal
WHERE
    id = ? AND
    environment_id = ?