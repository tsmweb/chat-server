package kafka

import (
	"context"
	skafka "github.com/segmentio/kafka-go"
	"github.com/tsmweb/chat-service/chat"
	"log"
	"time"
)

// kafka implementation for chat.Kafka interface with segmentio/kafka-go library.
type kafka struct {
	kafkaBrokerUrls []string
	clientId        string
	debug           bool
}

// NewKafka creates a new instance of chat.Kafka.
func New(kafkaBrokerUrls []string, clientId string) chat.Kafka {
	return &kafka{
		kafkaBrokerUrls: kafkaBrokerUrls,
		clientId:        clientId,
		debug:           false,
	}
}

// Debug enables logging of incoming events.
func (k *kafka) Debug(debug bool) {
	k.debug = debug
}

func (k *kafka) dialer() *skafka.Dialer {
	return &skafka.Dialer{
		ClientID:  k.clientId,
		DualStack: true,
		Timeout:   10 * time.Second,
	}
}

// NewProducer creates a new KafkaProducer to produce events on a topic.
func (k *kafka) NewProducer(topic string) chat.KafkaProducer {
	w := &skafka.Writer{
		Addr:         skafka.TCP(k.kafkaBrokerUrls...),
		Topic:        topic,
		RequiredAcks: skafka.RequireAll,
		BatchTimeout: time.Millisecond,
		Compression:  skafka.Snappy,
	}
	return newProducer(w)
}

// NewConsumer creates a new KafkaConsumer to consume events from a topic.
func (k *kafka) NewConsumer(groupID, topic string) chat.KafkaConsumer {
	config := skafka.ReaderConfig{
		Brokers:         k.kafkaBrokerUrls,
		GroupID:         groupID,
		Topic:           topic,
		Dialer:          k.dialer(),
		MinBytes:        10e3,        // 10KB
		MaxBytes:        10e6,        // 10MB
		MaxWait:         time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
		//CommitInterval: time.Second, // flushes commits to KafkaWrap every second
	}
	r := skafka.NewReader(config)
	return newConsumer(r, k.debug)
}

// Producer provide methods for producing events for a given topic.
type Producer struct {
	writer *skafka.Writer
}

func newProducer(w *skafka.Writer) chat.KafkaProducer {
	return &Producer{writer: w}
}

// Publish produces and sends an event for a kafka topic.
// The context passed as first argument may also be used to asynchronously
// cancel the operation.
func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	message := skafka.Message{
		Key: key,
		Value: value,
	}
	return p.writer.WriteMessages(ctx, message)
}

// Close flushes pending writes, and waits for all writes to complete before
// returning. Calling Close also prevents new writes from being submitted to
// the Writer, further calls to WriteMessages and the like will fail with
// io.ErrClosedPipe.
func (p *Producer) Close() {
	p.writer.Close()
}

// Consumer provide methods for consuming events on a given topic.
type Consumer struct {
	reader *skafka.Reader
	debug bool
}

func newConsumer(r *skafka.Reader, d bool) chat.KafkaConsumer {
	return &Consumer{
		reader: r,
		debug: d,
	}
}

// Subscribe consumes the events of a topic and passes the event to the
// informed callback function. The method call blocks until an error occurs.
// The program may also specify a context to asynchronously cancel the blocking operation.
func (c *Consumer) Subscribe(ctx context.Context, callbackFn func(event *chat.KafkaEvent, err error)) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			callbackFn(nil, err)
			break
		}

		if c.debug {
			c.logMessage(m)
		}

		callbackFn(c.makeEvent(m), nil)
	}
}

// Close closes the stream, preventing the program from reading any more
// events from it.
func (c *Consumer) Close() {
	c.reader.Close()
}

func (c *Consumer) makeEvent(m skafka.Message) *chat.KafkaEvent {
	return &chat.KafkaEvent{
		Topic:  m.Topic,
		Key:    m.Key,
		Value:  m.Value,
		Time:   m.Time,
	}
}

func (c *Consumer) logMessage(m skafka.Message) {
	log.Printf("[>] ReadMessage - TIME: %s | TOPIC: %s | PARTITION: %d | OFFSET: %d | HEADER: %s | SIZE PAYLOAD: %d\n",
		m.Time, m.Topic, m.Partition, m.Offset, m.Headers, len(m.Value))
}