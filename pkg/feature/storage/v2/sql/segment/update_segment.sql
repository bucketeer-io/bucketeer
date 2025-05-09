		UPDATE 
			segment
		SET
			name = ?,
			description = ?,
			rules = ?,
			created_at = ?,
			updated_at = ?,
			version = ?,
			deleted = ?,
			included_user_count = ?,
			excluded_user_count = ?,
			status = ?
		WHERE
			id = ? AND
			environment_id = ?