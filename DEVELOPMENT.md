## Setup development environment

It's recommended to use Github dev container to setup the development environment.
The dev container is based on Ubuntu 20.04 and contains all the necessary tools to build and run the project.
The dev container is also configured to use the latest version of the project.

There are two ways to setup the development environment by using dev container:

1. Use the dev container directly from Github
2. Build the dev container locally

### Use the dev container directly from Github

Using the dev container directly from Github is the easiest way to setup the development environment. There are
configuration file for dev container in the project. Github will automatically build the dev container and run it in the
cloud.
But it may need to make a billing of the dev container if you use it frequently. \
You can find more detail about the billing of Github dev
container [here](https://docs.github.com/en/github/developing-online-with-codespaces/about-billing-for-codespaces).

1. Open the [bucketeer project](https://github.com/bucketeer-io/bucketeer) in Github
2. Click the `Code` button and select `Open with Codespaces`
3. Select `New codespace` and click `Create codespace` (it's recommended to choose 4core/8GB RAM machine)
4. Wait for the dev container to be ready

First, open the terminal in the dev container and run the following command to install the dependencies:

```shell
make tidy-deps local-deps
```

This command will install the dependencies of the project and generate the `go.sum` file.

Then, you can setup the other service by using helm chart:

```shell
minikube start --memory=4g
cd manifests/localenv
helm install localenv .
```

This command will install the other service in the local kubernetes cluster (minikube). You can check the status of the
service by using the following command:

```shell
kubectl get pods
```

You'll see some output like the following:

```shell
localenv-bq-68f679b667-d68j4                     1/1     Running   2 (2m16s ago)   3d14h
localenv-mysql-0                                 1/1     Running   2 (2m16s ago)   3d14h
localenv-pubsub-7c5bf796cd-5d7ns                 1/1     Running   2 (2m16s ago)   3d14h
localenv-redis-master-0                          1/1     Running   2 (2m16s ago)   3d14h
localenv-vault-0                                 1/1     Running   2 (24h ago)     3d14h
localenv-vault-agent-injector-54c848cd44-88m6d   1/1     Running   4 (98s ago)     3d14h
```

As you can see, the project is using the following service:

* Google BigQuery
* Google Cloud Pub/Sub
* MySQL
* Redis
* Vault

Once all the service is ready, you can start to develop the project.

### Build the dev container locally

Also you can build the dev container locally. This way is more flexible than using the dev container directly from
Github. You can modify the configuration of the dev container and build it locally.

#### Prerequisites:

* VSCode
* Docker

Then you can follow the steps below to build the dev container locally:

1. Clone the project
2. Open the project in VSCode
3. Install
   the [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
   extension, if you haven't already. This extension lets you use a Docker container as a full-featured development
4. Command Shift P (Mac) or Ctrl Shift P (Windows/Linux) to open the Command Palette and type `Dev Container: Rebuild
   and Reopen in Container` and select it. This will build the dev container and open the project in the dev container.

## Build the project

There is a service called backend in the project. It's the main service of the project which contains all moudles of the
project.
You can build the backend by using the following command:

```shell
make build-backend
```

This command will build the backend and generate the binary file in the `bin` directory.