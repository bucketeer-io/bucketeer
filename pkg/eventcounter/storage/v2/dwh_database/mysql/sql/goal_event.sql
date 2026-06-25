WITH grouped_by_user_evaluation AS (
    SELECT
        user_id,
        variation_id,
        COUNT(id) as event_count,
        IFNULL(SUM(value), 0) as value_sum
    FROM
        goal_event
    WHERE
        timestamp BETWEEN ? AND ?
        AND environment_id = ?
        AND goal_id = ?
        AND feature_id = ?
        AND feature_version = ?
    GROUP BY
        user_id,
        variation_id
),
cap_level AS (
    -- Winsorization threshold: the configurable percentile (?, an integer in
    -- [1,100]) of per-user value_sum pooled across all variations. NTILE(100)
    -- divides the sorted rows into 100 buckets; MAX of buckets ≤ the percentile
    -- gives the cap. Values above it are whales whose outsized spend would
    -- otherwise make the Normal posterior overconfident (Kohavi, Tang & Xu,
    -- Trustworthy Online Controlled Experiments, 2020, §4).
    SELECT MAX(value_sum) AS cap
    FROM (
        SELECT
            value_sum,
            NTILE(100) OVER (ORDER BY value_sum) AS pct
        FROM grouped_by_user_evaluation
    ) ranked
    WHERE pct <= ?
),
capped_by_user AS (
    SELECT
        u.user_id,
        u.variation_id,
        u.event_count,
        LEAST(u.value_sum, IFNULL(c.cap, 0)) AS value_sum
    FROM grouped_by_user_evaluation u
    CROSS JOIN cap_level c
)
SELECT
    variation_id as variationID,
    COUNT(DISTINCT user_id) as goalUser,
    SUM(event_count) as goalTotal,
    SUM(value_sum) as goalValueTotal,
    AVG(value_sum) as goalValueMean,
    IFNULL(VAR_SAMP(value_sum), 0) as goalValueVariance
FROM
    capped_by_user
GROUP BY
    variation_id
ORDER BY
    variation_id
