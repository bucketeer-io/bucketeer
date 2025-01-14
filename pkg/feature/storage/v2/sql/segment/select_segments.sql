SELECT
    seg.id,
    seg.name,
    seg.description,
    seg.rules,
    seg.created_at,
    seg.updated_at,
    seg.version,
    seg.deleted,
    seg.included_user_count,
    seg.excluded_user_count,
    seg.status,
    (
        SELECT
            GROUP_CONCAT(id)
        FROM
            feature as ft
        WHERE
            ft.environment_id = seg.environment_id AND
            ft.rules LIKE concat("%%", seg.id, "%%")
    ) AS feature_ids
FROM
    segment as seg
%s %s %s %s
