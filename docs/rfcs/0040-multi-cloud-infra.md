## Overview.

This rfc explains how to support self-hosting on multi cloud services.

First, we'll start implementing AWS. Later, we'll support Azure.

### Other services

**PipeCD**

```yaml
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  stateKey: {RANDOM_STRING}
  datastore:
    type: FIRESTORE or MySQL
    config:
      namespace: pipecd
      environment: dev
      project: {YOUR_GCP_PROJECT_NAME}
      # Must be a service account with "Cloud Datastore User" and "Cloud Datastore Index Admin" roles
      # since PipeCD needs them to creates the needed Firestore composite indexes in the background.
      credentialsFile: /etc/pipecd-secret/firestore-service-account
  filestore:
    type: GCS or AWS S3 or MINIO
    config:
      bucket: {YOUR_BUCKET_NAME}
      # Must be a service account with "Storage Object Admin (roles/storage.objectAdmin)" role on the given bucket
      # since PipeCD need to write file object such as deployment log file to that bucket.
      credentialsFile: /etc/pipecd-secret/gcs-service-account
```

```sh
$ helm install pipecd oci://ghcr.io/pipe-cd/chart/pipecd --version v0.39.0 --namespace={NAMESPACE} \
  --set-file config.data=path-to-control-plane-configuration-file \
  --set-file secret.encryptionKey.data=path-to-encryption-key-file \
  --set-file secret.firestoreServiceAccount.data=path-to-service-account-file \
  --set-file secret.gcsServiceAccount.data=path-to-service-account-file
```

https://pipecd.dev/docs/installation/install-controlplane/#using-firestore-and-gcs

**Growth Book**

To open web console, docker-compose.yml is used. Users can't use the feature of feature flag at this time.

```yml
# docker-compose.yml
version: "3"
services:
  mongo:
    image: "mongo:latest"
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=password
  growthbook:
    image: "growthbook/growthbook:latest"
    ports:
      - "3000:3000"
      - "3100:3100"
    depends_on:
      - mongo
    environment:
      - MONGODB_URI=mongodb://root:password@mongo:27017/
    volumes:
      - uploads:/usr/local/src/app/packages/back-end/uploads
volumes:
  uploads:
```

To use the feature of feature flag, Growth Book loads `/usr/local/src/app/config/config.yml`.

```yml
datasources:
  warehouse:
    type: postgres # or "redshift" or "mysql" or "clickhouse"
    name: Main Warehouse
    # Connection params (different for each type of data source)
    params:
      host: localhost
      port: 5432
      user: root
      password: ${POSTGRES_PW} # use env for secrets
      database: growthbook
...
```

As an alternative plan, users can register the configuration in web console as follows:

![growth-book-dashboard](./images/0039-image6.png)
![growth-book-dashboard2](./images/0039-image7.png)

https://docs.growthbook.io/self-host#installation

**mastodon**

Users configure the following `./.env.production`. Then, run docker-compose.

```text

# Redis
# -----
REDIS_HOST=localhost
REDIS_PORT=6379

# PostgreSQL
# ----------
DB_HOST=/var/run/postgresql
DB_USER=mastodon
DB_NAME=mastodon_production
DB_PASS=
DB_PORT=5432

# Elasticsearch (optional)
# ------------------------
ES_ENABLED=true
ES_HOST=localhost
ES_PORT=9200
# Authentication for ES (optional)
ES_USER=elastic
ES_PASS=password

```

https://github.com/mastodon/mastodon/blob/main/.env.production.sample

**FeatureHub**

There are several deployment options for running FeatureHub. Please visit https://docs.featurehub.io/featurehub/latest/installation.html for further information.

FeatureHub supports both [kubernetes](https://github.com/featurehub-io/featurehub-helm) and [docker-compose](https://github.com/featurehub-io/featurehub-install).
Users configure the following application.properties files.

```text
db.url=jdbc:postgresql://db:5432/featurehub
db.username=featurehub
db.password=featurehub
db.connections=10
nats.urls=nats://nats:4222
dacha1.enabled=false
dacha2.enabled=true
```

**PostHog**

Users need to create values.yaml such as follows:

```yaml
cloud: 'aws'
ingress:
  hostname: <your-hostname>
  nginx:
      enabled: true
externalPostgresql:
  # -- External PostgreSQL service host.
  postgresqlHost:
  # -- External PostgreSQL service port.
  postgresqlPort: 5432
```

Then, deploying service as follows:

```console
$ helm repo add posthog https://posthog.github.io/charts-clickhouse/
$ helm repo update
$ helm upgrade --install -f values.yaml --timeout 30m --create-namespace --namespace posthog posthog posthog/posthog --wait --wait-for-jobs --debug
```

https://github.com/PostHog/charts-clickhouse

### Conclusion

Since Bucketeer uses Kubernetes, using YAML file and Helm fits into our cases.

## Self-host Bucketeer

Basically, we'll start support self-host by updating each values.yaml.

In the future, we might make it possible to deploy by updating only one YAML file as follows:

[Prototype of design for YAML file](./utils/sample.yml)

### Controversial topic

#### 1. Which middleware requires its own implementation?

| Currently used middlewares | Require to implement in each cloud                           |
| -------------------------- | ------------------------------------------------------------ |
| Cloud Pub/Sub              | Yes                                                          |
| Cloud SQL                  | No (If we support IAM authentication, we need to use AWS SDK) |
| BigQuery                   | Yes                                                          |
| Cloud KMS                  | Yes                                                          |
| MemoryStore                | No                                                           |

##### 2. How to configure values.yaml?

###### Cloud platform

```yaml
cloudPlatform: gcp
```

###### Cloud SQL

```yaml
relationalDatabase:
  type: mysql
  password:
  port:
```

##### MemoryStore

```yaml
inMemoryDataStore:
  password:
  port:
```

###### Cloud KMS

```yaml
keyManagement:
  gcp:
    kms:
      keyId: 
  aws:
    kms:
      ...
```

###### Cloud Pub/Sub

```yaml
messaging:
  gcp:
    pubsub:
      topic:
      subscription:
      pullerNumGoroutines: 5
      pullerMaxOutstandingMessages: "1000"
      pullerMaxOutstandingBytes: "1000000000"
```

###### BigQuery

```yaml
bigData:
  gcp:
    bigQuery:
      dataset:
```

##### 3. Design of directory 

###### pkg directory

The pkg directory will be as follows:

```
pkg
├── storage
│   ├── bigdata
│       ├── athena
│       └── bigquery
│   ├── inmemory
│       └── redis
│   └── rdb
│       ├── alloydb
│       ├── aurora
│       ├── mysql
│       └── postgres
├── kms
│   ├── aws
│   └── gcp
├── messaging
│   ├── pubsub
│   ├── sns
│   └── sqs

...
```

###### cmd package in each micro service

The cmd package in each service will have aws.go and gcp.go as follows:

```
pkg/auth/cmd
└── server
    ├── aws.go
    ├── gcp.go
    └── server.go
```

##### 4. Changes in deployment.yaml

deployment.yaml in each each service will be changed as follows:

```yaml
{{- if eq .Values.env.relationalDatabase.type "mysql" }}
  - name: BUCKETEER_EVENT_COUNTER_MYSQL_USER
    value: "{{ .Values.env.relationalDatabase.user }}"
{{- end }}
```

Also, we'll determine which cloud is used by checking `cloudPlatform` field in values.yaml as follows:

```yaml
{{- if eq .Values.env.cloudPlatform "gcp" }}
  ...
{{- else eq .Values.env.cloudPlatform "aws" }}
  ...
{{- end }}
```
