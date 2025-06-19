# Summary

Currently, when users configure rules on the Targeting tab in the console, they need to manually type custom attribute keys. This manual process can lead to typos and misconfigurations, potentially causing incorrect conditions when evaluating end-users.

Since these attributes are sent from the SDK to the server, we can automate this process by generating a list to display on the console. This will improve user experience and reduce configuration errors.

# Solution１ - Using the new PubSub topic for UserAttribute -

The UserData sent in the 'GetEvaluationsRequest' is compared with the UserAttributes stored in Redis, and only if new attributes are found, a PubSub topic for the newly created UserAttribute is published. The Subscriber stores the new UserAttribute in Redis.

1. Extract UserAttribute information from the GetEvaluationsRequest from the SDK
2. Compare with Redis UserAttribute cached data
3. Publish only new attribute information using PubSub
4. Save attributes to Redis with UserAttributeSubscriber.
5. Provide an API for the console to retrieve the attribute list

## Sequence
```mermaid
sequenceDiagram
    participant SDK as SDK/App Client
    participant BackendService as Backend Service
    participant PubSub as Google Cloud Pub/Sub
    participant Subscriber as UserAttributeSubscriber (Cloud Functions/Cloud Run)
    participant UserAttributeStore as Redis (Cache & Persistent Store)

    SDK->>BackendService: 1. User Action / Data Update Request
    activate BackendService
    BackendService->>UserAttributeStore: 2. Get UserAttribute Info
    UserAttributeStore-->>BackendService: 3. UserAttribute Info
    BackendService->>BackendService: 4. Compare with Current UserAttribute Info
    alt New or Updated UserAttribute Info
        BackendService->>BackendService: 5. Extract UserAttribute Info (id, key, value, environment_id, created_at)
        BackendService->>PubSub: 6. Publish UserAttribute Message (JSON with id, key, value, environment_id, created_at)
    else UserAttribute Info is Unchanged
        %% No action needed, or log that data is unchanged
    end
    deactivate BackendService

    PubSub-->>Subscriber: 8. Message Delivered (Push/Pull)
    activate Subscriber
    Subscriber->>Subscriber: 9. Parse Pub/Sub Message (Extracting id, key, value, environment_id, created_at)
    Subscriber->>UserAttributeStore: 10. SET UserAttribute (key: user_attribute:{id}, value: {key, value, environment_id, created_at})
    %% Storing as a Redis Hash or JSON string. This step updates the store, which also serves as the cache.

    alt On success
        UserAttributeStore-->>Subscriber: 11. Success (OK from Redis)
        Subscriber-->>PubSub: 12. Acknowledge Message (ACK)
    else On failure
        UserAttributeStore-->>Subscriber: 11. Failure (e.g., Redis Connection Error, Write Error)
        Subscriber-->>PubSub: 12. Not Acknowledge Message (NACK)
        PubSub-->PubSub: 13. Message Re-delivery / Dead-Letter Queue
    end
    deactivate Subscriber
```
### Topic
- Since 'GetEvaluationsRequest' is always sent by SDK users, it is possible to accurately detect UserAttributes. However, since the request requires low latency, storage operations and the like must be processed in a separate thread as much as possible, which increases costs.
- Creating a topic for UserAttribute will increase your PubSub costs.

# Solution２ - Using existing PubSub topics that contain UserAttributes -

This solution leverages an existing PubSub topic by leveraging the UserAttribute included in the 'EvaluationEvent' sent by the 'RegisterEventsRequest'.

1. Add a new subscription for UserAttribute to the existing Evaluation Event topic.
2. Save attributes to Redis with UserAttributeSubscriber.
3. Provide an API for the console to retrieve the attribute list.


