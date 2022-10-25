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

First, we'll start implementing AWS. Later, we'll support Azure.

## Current middlewares we use

| Currently used middlewares | Usage                                                 |
| -------------------------- | ----------------------------------------------------- |
| Cloud Pub/Sub & Bigtable   | Event-driven-architecture                             |
| Cloud SQL & Memorystore    | Almost Util data such as feature, segmentUser, apiKey |
| Druid (GCS) & Kafka        | Calculator                                            |
| Cloud KMS                  | Webhook                                               |

### Cloud Pub/Sub & Bigtable

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



### Druid (GCS) & Kafka

WIP



### Cloud KMS

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

