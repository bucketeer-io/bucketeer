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

## Current middlewares we use

| Currently used middlewares | Usage                                                 |
| -------------------------- | ----------------------------------------------------- |
| Cloud Pub/Sub & Bigtable   | Event-driven-architecture                             |
| Cloud SQL & Memorystore    | Almost Util data such as feature, segmentUser, apiKey |
| Druid (GCS) & Kafka        | Calculator                                            |
| Cloud KMS                  | Webhook                                               |

### Cloud Pub/Sub & Bigtable

![event-pipeline](./images/0039-image1.png)

#### Comparison Table

|                                    | GCP                                                          | AWS                                                          | Azure                                                        |
| ---------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
|                                    | Cloud Pub/Sub                                                | Amazon Simple Notification Service（SNS）、Amazon Simple Queueing Service（SQS） | Azure Service Bus Messaging                                  |
| total size of Publish Request      | 10MB                                                         | 0.25MB                                                       | 0.25MB for [Standard tier](https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-premium-messaging)<br/>100 MB for [Premium tier](https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-premium-messaging). |
| number of messages per transaction | 1,000                                                        | 10                                                           | 100                                                          |
| Topic                              | 10,000                                                       | 10,0000                                                      | 10,000                                                       |
| Subscription per topic             | 10,000                                                       | 12,500,000                                                   | 2,000                                                        |
| Expiration period in topic         | 31 days                                                      | 4 weeks                                                      | 14 days                                                      |
| message retention period           | 7 days                                                       | 14 days (default is 4 days)                                  | 14 days                                                      |
| API throttling(Tokyo region)       |                                                              | 1500 transactions per second                                 |                                                              |
| At-Least-Once                      | ◯                                                            | ◯                                                            | ◯                                                            |
| Retry logic about subscription     | redelivers every subsequent message with the same ordering key, including acknowledged messages | redelivers only the messages                                 | redelivers only the messages                                 |
| SLA(%)                             | >=99.95                                                      | >=99.9                                                       | >=99.9                                                       |



|      | GCP      | AWS             | Azure           |
| ---- | -------- | --------------- | --------------- |
|      | Bigtable | Amazon DynamoDB | Azure Cosmos DB |
|      |          |                 |                 |

#### Controversial topic

##### 1. Can we stop using Bigtable?

We use Bigtable to store userEvaluation because PubSub doesn't guarantee order in the default. However, we can stop it if we use the order guarantee feature.

**Comparison table about order guarantee**

|                              | Cloud Pub/Sub | Amazon Simple Notification Service（SNS）、Amazon Simple Queueing Service（SQS） | Azure Service Bus Messaging |
| ---------------------------- | ------------- | ------------------------------------------------------------ | --------------------------- |
| Order guarantee(optional)    | Yes           | Yes                                                          | Yes                         |
| Topic                        |               | 10,0000 -> 1,000                                             |                             |
| Subscription                 |               | 12,500,000 -> 100                                            |                             |
| API throttling(Tokyo region) |               | 1500 -> 300 transactions per second (3000 if high throughput) |                             |
| messages only once           |               | ◯                                                            |                             |

**Important notice about enabling order guarantee**

* GCP
  * Possibility of increasing latency
  * Increase the number of duplicates
  * We have to design the ordering key so that granularity is not increased.
* AWS
  * By default, it limits up to 300 messages per second.
    * However, it can be configured to handle 6000 messages by enabling throughput mode.
* Azure
  * We have to design the session id so that granularity is not increased.

In conclusion, we can enable an order guarantee because we can set the key for each evaluation(userEvaluations + goalEvaluations). In short, we can design the key with fine granularity.

**Conversion plan from existing subscriber to ordering subscriber**

1. Create the new pubsub topic(bucketeer-xxx-evaluation-goal-events) in terraform.

2. Create the new subscription(bucketeer-xxx-evaluation-goal-events-event-persister) with turning on ordering feature and target topic.

3. Implement for enabling message ordering
  * Publisher
    * https://cloud.google.com/pubsub/docs/publisher#using-ordering-keys
  * Subscriber
    * https://cloud.google.com/pubsub/docs/ordering#enabling_message_ordering
  * Stores all events into dummy table such as dummy_evaluation_event, dummy_goal_event

4. Check if all messages are correctly stored into dummy tables correctly.

5. Move evaluations from BigTable into RDB.

6. Remove evaluation event persister and goal event persister.

7. Stop using Bigtable


### Cloud SQL & Memorystore

![event-pipeline](./images/0039-image3.png)

#### Comparison Table

|      | GCP        | AWS                                                     | Azure                                                      |
| ---- | ---------- | ------------------------------------------------------- | ---------------------------------------------------------- |
|      | Cloud  SQL | Amazon Relational Database Service (RDS), Amazon Aurora | Azure Database for MySQL and Azure Database for PostgreSQL |

