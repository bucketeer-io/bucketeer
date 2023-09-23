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

There is a service called backend in the project. It's the main service of the project which contains all modules of the
project.
You can build the **backend** service by using the following command:

```shell
make vender

make build-backend
```

This command will build the **backend** service and generate the binary file in the `bin` directory.

## Run the project

As the Setup development environment section mentioned, the project is using the following service:

* Google BigQuery
* Google Cloud Pub/Sub
* MySQL
* Redis
* Vault

You can use the following command to get these service:

```shell
kubectl get services
```

You'll see some output like the following:

```shell
NAME                                TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
kubernetes                          ClusterIP   10.96.0.1        <none>        443/TCP             7d22h
localenv-bq                         ClusterIP   10.110.0.177     <none>        9060/TCP,9050/TCP   7d22h
localenv-mysql                      ClusterIP   10.96.134.38     <none>        3306/TCP            7d22h
localenv-mysql-headless             ClusterIP   None             <none>        3306/TCP            7d22h
localenv-pubsub                     ClusterIP   10.104.148.143   <none>        8089/TCP            7d22h
localenv-redis-headless             ClusterIP   None             <none>        6379/TCP            7d22h
localenv-redis-master               ClusterIP   10.111.175.18    <none>        6379/TCP            7d22h
localenv-vault                      ClusterIP   10.104.132.32    <none>        8200/TCP,8201/TCP   7d22h
localenv-vault-agent-injector-svc   ClusterIP   10.106.151.189   <none>        443/TCP             7d22h
localenv-vault-internal             ClusterIP   None             <none>        8200/TCP,8201/TCP   7d22h
```

So you can use these service to run the project. For example, you can start the **backend** service by the following
steps:

1. Generate certificate file for the **backend** service

```shell
cd tools/cert

make generate-tls-certificate
make generate-oauth

make service-cert-secret
make oauth-key-secret

GITHUB_TOKEN=${enter_your_token} make generate-github-token
```

These commands will generate the certificate file for the **backend** service and create the secret in the kubernetes,
so we can use the certificate file and the secret in the **backend** service chart.

2. Generate the service token for the **backend** service

```shell
ISSUER=https://accounts.google.com \
EMAIL=xxx@gmail.com \
OAUTH_KEY_PATH=/workspaces/bucketeer/tools/dev/cert/oauth-private.pem \
SERVICE_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/service-token \
make generate-service-token
```

This command will generate the service token in `tools/dev/cert` directory. You can find
the token in the `service-token` file, and you will use it in the service chart.

3. Set the environment variables in the helm chart `manifests/bucketeer/charts/backend/values.yaml`

```yaml
image:
  repository: ghcr.io/bucketeer-io/bucketeer-backend
  pullPolicy: IfNotPresent

fullnameOverride: "backend"

namespace: default

env:
  cloudService: hcv
  profile: false
  gcpTraceEnabled: false
  bucketeerTestEnabled: true
  bigqueryEmulatorHost: localenv-bq.default.svc.cluster.local:9050
  pubsubEmulatorHost: localenv-pubsub.default.svc.cluster.local:8089
  mysqlMigrationUser: bucketeer
  mysqlMigrationPass: bucketeer
  project: bucketeer-test
  mysqlUser: bucketeer
  mysqlPass: bucketeer
  mysqlHost: localenv-mysql-headless.default.svc.cluster.local
  mysqlPort: 3306
  mysqlDbName: bucketeer
  persistentRedis:
    serverName: backend
    addr: localenv-redis-headless.default.svc.cluster.local:6379
    poolMaxIdle: 25
    poolMaxActive: 25
  nonPersistentRedis:
    serverName: backend
    addr: localenv-redis-headless.default.svc.cluster.local:6379
    poolMaxIdle: 25
    poolMaxActive: 25
  bigQueryDataSet: bucketeer
  bigQueryDataLocation: bucketeer
  domainTopic: domain
  bulkSegmentUsersReceivedTopic: bulk-segment-users-received
  accountService: localhost:9001
  authService: localhost:9001
  environmentService: localhost:9001
  experimentService: localhost:9001
  featureService: localhost:9001
  healthCheckServicePort: 8000
  accountServicePort: 9091
  authServicePort: 9092
  auditLogServicePort: 9093
  autoOpsServicePort: 9094
  environmentServicePort: 9095
  eventCounterServicePort: 9096
  experimentServicePort: 9097
  featureServicePort: 9098
  migrateMysqlServicePort: 9099
  notificationServicePort: 9100
  pushServicePort: 9101
  metricsPort: 9002
  timezone: UTC
  emailFilter:
  githubUser: bucketeer
  githubMigrationSourcePath: /tmp/migration
  logLevel: info

tls:
  service:
    secret: bucketeer-service-cert
    cert:
    key:
  issuer:
    secret: bucketeer-service-cert
    cert:

serviceToken:
  secret:
  token:
oauth:
  key:
    secret: bucketeer-oauth-key
    public:
  clientId: bucketeer
  clientSecret: oauth-client-secret
  redirectUrls: https://google.com
  issuer: https://accounts.google.com

webhook:
  baseURL: http://localhost:9000
  kmsResourceName: vault

```

