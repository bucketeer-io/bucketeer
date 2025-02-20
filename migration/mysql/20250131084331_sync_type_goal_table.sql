# populate goal with experiment type
UPDATE
    goal
SET connection_type = 1
WHERE
    connection_type = 0;
