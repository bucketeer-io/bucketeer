# Evaluation for updated FeatureFlag


## Proposal

- This feature evaluates only those feature flags that have been updated since the last evaluation.
  - Computational efficiency will be improved because the number of the feature flags to be evaluated will be reduced.
- We can make the tag optional.
  - Previously, specifying a tag was required when executing GetEvaluations to reduce the response size.
  - But, with this proposal, the response size can be kept small, so specifying a tag can be made optional.




## Implementation

The SDK stores the timestamp of the evaluation's execution.
When GetEvaluations is called again, the timestamp is sent to the server as a request parameter.

The server checks it against the updatedAt value of the feature flags, and only evaluates those that have been updated since the previous evaluation.

As an exception, the following feature flags must be evaluated regardless of the timestamp value.
- Feature flags that depend on the feature flags that need to be evaluated


In addition, the following changes are required in both the server and SDK implementations:
- Since only the updated feature flags are evaluated and returned, the SDK's implementation for updating local data needs to be changed.
- The server needs to return information about archived feature flags to the SDK.
  - So the server must put the archived feature flags to Redis.
- The validation of the request parameter `Tag` needs to be removed.


### Evaluation on Server
The following diagram shows the dependency relationship between multiple feature flags:

![](http://g.gravizo.com/g?
  digraph {
    rankdir=LR;
    node [shape = circle];
    D
    C -> L -> M
    L -> N
    B
    A -> E -> G -> H -> I -> K
    H -> J
    A -> F
  }
)


#### Pattern1
Assuming only featureA, featureB, featureC, and featureD were updated after the last evaluation.

Only featuresA, featureB, featureC, and featureD are evaluated, since no features depend on these features.

| updated                                | evaluated                               |
|----------------------------------------|-----------------------------------------|
| featureA, featureB, featureC, featureD | featureA, featureB, featureC, featureD  |


#### Pattern2
Assuming only featureF was updated after the last evaluation.

Since featureA specifies featureF as Prerequisite, the evaluation of featureA may also change.

Therefore, featureA also needs to be re-evaluated.

| updated  | evaluated          |
|----------|--------------------|
| featureF | featureA, featureF |

#### Pattern3
Assuming only featureE and featureK were updated after the last evaluation.

featureK is a part of the chain of dependencies, where the featureA depends on featureE, and featureE depends on featureG...

In this case, we must evaluate all feature flags involved in dependencies.

| updated            | evaluated                                                  |
|--------------------|------------------------------------------------------------|
| featureE, featureK | featureA, featureE, featureG, featureH, featureI, featureK |


### SDK

#### Timestamp of the evaluation's execution
SDK stores the timestamp which is contained in the response of GetEvaluations in local storage.
The saved timestamp will be included in the next GetEvaluations request.

#### Storing the evaluation results in local storage
In the current implementation, the server returns all feature flags evaluation results, and the SDK inserts those in local storage after deleting the local data.
However, this proposal would change that behavior so that only some feature flags evaluation results are returned, not all.
It will also return the evaluation results of archived feature flags.
We must modify the implementation of SDK to address these changes.

#### UserEvaluationsID

In the current implementation, we use the `UserEvaluationsID` to detect the feature flags updates.
But we will be able to detect it via the timestamp, so the `UserEvaluationsID` is no longer necessary.

## Release Steps

During implementation changes, we ensure that any version of the SDK will work properly.

1. Add a function to get the feature flag dependencies.
2. Change the server behavior to put archived feature flags to Redis.(The server does not return archived flags to the SDK at this time.)
3. Add a `EvaluatedAt` field to GetEvaluationsRequest object.
4. Add implementation to check the timestamp `EvaluatedAt` against feature flag's `UpdatedAt` field.(`UserEvaluationsID` is also continue to be accepted.)
5. Modify each SDK to support `EvaluatedAt` and differential update the local data.
6. Make GetEvaluationsRequest's `UserEvaluationsID` field deprecated and `Tag` field optional.
