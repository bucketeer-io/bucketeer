SELECT
    hg.id,
    hg.name,
    hg.status,
    hgf.feature_id,
    hgf.holdout_variation_id
FROM
    holdout_group_feature hgf
    INNER JOIN holdout_group hg
        ON hg.id = hgf.holdout_group_id
        AND hg.environment_id = hgf.environment_id
WHERE
    hgf.environment_id = ?
    AND hgf.feature_id = ?
    AND hgf.removed_at = 0
    AND hg.deleted = ?
    AND hg.status IN (?, ?, ?)
