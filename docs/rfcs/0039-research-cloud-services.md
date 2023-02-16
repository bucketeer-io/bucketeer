## Overview

## Current middlewares we use

| Currently used middlewares | Usage                                                 |
| -------------------------- | ----------------------------------------------------- |
| Cloud Pub/Sub              | Event-driven-architecture                             |
| Cloud SQL                  | Almost Util data such as feature, segmentUser, apiKey |
| BigQuery                   | A/B test                                              |
| Cloud KMS                  | Webhook                                               |
| MemoryStore                | Feature Flag such as auto ops and evaluation count    |

### Cloud Pub/Sub

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

**Comparison table about order guarantee**

|                              | Cloud Pub/Sub | Amazon Simple Notification Service（SNS）、Amazon Simple Queueing Service（SQS） | Azure Service Bus Messaging |
| ---------------------------- | ------------- | ------------------------------------------------------------ | --------------------------- |
| Order guarantee(optional)    | Yes           | Yes                                                          | Yes                         |
| Topic                        |               | 10,0000 -> 1,000                                             |                             |
| Subscription                 |               | 12,500,000 -> 100                                            |                             |
| API throttling(Tokyo region) |               | 1500 -> 300 transactions per second (3000 if high throughput) |                             |
| messages only once           |               | ◯                                                            |                             |

### Cloud SQL

#### Comparison Table

|      | GCP        | AWS                                                     | Azure                                                      |
| ---- | ---------- | ------------------------------------------------------- | ---------------------------------------------------------- |
|      | Cloud  SQL | Amazon Relational Database Service (RDS), Amazon Aurora | Azure Database for MySQL and Azure Database for PostgreSQL |

**Comparison table about AlloyDB**

|        | GCP     | AWS           | Azure              |
| ------ | ------- | ------------- | ------------------ |
|        | AlloyDB | Amazon Aurora | Azure SQL Database |
| SLA(%) | >=99.99 | >=99.99       | >=99.995           |

### BigQuery

#### Comparison Table

|      | GCP      | AWS                                                      | Azure                   |
| ---- | -------- | -------------------------------------------------------- | ----------------------- |
|      | BigQuery | Amazon Athena, Amazon Redshift, Amazon Redshift Spectrum | Azure Synapse Analytics |

### MemoryStore

#### Comparison Table

|        | GCP         | AWS                | Azure                  |
| ------ | ----------- | ------------------ | ---------------------- |
|        | Memorystore | Amazon ElastiCache | Azure Cache            |
| SLA(%) | >=99.9      | >=99.9             | >=99.9 (from Standard) |

#### Controversial topic
##### 1. Can we configure Memorystore as optional?

Yes. We can configure memory store as a optional in YAML file.

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
