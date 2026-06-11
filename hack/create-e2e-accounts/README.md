# create-e2e-accounts

Bootstraps the accounts used by local development and e2e tests, and generates
an access token for each of them. It upserts rows directly into the
`account_v2` MySQL table and writes the signed access tokens to files.

| Account | Organization (role) | Environment role | Access token |
|---------|---------------------|------------------|--------------|
| `--sys-admin-email` | e2e org (`OWNER`) | — | `--sys-admin-token-output` (system admin) |
| `--org-admin-email` | default org (`ADMIN`) + e2e org (`ADMIN`) | — (org `ADMIN` implies access to every environment in the org) | `--org-admin-token-output` |
| `--env-write-email` | default org (`MEMBER`) | `EDITOR` on the e2e environment (`--e2e-environment-id`) | `--env-write-token-output` |
| `--env-read-email` | default org (`MEMBER`) | `VIEWER` on the e2e environment (`--e2e-environment-id`) | `--env-read-token-output` |

The e2e environment is owned by the default organization, so the editor/viewer
accounts are members of the default org carrying an environment role for the
e2e environment.

Rows are upserted with `INSERT ... ON DUPLICATE KEY UPDATE` (MySQL 8.0
row-alias form) against the composite primary key `(email, organization_id)`,
so re-running against a DB that already has the rows is safe.

The system admin token is scoped to the e2e organization and is minted with
`is_system_admin=true`, so it can call system-admin-only APIs and read across
organizations. The other three tokens are scoped to the default organization
and are **not** system admins: the org admin relies on its organization
`ADMIN` role and the editor/viewer rely on their environment roles, so they
exercise the real RBAC path. None of the tokens is a service token, and all
are minted with a far-future expiry.

## Run Command

```
go run ./hack/create-e2e-accounts create \
  --mysql-user=<MYSQL_USER> \
  --mysql-pass=<MYSQL_PASS> \
  --mysql-host=<MYSQL_HOST> \
  --mysql-port=<MYSQL_PORT> \
  --mysql-db-name=<MYSQL_DB_NAME> \
  --sys-admin-email=<SYS_ADMIN_EMAIL> \
  --org-admin-email=<ORG_ADMIN_EMAIL> \
  --env-write-email=<ENV_WRITE_EMAIL> \
  --env-read-email=<ENV_READ_EMAIL> \
  --default-organization-id=<DEFAULT_ORGANIZATION_ID> \
  --e2e-organization-id=<E2E_ORGANIZATION_ID> \
  --e2e-environment-id=<E2E_ENVIRONMENT_ID> \
  --oauth-key=<OAUTH_PRIVATE_KEY_PATH> \
  --issuer=<ISSUER> \
  --sys-admin-token-output=<SYS_ADMIN_TOKEN_PATH> \
  --org-admin-token-output=<ORG_ADMIN_TOKEN_PATH> \
  --env-write-token-output=<ENV_EDITOR_TOKEN_PATH> \
  --env-read-token-output=<ENV_VIEWER_TOKEN_PATH> \
  --no-profile \
  --no-gcp-trace-enabled
```

### Flags

| Flag | Description |
|------|-------------|
| `--mysql-user` | MySQL user. |
| `--mysql-pass` | MySQL password. |
| `--mysql-host` | MySQL host. |
| `--mysql-port` | MySQL port. |
| `--mysql-db-name` | MySQL database name. |
| `--sys-admin-email` | Email of the system admin account (`OWNER` of the e2e org; its token is a system admin). |
| `--org-admin-email` | Email of the organization admin account (`ADMIN` in the default and e2e orgs). |
| `--env-write-email` | Email of the environment editor account (`MEMBER` of default org, `EDITOR` on the e2e environment). |
| `--env-read-email` | Email of the environment viewer account (`MEMBER` of default org, `VIEWER` on the e2e environment). |
| `--default-organization-id` | ID of the default organization that owns the e2e environment. |
| `--e2e-organization-id` | ID of the e2e organization where the org admin account also gets `ADMIN`. |
| `--e2e-environment-id` | ID of the e2e environment used for the editor/viewer environment roles. |
| `--oauth-key` | Path to the OAuth RSA private key used to sign the access tokens. |
| `--issuer` | Issuer URL set in the generated access tokens (must match the gateway config). |
| `--audience` | OAuth audience set in the generated access tokens (default `bucketeer`). |
| `--sys-admin-token-output` | Path of the file to write the system admin access token. |
| `--org-admin-token-output` | Path of the file to write the org admin access token. |
| `--env-write-token-output` | Path of the file to write the environment editor access token. |
| `--env-read-token-output` | Path of the file to write the environment viewer access token. |

## When to use

- Fresh Minikube / Kubernetes cluster where the MySQL init SQL wasn't applied.
- Re-seeding after deleting the `account_v2` rows.
- Before running e2e tests: `make e2e` does not bootstrap the accounts, so run
  this first (e.g. via `make create-dev-container-e2e-accounts` or
  `make docker-compose-create-e2e-accounts`).

## Create docker image

```
make deps

export PAT=<PERSONAL_ACCESS_TOKEN>
export GITHUB_USER_NAME=<GITHUB_USER_NAME>
export TAG=<TAG>

make docker-build
make docker-push
```

Personal Access Token needs to have `write:packages` permission.
