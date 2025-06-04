# Summary

Currently, when users configure rules on the Targeting tab in the console, they need to manually type custom attribute keys. This manual process can lead to typos and misconfigurations, potentially causing incorrect conditions when evaluating end-users.

Since these attributes are sent from the SDK to the server, we can automate this process by generating a list to display on the console. This will improve user experience and reduce configuration errors.

# Background

The current manual attribute key input process has several issues:
- Prone to typos and human error
- Time-consuming for users to type attribute keys
- No validation of attribute key existence
- Potential for misconfiguration leading to incorrect feature flag evaluations

# Solution

We will implement an automated system to discover and manage user attributes:

1. Extract user attribute information from SDK requests
2. Compare with cached data in Redis
3. Publish only new attribute information using PubSub
4. Saving attributes to a database via PubSub
5. Provide an API for the console to retrieve the attribute list

## System Architecture

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

# Implementation Details

## Cache

- Create `UserAttributesCache` in the cache package
  - Key: environment_id
  - Value: []string (user_attribute_keys)

## PubSub

- Create new topic: `user-attribute-event`
- Create new subscription: `user-attribute-event-persister`
- Add topic and subscription definitions to YAML configuration
- Implement `UserAttributePersister` in the Processor

## API Server

- Implement attribute discovery in `getEvaluations` and `getEvaluation` requests
- Compare with UserAttributeCache
- Publish new attributes to `user-attribute-event` topic

## API

Add new API to the Environment Package:

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

1. I intentionally do not implement attribute deletion functionality
   - The deleted attribute may be needed again in the future, but there is currently no way to undo the deletion using the console.
   - However, the console takes into account the large number of user attributes by providing incremental search to improve usability.


# Testing

- The e2e test is performed in the following steps:
- Test flow:
  1. Send request via `GetEvaluation`
  2. Wait for processing
  3. Verify attributes via `ListUserAttributes` API

# Release Steps

1. Cache Implementation
   - Implement UserAttributesCache
   - Implement UserAttributesCacher

2.  PubSub Implementation
   - Implement Processor (without Subscriber connection)

3. API Server Implementation and Set up topic and subscription
   - Add publishing logic
   - Configure PubSub in Dev environment
   - Connect UserAttributePersister
   - I run e2e tests and, if there are no problems, we release them to the production environment.

4. API Implementation
   - Implement ListUserAttributes endpoint
   - Add API documentation