As you can see, we set the environment variables in the helm chart. We use the already created secret and certificate
file in the kubernetes.

> Note: This is just an example, but it can be used in many cases. You can modify the environment variables according to
> your needs.

4. Install the **backend** service

```shell
cd manifests/bucketeer
helm install backend charts/backend/ --set global.image.tag=${tag}
```

You should replace the `${tag}` with the tag of the **backend** service image.

5. Check the status of the **backend** service

```shell
# check the status of the **backend** service
kubectl get pods
# output
NAME                                            READY   STATUS    RESTARTS   AGE
backend-7bb7b99df-zqmwd                         2/2     Running   0          15m
localenv-bq-7f45dd7fb9-qnhvr                    1/1     Running   0          101m
localenv-mysql-0                                1/1     Running   0          101m
localenv-pubsub-596c76779b-tzk8p                1/1     Running   0          101m
localenv-redis-master-0                         1/1     Running   0          101m
localenv-vault-0                                1/1     Running   0          101m
localenv-vault-agent-injector-79bf89fc6-lhksl   1/1     Running   0          101m

# show the logs of the **backend** service
kubectl logs backend-7bb7b99df-zqmwd -c backend
{"severity":"INFO","eventTime":1693838947.3376513,"logger":"bucketeer-backend.server.metrics","caller":"metrics/metrics.go:119","message":"Run started","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.3377233,"logger":"bucketeer-backend.server","caller":"cli/app.go:163","message":"Running bucketeer-backend.server","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.3380814,"logger":"bucketeer-backend.server.http","caller":"rest/server.go:96","message":"Rest server is running on 8000","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4232283,"logger":"bucketeer-backend.server.rpc-server.account-server","caller":"rpc/server.go:121","message":"Running on 9091","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4232364,"logger":"bucketeer-backend.server.rpc-server.auth-server","caller":"rpc/server.go:121","message":"Running on 9092","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4295552,"logger":"bucketeer-backend.server.rpc-server.audit-log-server","caller":"rpc/server.go:121","message":"Running on 9093","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4315414,"logger":"bucketeer-backend.server.rpc-server.push-server","caller":"rpc/server.go:121","message":"Running on 9101","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4324594,"logger":"bucketeer-backend.server.rpc-server.event-counter-server","caller":"rpc/server.go:121","message":"Running on 9096","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.432471,"logger":"bucketeer-backend.server.rpc-server.environment-server","caller":"rpc/server.go:121","message":"Running on 9095","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4333467,"logger":"bucketeer-backend.server.rpc-server.migrate-mysql-server","caller":"rpc/server.go:121","message":"Running on 9099","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4333768,"logger":"bucketeer-backend.server.rpc-server.feature-server","caller":"rpc/server.go:121","message":"Running on 9098","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4339676,"logger":"bucketeer-backend.server.rpc-server.experiment-server","caller":"rpc/server.go:121","message":"Running on 9097","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4341455,"logger":"bucketeer-backend.server.rpc-server.notification-server","caller":"rpc/server.go:121","message":"Running on 9100","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
{"severity":"INFO","eventTime":1693838947.4344642,"logger":"bucketeer-backend.server.rpc-server.auto-ops-server","caller":"rpc/server.go:121","message":"Running on 9094","serviceContext":{"service":"bucketeer-backend.server","version":"-"}}
```

So until now, the **backend** service is running successfully.