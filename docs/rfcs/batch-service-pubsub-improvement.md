# Summary

Currently, we are calling the batch service to process the events from the PubSub every 1 minute. But it would be better
if we can process the events as soon as possible.

To reduce the latency of processing the events and let end users receive Slack notifications more quickly, we will make
improvements to the batch service event processing job.

## Implementation

Basically, we will create a new puller to process the events from the PubSub for each topic. The new architecture should
be able to add new topics easily, because we will have more topics in the future.

### Architecture

The multi pub-sub architecture will have the following components:

- **Configuration**: The configuration for the subscriber, which includes the topic and the subscription name.
- **Subscriber**: The subscriber that will pull the messages from the PubSub.
- **Processor**: The processor that will process the messages from the PubSub.

We will create a multi pub-sub Golang struct to handle the configuration, subscriber, and processor. The multi pub-sub
will have the following methods:

- **NewMultiPubSub**: The constructor for the multi pub-sub.
- **Start**: The method to start the multi pub-sub.
- **Stop**: The method to stop the multi pub-sub.
- **AddSubscriber**: The method to add a new subscriber to the multi pub-sub.

#### Configuration

The configuration will be used to configure `topic`, `subscription`, and other settings for the subscriber.

```go
package subscriber

type Configuration struct {
	project                      string
	subscription                 string
	topic                        string
	pullerNumGoroutines          int
	pullerMaxOutstandingMessages int
	pullerMaxOutstandingBytes    int
}

```

#### Processor

The processor will process the messages from the PubSub. The processor will be a function that will be called when the
subscriber receives a message from the PubSub.

```go
package subscriber

type Processor func(msg *puller.Message)

```

#### Subscriber

We will create a new subscriber for each configuration in multi pub-sub Golang struct. When the multi pub-sub starts, it
will start the subscriber for each configuration.

```go
package subscriber

type Subscriber struct {
	name          string
	configuration Configuration
	processor     Processor
}

```

The Subscriber struct have a `Run` function to start the subscriber. The `Run` function will pull the messages from the
PubSub and call the processor function to process the messages.

```go
package subscriber

func (s Subscriber) Run(ctx context.Context) {
	puller := createPuller(ctx, s.configuration)
	group := errgroup.Group{}
	group.Go(func() error {
		puller.Pull(ctx, func(ctx context.Context, msg *puller.Message) {
			s.processor(msg)
		})
		return nil
	})
	group.Go(func() error {
		for {
			select {
			case msg, ok := <-puller.MessageCh():
				if !ok {
					logger.error("receive msg error")
					return nil
				}
				s.processor(msg)
			case <-ctx.Done():
				return nil
			}
		}
	})
}

```

### Multi Pub-Sub

The multi pub-sub struct will combine the configuration, subscriber, and processor. The multi pub-sub will look like
this:

```go
package subscriber

import "go.uber.org/zap"

type MultiPubSub struct {
	subscribers []Subscriber
	logger      zap.Logger
}

```

We will create a new multi pub-sub for each topic. When the multi pub-sub starts, it will start the subscriber for each
topic, like this:

```go
package subscriber

func NewMultiPubSub(logger zap.logger) *MultiPubSub {
	return *MultiPubSub{
		subscribers: make([]Subscriber),
		logger:      logger,
	}
}

func (m MultiPubSub) AddSubscriber(subscriber Subscriber) {
	m.subscribers = append(m.subscribers, subscriber)
}

func (m MultiPubSub) Start(multiPubSub *MultiPubSub) {
	for _, subscriber := range multiPubSub.subscribers {
		go subscriber.Start()
	}
}
```