|        | GCP         | AWS                | Azure                  |
| ------ | ----------- | ------------------ | ---------------------- |
|        | Memorystore | Amazon ElastiCache | Azure Cache            |
| SLA(%) | >=99.9      | >=99.9             | >=99.9 (from Standard) |

#### Controversial topic

##### 1. Can we configure Memorystore as optional?

Yes. We can configure memory store as a optional in YAML file.

### Druid (GCS) & Kafka

![data-pipeline](./images/0039-image4.png)

We use Druid as a relay DB and usual DB(fething data directory) and Kafka as an intereface for Druid. The amount of data is a huge and we have to handle them as a high performance when fething data directory.

#### Controversial topic

##### 1. Can we stop using Druid & Kafka?

We have the following problems when using Druid & Kafka.

* It's hard to maintainance self-hosted service.
* It's hard to solve the problem when something such as error occurs.

Therefore, We have the following options:

**1.  Replace with managed service**

|        | GCP     | AWS           | Azure              |
| ------ | ------- | ------------- | ------------------ |
|        | AlloyDB | Amazon Aurora | Azure SQL Database |
| SLA(%) | >=99.99 | >=99.99       | >=99.995           |

Pros

* If we can use AlloyDB, there is a possibility that we can use single DB.

Cons

* AlloyDB may not match with our use cases. For example, it takes time longer than Druid.

**2.  Preprocess evaluation data** 

![data-pipeline2](./images/0039-image5.png)

Pros

* This architecture can be consistent with Calculator.
* EventCounter service can access to only single DB.
* We don't need high performance DB.

Cons

* It takes time longer than 1 to finish this task.

##### Conclusion

We decided not to conclude this topic. Instead, we decided to store Evaluation Event to another DB, too. Therefore we can divide features into single feature flag service and A/B test service. 

Because Evaluation Event is a large data, we need to switch MySQL and PostgreSQL. Since PostgreSQL is ORDBMS, we need to define table as follows:

MySQL

```mysql
CREATE TABLE IF NOT EXISTS `feature` (
  `id` VARCHAR(255) NOT NULL,
  `tags` JSON NOT NULL,
  ...
  PRIMARY KEY (`id`, `environment_namespace`)
);
CREATE TABLE IF NOT EXISTS `tag` (
  `id` VARCHAR(255) NOT NULL,
  `created_at` BIGINT(20) NOT NULL,
  `updated_at` BIGINT(20) NOT NULL,
  `environment_namespace` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`, `environment_namespace`)
);
```

PostgreSQL

```postgresql
CREATE TABLE IF NOT EXISTS "feature" (
  "id" VARCHAR(255) NOT NULL,
  "tags" tags NOT NULL,
  ...
  PRIMARY KEY ("id", "environment_namespace")
);

INSERT INTO feature VALUES ('id', ROW('id', 1667370510, 1667370510, 'production'), ...);
...
```



### Cloud KMS

![webhook](./images/0039-image2.png)

#### Comparison Table

|        | GCP       | AWS                              | Azure           |
| ------ | --------- | -------------------------------- | --------------- |
|        | Cloud KMS | AWS Key Management Service (KMS) | Azure Key Vault |
| SLA(%) | >=99.99   | >=99.999                         | >=99.99         |

## Ref

Order guarantee

https://cloud.google.com/pubsub/docs/ordering#console

https://medium.com/google-cloud/google-cloud-pub-sub-ordered-delivery-1e4181f60bc8

https://aws.amazon.com/jp/about-aws/whats-new/2022/10/amazon-sqs-increased-throughput-quota-fifo-high-throughput-ht-mode-6000-transactions-per-second-tps/

https://learn.microsoft.com/en-us/azure/service-bus-messaging/message-sessions

https://devblogs.microsoft.com/premier-developer/ordering-messages-in-azure-service-bus/

https://christina04.hatenablog.com/entry/gcp-cloud-pubsub-ordering-key-concern



Quotas

https://cloud.google.com/pubsub/quotas

https://docs.aws.amazon.com/general/latest/gr/sns.html

https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/quotas-messages.html

https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quotas

SLA

https://aws.amazon.com/messaging/sla/

https://cloud.google.com/pubsub/sla

https://azure.microsoft.com/ja-jp/support/legal/sla/service-bus/v1_1/



https://azure.microsoft.com/ja-jp/support/legal/sla/cache/v1_1/

https://aws.amazon.com/elasticache/sla/

https://cloud.google.com/memorystore/sla



https://aws.amazon.com/kms/sla/

https://cloud.google.com/kms/sla

https://azure.microsoft.com/ja-jp/updates/akv-sla-raised-to-9999/



https://azure.microsoft.com/ja-jp/support/legal/sla/azure-sql-database/v1_8/



https://aws.amazon.com/rds/aurora/sla/

https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues-understanding-logic.html

