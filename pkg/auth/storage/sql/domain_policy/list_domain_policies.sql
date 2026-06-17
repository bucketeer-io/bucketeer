SELECT domain,
    auth_policy,
    enabled,
    created_at,
    updated_at
FROM domain_auth_policy
%s %s %s