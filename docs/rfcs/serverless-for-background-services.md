# Summary

Currently, we are running the background services as a background process in the compute engine. It will subscribe the
events from the PubSub and then process the events. But it will waste the compute engine and PubSub resources when there
is no events to process.

To reduce compute engine and PubSub costs, we will implement serverless architecture for background services:

- AutoOps
    - event-persister-evaluation-events-ops
    - event-persister-goal-events-ops
- Experiments
    - event-persister-evaluation-events-dwh
    - event-persister-goal-events-dwh
    - experiment-calculator

## Implementation

There are two ways to implement serverless architecture for background services:

1. Use Faas (Function as a Service) to implement the background services.
2. Implement the background services as a cron job to check the events from the PubSub.

### Faas

We have evaluated the OpenSource Faas frameworks, such
as [OpenFaaS](https://www.openfaas.com/), [Kubeless](https://kubeless.io/), [Knative](https://knative.dev/), etc.

But we found that they are not mature enough to support our requirements and bring more complexity to our system.
So we will implement the background services as a cron job.

### Cron Job

TL;DR: We will turn the background services as a cron job to check if there are experiments/auto ops. If there are
experiments/auto ops, it will process the events.

#### Process flow

1. Use the `Ticker` from the `time` package to run the cron job every 1 minute.
2. Check if there are experiments/auto ops to process.
3. If there are experiments/auto ops, then we subscribe PubSub to get the events and process them.
4. If there are no experiments/auto ops, then we will unsubscribe PubSub.

#### How we check if there are experiments/auto ops to process?

* Experiments

  Check the `experiment` table to find the experiments which `status` is `RUNNING`.
* AutoOps

  Check the `auto_ops` and `auto_ops_count` table to find the auto ops which `ops_event_count` and `evaluation_count` is
  greater than 0.

#### How we subscribe/unsubscribe PubSub?

We can create a new puller while there are experiments/auto ops to process as the following code:

```go
package main

import "go.uber.org/zap"

type Persister struct {
	client pubsub.Client
	opts   *options
	logger *zap.Logger
	ctx    context.Context
	cancel func()
}

func (p Persister) subscribe(ctx context.Context) {
	puller := client.createPuller(ctx)
	cctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	puller.Run(cctx, func() {
		// process events
	})
}

func (p Persister) unsubscribe() {
	p.cancel()
	p.cancel = nil
}

func (p Persister) IsRunning() bool {
	return p.cancel != nil
}

```

Explanation:

- When `Ticker` triggers the cron job, we call `IsRunning` to check if there are already a puller running, as we use
  `Ticker` to run the cron job every 1 minute, so we don't need to use `Mutex` to protect the `cancel` field.
- If there are no puller running, then we check if there are experiments/auto ops to process, if there are
  experiments/auto ops to process, then we call `subscribe` to subscribe PubSub.
- If there are puller running, then we check if there are experiments/auto ops to process, if there are no
  experiments/auto ops to process, then we call `unsubscribe` to unsubscribe PubSub, otherwise, if there are
  experiments/auto ops to process, then we do nothing.