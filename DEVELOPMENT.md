## Setup development environment

It's recommended to use Github dev container to setup the development environment.
The dev container is based on Ubuntu 20.04 and contains all the necessary tools to build and run the project.
The dev container is also configured to use the latest version of the project.

There are two ways to setup the development environment by using dev container:

1. Use the dev container directly from GitHub Codespaces
2. Build the dev container locally using VSCode `Dev Containers` extension

### Use the dev container directly from GitHub

Using the dev container directly from Github is the easiest way to setup the development environment. There are
configuration file for dev container in the project. Github will automatically build the dev container and run it in the
cloud.
But it may need to make a billing of the dev container if you use it frequently. \
You can find more detail about the billing of GitHub dev
container [here](https://docs.github.com/en/github/developing-online-with-codespaces/about-billing-for-codespaces).

1. Open the [bucketeer project](https://github.com/bucketeer-io/bucketeer) in GitHub
2. Click the `Code` button and select `Open with Codespaces`
3. Select `New codespace` and click `Create codespace` (We set a minimal machine type `Basic (4 vCPU, 8 GB RAM)` for the
   dev container)
4. Wait for the dev container to be ready

### Setup Minikube and services that Bucketeer depends on

1. Open the terminal (in project root directory) in the dev container and run the following command to install the
   dependencies:

```shell
make update-repos local-deps 
```

This command will install the Golang packages that Bucketeer depends on.

2. Setup minikube and services that Bucketeer depends on:

```shell
make start-minikube
```

> Note: When you restart the minikube cluster, you will need to use `make start-minikube` to start the
> cluster, do not use `minikube start` directly.

This command will set up minikube and services that Bucketeer depends on:

* MySQL
* Redis
* Google Pub/Sub (Emulator)
* Google Big Query (Emulator)
* Hashicorp Vault

It will add 2 hosts to `/etc/hosts` that point to the minikube IP address:

* `api-gateway.bucketeer.org` for API Gateway Service
* `web-gateway.bucketeer.org` for Web Gateway Service

Additionally, this command will:

* Initialize the Hashicorp Vault Transit Engine

* Create tables for Google Big Query (Emulator)

3. Generate the certificates for local development:

```shell
# generate TLS certificate and OAuth key in minikube
make generate-tls-certificate generate-oauth service-cert-secret oauth-key-secret
# create Github token in minikube
GITHUB_TOKEN="your token"  make generate-github-token
# generate service token that gRPC service uses to authenticate
ISSUER=https://accounts.google.com \
EMAIL="your email" \
OAUTH_KEY_PATH=/workspaces/bucketeer/tools/dev/cert/oauth-private.pem \
SERVICE_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/service-token \
make generate-service-token
```

The commands above will generate TLS certificate, OAuth key, GitHub token and service token in minikube. And the service
token will be used to authenticate the gRPC service (we will use this token in Helm Charts `values.dev.yaml` later).

4. Build the project:

```shell
# build all go applications
make build-go
# build all docker images, make sure to set the TAG environment variable
TAG=test make build-docker-images
# load the docker images into minikube, make sure to set the TAG environment variable
TAG=test make minikube-load-images
```

5. Deploy the project:

Make sure to set the service token in `values.dev.yaml` file before deploying the project.

```yaml
serviceToken:
  secret:
  # set the service token here     
  token:
```

For example, we will deploy the `backend` service:

```shell
# deploy the backend service (in project root directory)
helm install backend manifests/bucketeer/charts/backend/ --values manifests/bucketeer/charts/backend/values.dev.yaml
```

As you can see, we use the `values.dev.yaml` file to override the default values in `values.yaml` file. And we use the
service token (`/workspaces/bucketeer/tools/dev/cert/service-token`) that we generated in step 3 to authenticate the
gRPC service.


> Pro-tip: You can use `make deploy-service-to-minikube` to deploy services.
> For example, we will deploy the `backend` service:
> ```shell
> SERVICE=backend make deploy-service-to-minikube
> ```
> This command will deploy the `backend` service to minikube. And it will use the service token that we generated in
*step 3* to authenticate the gRPC service. \
> But make sure you are in the project root directory when you run this command.
>
> Also, you can use `make deploy-all-services-to-minikube` to deploy all services to minikube, if you don't want to
> deploy services one by one.

### Deploy Bucketeer in one command
As you can see, there are many steps to setup the development environment after start minukube. You can use the following command to complete all the steps above in one command:

```shell
EMAIL="your email" \ 
ISSUER=https://accounts.google.com \
TAG=test \ 
make deploy-bucketeer
```

### Run the project e2e tests

* Create api key for e2e tests

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.org \
GATEWAY_URL=api-gateway.bucketeer.org \
WEB_GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/service-token \
API_KEY_PATH=/workspaces/bucketeer/apitoken ENVIRONMENT_NAMESPACE=e2e \
API_KEY_ROLE=SDK_CLIENT \
make create-api-key 
```
> Note: The `API_KEY_ROLE` is the role of the api key, you can set it to `SDK_CLIENT` or `SDK_SERVER`.

* Run e2e tests

```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.org \
GATEWAY_URL=api-gateway.bucketeer.org \
WEB_GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
GATEWAY_CERT_PATH=/workspaces/bucketeer/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/service-token \
API_KEY_PATH=/workspaces/bucketeer/apitoken ENVIRONMENT_NAMESPACE=e2e \
make e2e
```