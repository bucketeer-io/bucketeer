## Setup development environment

It's recommended to use Github dev container to setup the development environment.
The dev container is based on Ubuntu 20.04 and contains all the necessary tools to build and run the project.
The dev container is also configured to use the latest version of the project.

There are two ways to setup the development environment by using dev container:

1. Use the dev container directly from Github Codespaces
2. Build the dev container locally using VSCode Remote - Containers extension

### Use the dev container directly from Github

Using the dev container directly from Github is the easiest way to setup the development environment. There are
configuration file for dev container in the project. Github will automatically build the dev container and run it in the
cloud.
But it may need to make a billing of the dev container if you use it frequently. \
You can find more detail about the billing of Github dev
container [here](https://docs.github.com/en/github/developing-online-with-codespaces/about-billing-for-codespaces).

1. Open the [bucketeer project](https://github.com/bucketeer-io/bucketeer) in Github
2. Click the `Code` button and select `Open with Codespaces`
3. Select `New codespace` and click `Create codespace` (We set a minimal machine type `Basic (4 vCPU, 8 GB RAM)` for the
   dev container)
4. Wait for the dev container to be ready

### Setup Minikube and services that Bucketeer depends on

1. Open the terminal (in project root directory) in the dev container and run the following command to install the
   dependencies:

```shell
make tidy-deps local-deps
```

This command will install the Golang packages that Bucketeer depends on.

2. Setup minikube and services that Bucketeer depends on:

```shell
make setup-minikube
```

> Note: If you setup the minikube cluster for the first time, next time you can just run `make start-minikube` to start the cluster.

This command will setup minikube and services that Bucketeer depends on:

* MySQL
* Redis
* Google Pub/Sub (Emulator)
* Google Big Query (Emulator)
* Hashicorp Vault

Also it will add 2 hosts to `/etc/hosts` that point to the minikube IP address:

* `api-gateway.bucketeer.org` for API Gateway Service
* `web-gateway.bucketeer.org` for Web Gateway Service


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

The commands above will generate TLS certificate, OAuth key, Github token and service token in minikube. And the service
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
helm install backend manifests/bucketeer/charts/backend/ --values manifests/bucketeer/charts/backend/****values.dev.yaml
```

As you can see, we use the `values.dev.yaml` file to override the default values in `values.yaml` file. And we use the
service token (`/workspaces/bucketeer/tools/dev/cert/service-token`) that we generated in step 3 to authenticate the
gRPC service.


> Pro tip: You can use `make deploy-service-to-minikube` to deploy services.
> For example, we will deploy the `backend` service:
> ```shell
> SERVICE=backend make deploy-service-to-minikube
> ```
> This command will deploy the `backend` service to minikube. And it will use the service token that we generated in *step 3* to authenticate the gRPC service. \
> But make sure you are in the project root directory when you run this command.