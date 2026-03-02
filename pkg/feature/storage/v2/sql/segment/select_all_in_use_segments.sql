SELECT 
    seg.id AS segment_id,
    seg.environment_id AS environment_id,
    seg.updated_at AS segment_updated_at
FROM segment AS seg
WHERE seg.deleted = 0
    AND EXISTS (
        SELECT 1 FROM feature AS ft 
        WHERE ft.environment_id = seg.environment_id 
            AND ft.rules LIKE CONCAT('%', seg.id, '%')
            AND ft.deleted = 0
    )
ORDER BY seg.environment_id, seg.id;
