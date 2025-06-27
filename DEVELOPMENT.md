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

* MySQL
* Redis
* Google Pub/Sub (Emulator)
* Google Big Query (Emulator)

```shell
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

```shell
make deploy-bucketeer
```

If you need to deploy a single service, you can do as follows.

```shell
# Deploy the backend service (in the project root directory)
helm install backend manifests/bucketeer/charts/backend/ --values manifests/bucketeer/charts/backend/values.dev.yaml
```

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
make start-docker-compose
```

2. Or start with specific versions:
```shell
# Use a different Bucketeer version for all services
export BUCKETEER_VERSION=v1.4.0
make start-docker-compose

# Or customize individual service versions
export BUCKETEER_WEB_VERSION=v1.4.0
export BUCKETEER_API_VERSION=v1.3.0
make start-docker-compose
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

### Services and Ports

The Docker Compose setup includes:

- **MySQL**: Port 3306 (Database)
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

# View logs for specific service
docker-compose -f docker-compose/docker-compose.yml logs -f web

# Restart a specific service
docker-compose -f docker-compose/docker-compose.yml restart api
```

# Running E2E Tests

## For Minikube Setup

To run E2E tests you must create API Keys for Server and Client SDKs.
Please note that you only need to create them once.

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

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
GATEWAY_URL=api-gateway.bucketeer.io \
WEB_GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/service-token \
API_KEY_PATH=/workspaces/bucketeer/tools/dev/cert/api_key_client \
API_KEY_SERVER_PATH=/workspaces/bucketeer/tools/dev/cert/api_key_server \
ENVIRONMENT_ID=e2e \
ORGANIZATION_ID=default \
make e2e
```

### Delete E2E data

```shell
make delete-dev-container-mysql-data
```

## For Docker Compose Setup

When using Docker Compose instead of Minikube, you can run E2E tests with modified endpoints:

### Create API Keys for Docker Compose

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
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
SERVICE_TOKEN_PATH=$PWD/tools/dev/cert/service-token \
API_KEY_NAME="e2e-test-$(date +%s)-server" \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_server \
API_KEY_ROLE=SDK_SERVER \
ENVIRONMENT_ID=e2e \
make create-api-key
```

### Run E2E Tests Against Docker Compose

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
GATEWAY_URL=api-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
GATEWAY_PORT=443 \
SERVICE_TOKEN_PATH=$PWD/tools/dev/cert/service-token \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_client \
API_KEY_SERVER_PATH=$PWD/tools/dev/cert/api_key_server \
ENVIRONMENT_ID=e2e \
ORGANIZATION_ID=default \
make e2e
```
