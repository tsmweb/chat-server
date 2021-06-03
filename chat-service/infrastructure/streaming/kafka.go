package streaming

import (
	"context"
	"github.com/tsmweb/chat-service/common/ebus"
	"github.com/tsmweb/go-helper-api/integration"
)

// kafka implementation for core.EBus interface.
type kafka struct {
	kafka *integration.Kafka
}

// NewKafka creates a new instance of core.EBus.
func NewKafka(kafkaBrokerUrls []string, clientId string) ebus.EBus {
	return &kafka{
		kafka: integration.NewKafka(kafkaBrokerUrls, clientId),
	}
}

// Dispatch produces and sends an event for a kafka topic.
func (k *kafka) Dispatch(ctx context.Context, topic string, key string, value []byte) error {
	w := k.kafka.NewWriter(topic)
	defer w.Close()

	return w.SendEvent(ctx, []byte(key), value)
}

// Subscribe consumes the events of a kafka topic and passes the event to the
// informed callback function.
func (k *kafka) Subscribe(ctx context.Context, groupID, topic string, callbackFn func(event *ebus.Event, err error)) {
	r := k.kafka.NewReader(groupID, topic)
	defer r.Close()

	r.SubscribeTopic(ctx, func(event *integration.Event, err error) {
		if err != nil {
			callbackFn(nil, err)
		} else {
			e := &ebus.Event{
				Key:   string(event.Key),
				Value: event.Value,
				Time:  event.Time,
			}
			callbackFn(e, nil)
		}
	})
}
