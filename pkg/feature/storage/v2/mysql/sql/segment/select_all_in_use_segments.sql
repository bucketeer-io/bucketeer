SELECT 
    seg.id AS segment_id,
    seg.environment_id AS environment_id,
    seg.updated_at AS segment_updated_at,
    seg.rules AS segment_rules
FROM segment AS seg
WHERE seg.deleted = 0
    AND (
        EXISTS (
            SELECT 1 FROM feature AS ft 
            WHERE ft.environment_id = seg.environment_id 
                AND ft.rules LIKE CONCAT('%', seg.id, '%')
                AND ft.deleted = 0
        )
        OR EXISTS (
            SELECT 1 FROM holdout_group AS hg
            WHERE hg.environment_id = seg.environment_id
                AND hg.segment_id = seg.id
                AND hg.status IN (?, ?, ?)
                AND hg.deleted = 0
        )
    )
ORDER BY seg.environment_id, seg.id;
