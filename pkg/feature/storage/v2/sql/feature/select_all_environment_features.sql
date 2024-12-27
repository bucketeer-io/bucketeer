SELECT 
    env.id AS environment_id,
    ft.id AS feature_id,
    ft.name AS feature_name,
    ft.description AS feature_description,
    ft.enabled AS feature_enabled,
    ft.archived AS feature_archived,
    ft.deleted AS feature_deleted,
    ft.version AS feature_version,
    ft.created_at AS feature_created_at,
    ft.updated_at AS feature_updated_at,
    ft.variation_type AS feature_variation_type,
    ft.variations AS feature_variations,
    ft.targets AS feature_targets,
    ft.rules AS feature_rules,
    ft.default_strategy AS feature_default_strategy,
    ft.off_variation AS feature_off_variation,
    ft.tags AS feature_tags,
    ft.maintainer AS feature_maintainer,
    ft.sampling_seed AS feature_sampling_seed,
    ft.prerequisites AS feature_prerequisites
FROM feature ft
JOIN environment_v2 env ON ft.environment_id = env.id
ORDER BY env.id, ft.id;
