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

## Set up Minikube and services that Bucketeer depends on

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

## Deploy Bucketeer

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

## Run E2E tests

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
