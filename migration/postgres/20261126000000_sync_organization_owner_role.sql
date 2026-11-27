-- Give every organization's owner_email account the OWNER role (RFC 0049, owner audit).
-- An ADMIN or higher role carries no environment roles, as in AccountV2.ChangeOrganizationRole.
-- Organizations whose owner_email has no account are left for a manual check.
UPDATE account_v2
SET organization_role = 3,
    environment_roles = '[]',
    updated_at = EXTRACT(EPOCH FROM NOW())::BIGINT
FROM organization
WHERE organization.id = account_v2.organization_id
    AND organization.owner_email = account_v2.email
    AND account_v2.organization_role <> 3;
