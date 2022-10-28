## Overview

#### Procedure

1. Read YAML file.
2. Create a new instance for each cloud service such as AWS, GCP, and Azure.
3. Call common method from the new instance. The caller doesn't have to be conscious of which cloud service is used. (Dependency Injection)

##### Example for YAML file

```yaml
cloud: 'aws'
project: 'foobar'
...
```

#### Implementation 

1. Decide the interface(common method name) and abstract the around middleware implementation.
2. (if needed,) Implement GCP middleware part with SDK library.
3. Implement AWS middleware part with SDK library.
4. Decide YAML interface, then create Helm file and implement reading YAML file.

First, we'll start implementing AWS. Later, we'll support Azure.

## Self-host Bucketeer

### Comparision

**PostHog**

```yaml
cloud: 'aws'
ingress:
    hostname: <your-hostname>
    nginx:
        enabled: true
cert-manager:
    enabled: true
```

https://posthog.com/docs/self-host/deploy/aws#chart-configuration

```yaml
cloud: 'gcp'
ingress:
    hostname: <your-hostname>
```

https://posthog.com/docs/self-host/deploy/gcp#chart-configuration

They use Helm and above YAML file.

https://github.com/PostHog/charts-clickhouse

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

https://pipecd.dev/docs/installation/install-controlplane/#using-firestore-and-gcs

They use Helm and above YAML file.

```yaml
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

https://docs.growthbook.io/self-host#installation

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

|                                    | GCP           | AWS                                                          | Azure                                                        |
| ---------------------------------- | ------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
|                                    | Cloud Pub/Sub | Amazon Simple Notification Service（SNS）、Amazon Simple Queueing Service（SQS） | Azure Service Bus Messaging                                  |
| total size of Publish Request      | 10MB          | 0.25MB                                                       | 0.25MB for [Standard tier](https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-premium-messaging)<br/>100 MB for [Premium tier](https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-premium-messaging). |
| number of messages per transaction | 1,000         | 10                                                           | 100                                                          |
| Topic                              | 10,000        | 10,0000                                                      | 10,000                                                       |
| Subscription per topic             | 10,000        | 12,500,000                                                   | 2,000                                                        |
| Expiration period in topic         | 31 days       | 4 weeks                                                      | 14 days                                                      |
| Expiration period in subscription  | 31 days       | 14 days (default is 4 days)                                  | 14 days                                                      |
| API throttling(Tokyo region)       |               | 1500 transactions per second                                 |                                                              |
| At-Least-Once                      | ◯             | ◯                                                            | ◯                                                            |
| SLA(%)                             | >=99.9        | >=99.9                                                       | >=99.9                                                       |



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
* We don't need high performance DB.

Cons

* It takes time longer than 1 to finish this task.

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
