# Set up the DEV Container

It's recommended that Development Container (Dev Container) be used to set up the development environment.
The dev container is based on Ubuntu 20.04 and contains all the necessary tools to build and run the project.
It is also configured to use the latest version of the project.

There are two ways to set up the development environment by using a dev container:

1. Use the dev container directly from GitHub Codespaces
2. Build the dev container locally using the VSCode `Dev Containers` extension

## Use the dev container directly from GitHub

Using the dev container directly from GitHub is the easiest way to set up the development environment. There are
configuration files for the dev container in the project. GitHub will automatically build the dev container and run it in the
cloud.
But it may need to make a billing of the dev container if you use it frequently. \
You can find more details about the billing of GitHub dev
container [here](https://docs.github.com/en/github/developing-online-with-codespaces/about-billing-for-codespaces).

1. Open the [bucketeer project](https://github.com/bucketeer-io/bucketeer) on GitHub
2. Click the `Code` button and select `Open with Codespaces`
3. Select `New codespace` and click `Create codespace` (We set a minimal machine type `Basic (4 vCPU, 8 GB RAM)` for the
   dev container)
4. Wait for the dev container to be ready

# Local Development Setup

You can set up Bucketeer locally using one of two methods:

1. **Minikube Setup** (Kubernetes-based)
2. **Docker Compose Setup** (Container-based)

Choose the method that best fits your development environment and requirements.

## Option 1: Minikube Setup

### Set up Minikube and services that Bucketeer depends on

The following command will set up the Minikube and services that Bucketeer depends on:

- MySQL
- PostgresQL (optional)
- Redis
- Google Pub/Sub (Emulator)
- Google Big Query (Emulator - optional)

**Note:** Postgres and BigQuery are optional in Minikube. They are detected automatically by checking the dataWarehouse configuration in the `values.dev.yaml`.

```shell
# Default (Postgres & Big Query disabled)
make start-minikube

# Enable Postgres/Big Query
make start-minikube
```

**Note:** When you restart the Minikube cluster, you must use `make start-minikube` to start it. Do not use `minikube start` directly.

It will add 2 hosts to `/etc/hosts` that point to the minikube IP address:

* `api-gateway.bucketeer.io` for API Gateway Service
* `web-gateway.bucketeer.io` for Web Gateway Service

Additionally, this command will:

- Create tables for MySQL
- Create tables for Google Big Query (Emulator)

### Deploy Bucketeer

The following command will deploy all the Bucketeer services at once.

**Note:** Postgres and BigQuery are optional in Minikube. They are detected automatically by checking the dataWarehouse configuration in the `values.dev.yaml`.

```shell
make deploy-bucketeer
```

If you need to deploy a single service, you can do as follows.

```shell
# Deploy the backend service (in the project root directory)
helm install backend manifests/bucketeer/charts/backend/ --values manifests/bucketeer/charts/backend/values.dev.yaml
```

**Note:** You can switch between data warehouses (MySQL, PostgreSQL, BigQuery) but remember to update the `values.dev.yaml` file to match the data warehouse you are using as the events persister and web service must use same event store service.

**Note:** We use the `values.dev.yaml` file to override the default values in `values.yaml` file.

## Option 2: Docker Compose Setup

As an alternative to the Minikube setup, you can use Docker Compose to run Bucketeer locally. This approach is simpler and doesn't require Kubernetes knowledge.

All Docker Compose related files are organized in the `docker-compose/` directory. See `docker-compose/README.md` for detailed documentation including how to customize Docker image versions.

### Prerequisites

- Docker and Docker Compose installed
- Git repository cloned
- `/etc/hosts` file updated with the following entries:
  ```
  127.0.0.1 web-gateway.bucketeer.io
  127.0.0.1 api-gateway.bucketeer.io
  ```

### Quick Start with Docker Compose

1. Start all services with default versions:
```shell
make docker-compose-up
```

2. Or start with specific versions:
```shell
# Use a different Bucketeer version for all services
export BUCKETEER_VERSION=v1.4.0
make docker-compose-up

# Or customize individual service versions
export BUCKETEER_WEB_VERSION=v1.4.0
export BUCKETEER_API_VERSION=v1.3.0
make docker-compose-up
```

The start command will:
- Set up the required directories and configuration files
- Generate development certificates if they don't exist
- Start all Bucketeer services using Docker Compose

3. Check service status:
```shell
make docker-compose-status
```

4. View logs:
```shell
make docker-compose-logs
```

5. Stop services:
```shell
make docker-compose-down
```

6. Create MySQL event tables for data warehouse functionality:
```shell
make docker-compose-create-mysql-event-tables
```

### Services and Ports

The Docker Compose setup includes:

- **MySQL**: Port 3306 (Database)
- **PostgresQL**: Port 5432 (datawarehouse option)
- **Redis**: Port 6379 (Cache and pub/sub)
- **Nginx**: Ports 80 (HTTP) & 443 (HTTPS) for routing to:
  - **Web Service**: Admin UI and internal APIs via `web-gateway.bucketeer.io`
  - **API Gateway**: SDK endpoints via `api-gateway.bucketeer.io`
- **Batch Service**: Background processing (internal service)
- **Subscriber Service**: Event processing (internal service)

The subscriber service uses JSON configuration files located in the `docker-compose/config/subscriber-config/` directory:
- `subscribers.json`: Main subscriber configurations
- `onDemandSubscribers.json`: On-demand subscriber configurations
- `processors.json`: Processor configurations
- `onDemandProcessors.json`: On-demand processor configurations

### Additional Commands

```shell
# Clean up all containers, networks, and volumes
make docker-compose-clean

# View logs for specific service (Docker Compose v2)
docker compose -f docker-compose/compose.yml logs -f web
# Or for Docker Compose v1
docker-compose -f docker-compose/compose.yml logs -f web

# Restart a specific service (Docker Compose v2)
docker compose -f docker-compose/compose.yml restart api
# Or for Docker Compose v1  
docker-compose -f docker-compose/compose.yml restart api
```

# Running Unit Tests

Before running unit tests, ensure that the httpstan container is running. The experiment package unit tests depend on httpstan for Bayesian analysis.

### Start httpstan

```shell
make start-httpstan
```

This command is idempotent - it will:
- Skip if the container is already running
- Start the existing container if it was stopped
- Create and start a new container if it doesn't exist

### Run Unit Tests

```shell
make test-go
```

### Stop httpstan (Optional)

When you're done running tests:

```shell
make stop-httpstan
```

# Running E2E Tests

## For Minikube Setup

To run E2E tests you must create API Keys for Server and Client SDKs.
Please note that you only need to create them once.

### Bootstrap the e2e accounts (required before running E2E)

E2E API calls authenticate as four accounts that exercise the real RBAC path:

| Account | Organization (role) | Environment role | Token |
|---------|---------------------|------------------|-------|
| `sysadmin@bucketeer.io` | `e2e` (`OWNER`) | — | system admin |
| `orgadmin@bucketeer.io` | `default` (`ADMIN`) + `e2e` (`ADMIN`) | — | org admin |
| `envwrite@bucketeer.io` | `default` (`MEMBER`) | `EDITOR` on the `e2e` environment | env editor |
| `envread@bucketeer.io` | `default` (`MEMBER`) | `VIEWER` on the `e2e` environment | env viewer |

Bootstrap them — and generate an access token for each under
`tools/dev/cert/` — before running E2E. `make e2e` no longer does this for
you; run it explicitly:

```shell
make create-dev-container-e2e-accounts
```

The target connects to the Minikube MySQL pod, upserts into `account_v2`
with `INSERT ... ON DUPLICATE KEY UPDATE`, and writes the four access tokens
to `tools/dev/cert/{sys-admin,org-admin,env-editor,env-viewer}-token`, so it
is safe to run repeatedly (re-run it after `make delete-dev-container-mysql-data`
or against a cluster where the init SQL never ran).

### Create API keys

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
WEB_GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/service-token \
API_KEY_NAME="e2e-test-$(date +%s)-client" \
API_KEY_PATH=/workspaces/bucketeer/tools/dev/cert/api_key_client \
API_KEY_ROLE=SDK_CLIENT \
ENVIRONMENT_ID=e2e \
make create-api-key
```

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
WEB_GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/service-token \
API_KEY_NAME="e2e-test-$(date +%s)-server" \
API_KEY_PATH=/workspaces/bucketeer/tools/dev/cert/api_key_server \
API_KEY_ROLE=SDK_SERVER \
ENVIRONMENT_ID=e2e \
make create-api-key
```

### Run E2E tests

The tests authenticate with the four access tokens written by
`make create-dev-container-e2e-accounts` (run that first — see above); they no
longer take a service token. Point the run at those token files:

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
GATEWAY_URL=api-gateway.bucketeer.io \
WEB_GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
API_KEY_PATH=/workspaces/bucketeer/tools/dev/cert/api_key_client \
API_KEY_SERVER_PATH=/workspaces/bucketeer/tools/dev/cert/api_key_server \
SYS_ADMIN_ACCESS_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/sys-admin-token \
ORG_ADMIN_ACCESS_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/org-admin-token \
ENV_EDITOR_ACCESS_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/env-editor-token \
ENV_VIEWER_ACCESS_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/env-viewer-token \
ENVIRONMENT_ID=e2e \
ORGANIZATION_ID=default \
make e2e
```

### Delete E2E data

```shell
make delete-dev-container-mysql-data
make delete-mysql-data-warehouse-data
make delete-redis-retry-keys
```

## For Docker Compose Setup

When using Docker Compose instead of Minikube, you can run E2E tests with modified endpoints:

### Bootstrap the e2e accounts (required before running E2E)

Bootstrap the four e2e accounts (`sysadmin@bucketeer.io`,
`orgadmin@bucketeer.io`, `envwrite@bucketeer.io`, `envread@bucketeer.io`) and
their access tokens before running E2E. `make e2e` no longer does this for you;
run it explicitly:

```shell
make docker-compose-create-e2e-accounts
```

The target connects to the Docker Compose MySQL on `localhost:3306`, upserts
into `account_v2` with `INSERT ... ON DUPLICATE KEY UPDATE`, and writes the
four access tokens to
`tools/dev/cert/{sys-admin,org-admin,env-editor,env-viewer}-token`, so it is
safe to run repeatedly (re-run it after `make docker-compose-delete-data` or
any time the `account_v2` rows are wiped).
See [`hack/create-e2e-accounts/README.md`](./hack/create-e2e-accounts/README.md)
for details and the prebuilt Docker image option.

### Create API Keys for Docker Compose

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
WEB_GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=$PWD/tools/dev/cert/service-token \
API_KEY_NAME="e2e-test-$(date +%s)-client" \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_client \
API_KEY_ROLE=SDK_CLIENT \
ENVIRONMENT_ID=e2e \
make create-api-key
```

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
WEB_GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=$PWD/tools/dev/cert/service-token \
API_KEY_NAME="e2e-test-$(date +%s)-server" \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_server \
API_KEY_ROLE=SDK_SERVER \
ENVIRONMENT_ID=e2e \
make create-api-key
```

### Run E2E Tests Against Docker Compose

The tests authenticate with the four access tokens written by
`make docker-compose-create-e2e-accounts` (run that first — see above); they no
longer take a service token. Point the run at those token files:

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
GATEWAY_URL=api-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
GATEWAY_PORT=443 \
WEB_GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_client \
API_KEY_SERVER_PATH=$PWD/tools/dev/cert/api_key_server \
SYS_ADMIN_ACCESS_TOKEN_PATH=$PWD/tools/dev/cert/sys-admin-token \
ORG_ADMIN_ACCESS_TOKEN_PATH=$PWD/tools/dev/cert/org-admin-token \
ENV_EDITOR_ACCESS_TOKEN_PATH=$PWD/tools/dev/cert/env-editor-token \
ENV_VIEWER_ACCESS_TOKEN_PATH=$PWD/tools/dev/cert/env-viewer-token \
ENVIRONMENT_ID=e2e \
ORGANIZATION_ID=default \
make e2e
```
