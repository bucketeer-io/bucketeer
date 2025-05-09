		SELECT
			id,
			name,
			description,
			rules,
			created_at,
			updated_at,
			version,
			deleted,
			included_user_count,
			excluded_user_count,
			status,
			(
				SELECT 
					GROUP_CONCAT(id)
				FROM 
					feature
				WHERE
					environment_id = ? AND
					rules LIKE concat("%%", segment.id, "%%")
			) AS feature_ids
		FROM
			segment
		WHERE
			id = ? AND
			environment_id = ?