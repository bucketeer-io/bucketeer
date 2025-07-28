package processor

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	domainevent "github.com/bucketeer-io/bucketeer/proto/event/domain"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type demoOrganizationCreationNotifier struct {
	sender sender.Sender
	logger *zap.Logger
}

func NewDemoOrganizationCreationNotifier(
	sender sender.Sender,
	logger *zap.Logger,
) subscriber.PubSubProcessor {
	return &demoOrganizationCreationNotifier{
		sender: sender,
		logger: logger,
	}
}

func (d demoOrganizationCreationNotifier) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				d.logger.Error("domainEventInformer: message channel closed")
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberDemoOrganizationEvent).Inc()
			d.handleMessage(msg)
		case <-ctx.Done():
			d.logger.Debug("subscriber context done, stopped processing messages")
			return nil
		}
	}
}

func (d demoOrganizationCreationNotifier) handleMessage(msg *puller.Message) {
	if id := msg.Attributes["id"]; id == "" {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberDemoOrganizationEvent, codes.BadMessage.String()).Inc()
		return
	}
	domainEvent, err := d.unmarshalMessage(msg)
	if err != nil {
		subscriberHandledCounter.WithLabelValues(subscriberDemoOrganizationEvent, codes.BadMessage.String()).Inc()
		msg.Ack()
		return
	}
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("Processing domain event: %s\n", domainEvent.Type)
}

func (d demoOrganizationCreationNotifier) unmarshalMessage(msg *puller.Message) (*domainevent.Event, error) {
	event := &domaineventproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		d.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
		return nil, err
	}
	return event, nil
}
