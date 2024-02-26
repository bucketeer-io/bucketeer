# Summary

Currently, we are calling the batch service to process the events from the PubSub every 1 minute. But it would be better
if we can process the events as soon as possible.

To reduce the latency of processing the events and let end users receive Slack notifications more quickly, we will make
improvements to the batch service event processing job.

## Implementation

Basically, we will create a new puller to process the events from the PubSub for each topic. The new architecture should
be able to add new topics easily, because we will have more topics in the future.
