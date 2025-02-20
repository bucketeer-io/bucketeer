SELECT
    COUNT(
            CASE
                WHEN (
                        SELECT
                            COUNT(1)
                        FROM
                            experiment
                        WHERE
                            environment_id = ? AND
                            goal_ids LIKE concat("%%", goal.id, "%%")
                    ) %s
                END
    )
FROM
    goal
%s
