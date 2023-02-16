## Overview

#### Procedure

1. Read YAML file.
2. Create a new instance for each cloud service such as AWS, GCP, and Azure.
3. Call common method from the new instance. The caller doesn't have to be conscious of which cloud service is used. (Dependency Injection)

#### Implementation 

1. Decide the interface(common method name) and abstract the around middleware implementation.
2. (if needed,) Implement GCP middleware part with SDK library.
3. Implement AWS middleware part with SDK library.
4. Decide YAML interface, then create Helm file and implement reading YAML file.

First, we'll start implementing AWS. Later, we'll support Azure.

## Self-host Bucketeer

[Prototype of design for YAML file](./utils/sample.yml)

### Controversial Topic

#### How to configure database info?

1.
```yaml
params:
  host: localhost
  port: 5432
  usernameFile: path/to/usernameFile # or env?
  passwordFile: path/to/passwordFile # or env?
```

2.
```yaml
params:
  url: mysql://<username>:<password>@<host>:<port>/<db_name>
  pemFile: path/to/pemfile # for IAM
```

#### How to configure secret info?

1
```yaml
params:
  usernameFile: path/to/usernameFile # or env?
  passwordFile: path/to/passwordFile # or env?
```

2
```yaml
params:
  username: ${POSTGRES_USER} # or env?
  password: ${POSTGRES_PW} # or env?
```

#### Should we use the specific cloud's name or the commonalized name?

1. Use the specific cloud's name

```yaml
kms:
  type: cloudKms
  params:
    keyName: keyName
```

```yaml
kms:
  type: awsKms
  params:
    keyID: keyID
    region: # for AWS
```

```yaml
messageService:
  type: pubsub 
  metricsEvent:
    topic: bucketeer-metrics-events
    subscription:
      metricsEventPersister: bucketeer-metrics-events-metrics-event-persister
```

```yaml
messageService:
  type: sns/sqs
  metricsEvent:
    topic: bucketeer-metrics-events
    queue:
      metricsEventPersister: bucketeer-metrics-events-metrics-event-persister
```

2. Use the commonalized name

```yaml
kms:
  type: cloudKms
  params:
    keyName: keyName
```

```yaml
kms:
  type: awsKms
  params:
    keyName: keyName
    region: # for AWS
```

```yaml
messageService:
  type: pubsub 
  metricsEvent:
    publisher: bucketeer-metrics-events
    subscriber:
      metricsEventPersister: bucketeer-metrics-events-metrics-event-persister
```

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

