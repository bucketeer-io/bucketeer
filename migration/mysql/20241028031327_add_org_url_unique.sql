ALTER TABLE `project`
    DROP KEY `unique_url_code`,
    ADD UNIQUE KEY `unique_organization_url_code` (`organization_id`, `url_code`);