## Sequence
```mermaid
sequenceDiagram
    participant SDK as SDK/App Client
    participant BackendService as Backend Service
    participant PubSub as Google Cloud Pub/Sub
    participant Subscriber as UserAttributeSubscriber (Cloud Functions/Cloud Run)
    participant UserAttributeStore as Redis (Cache & Persistent Store)

    rect rgb(220, 220, 220)
        note over SDK,PubSub: This Sequence that already exists
    SDK->>BackendService: User Action / Send RegisterEventsRequest
    activate BackendService
        BackendService->>PubSub: Publish Evaluation Event Message
    deactivate BackendService
    end
    PubSub-->>Subscriber: Message Delivered (Push/Pull)
    activate Subscriber
    Subscriber->>Subscriber: Parse the Pub/Sub message and extract the UserAttribute.
    Subscriber->>UserAttributeStore: SET new UserAttribute

    note over Subscriber,PubSub: Acking regardless of success or failure
        Subscriber-->>PubSub: Acknowledge Message
    deactivate Subscriber
```
### Topic
- It leverages the existing PubSub topic for EvaluetionEvent, but will increase PubSub costs by 40-50%.
- Development costs are low by utilizing existing sequences.

# Solution3 - Using the existing EvaluationCountEventPersister -

This solution simply leverages the UserAttribute contained in the `EvaluationEvent` sent by the `RegisterEventsRequest` and adds a UserAttribute save operation to the existing data save sequence by leveraging the existing PubSub subscription.

1. Create `UserAttributesCache` for Redis.
2. In the `EvalEvtWriter` used by `EventsPersister` (which updates data related to `EvaluationEvent`), add a process to determine whether the attribute is new and save it using `UserAttributesCache`.
3. Provide an API for the console to retrieve the attribute list.


## Sequence
```mermaid
sequenceDiagram
    participant SDK as SDK/App Client
    participant BackendService as Backend Service
    participant PubSub as Google Cloud Pub/Sub
    participant Subscriber as EvaluationCountEventPersister
    participant UserAttributeStore as Redis

    rect rgb(220, 220, 220)
        note over SDK,Subscriber: This Sequence that already exists
    SDK->>BackendService: User Action / Send RegisterEventsRequest
    activate BackendService
        BackendService->>PubSub: Publish Evaluation Event Message
    deactivate BackendService
    PubSub-->>Subscriber: Message Delivered (Push/Pull)
    activate Subscriber
    end
    Subscriber->>Subscriber: Parse the Pub/Sub message and extract the UserAttribute.
    Subscriber->>UserAttributeStore: SET new UserAttribute

    note over Subscriber,PubSub: Acking regardless of success or failure
        Subscriber-->>PubSub: Acknowledge Message
    deactivate Subscriber
```
### Topic
- Existing PubSub subscriptions will be leveraged, so your PubSub costs will not increase.
- Development costs are low by utilizing existing sequences.

# Conclustion
I adopt **Solution3** because it will not increase development costs or PubSub costs.

# Solution3 Implementation Details

## Cache

- Create `UserAttributesCache` in the cache package
  - Key: string (environment_id:user-attributes)
  - Value: []string (user_attribute_keys)
  - Interface: Put and Get

## Processor
- Add a process to the `EvalEvtWriter` to extract new attributes from the `EvaluationEvent`'s `UserData` and save them in the `UserAttributesCache`.

## API

Add a new API to get UserAttributes in the environment:

```protobuf
message ListUserAttributesRequest {
    string environment_id = 1;
}

message ListUserAttributesResponse {
    repeated string userAttributes = 1;
}
```

Note: Pagination is not implemented for this API.

# Important Considerations

-  Intentionally not implementing user attribute delete API
   - The deleted attribute may be needed again in the future, but there is currently no way to undo the deletion using the console.
   - However, the console takes into account the large number of user attributes by providing incremental search to improve usability.


# Testing

- The e2e test is performed in the following steps:
- Test flow:
  1. Send request via `RegisterEventsRequest`
  2. Wait for processing
  3. Verify attributes via `ListUserAttributes` API

# Release Steps

1. Implement Cache
   - Implement `UserAttributesCacher`.

2.  Implement Processor and API. and verify its operation through e2e test
   - Implement `UserAttributeCache` saving process in Processor
   - Implement `ListUserAttributes` API
   - Test e2e for Dev Env
