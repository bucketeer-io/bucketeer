# create-localenv-account

Bootstraps an account used by local development and e2e tests by upserting
two rows directly into the `account_v2` MySQL table:

| Organization | Role | Notes |
|--------------|------|-------|
| default org (`--default-organization-id`) | `ADMIN` | Grants org-scoped write permissions so admin APIs (`CreateAPIKey`, etc.) succeed. Org-level ADMIN implicitly grants access to every environment in the org. |
| e2e org (`--e2e-organization-id`) | `OWNER` | Expected to be flagged `system_admin=1`, so any member gets system admin privileges. |

The two rows are upserted with `INSERT ... ON DUPLICATE KEY UPDATE` (using
the MySQL 8.0 row-alias form) against the composite primary key
`(email, organization_id)`. Re-running the command against a DB that already
has the rows is safe — the upsert refreshes `name`, `organization_role`,
`environment_roles`, `disabled`, and `updated_at` so a previously-disabled
account is also re-enabled.

## Run Command

```
go run ./hack/create-localenv-account create \
  --mysql-user=<MYSQL_USER> \
  --mysql-pass=<MYSQL_PASS> \
  --mysql-host=<MYSQL_HOST> \
  --mysql-port=<MYSQL_PORT> \
  --mysql-db-name=<MYSQL_DB_NAME> \
  --email=<ACCOUNT_EMAIL> \
  --default-organization-id=<DEFAULT_ORGANIZATION_ID> \
  --e2e-organization-id=<E2E_ORGANIZATION_ID> \
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
| `--email` | Email of the account to create. |
| `--default-organization-id` | ID of the default organization where the account gets ADMIN role. |
| `--e2e-organization-id` | ID of the e2e (system-admin) organization where the account gets OWNER role. |

## When to use

- Fresh Minikube / Kubernetes cluster where the MySQL init SQL wasn't applied.
- Re-seeding after deleting the `account_v2` rows.
- CI pipelines that bring up a blank DB and need the account bootstrapped
  before running e2e tests.

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
