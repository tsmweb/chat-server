package ebus

import (
	"sync"
)

// DataEvent represents an event posted to a topic.
type DataEvent struct {
	Data  interface{}
	Topic string
}

// DataChannel is a channel which can accept a DataEvent.
type DataChannel chan DataEvent

// DataChannelSet is a set of DataChannels.
type DataChannelSet map[DataChannel]bool

// Subscription represents a subscription to a topic.
type Subscription struct {
	Event       <-chan DataEvent
	Unsubscribe func()
}

// EventBus stores the information about subscribers interested for a particular topic.
type EventBus struct {
	subscribers map[string]DataChannelSet
	mu          sync.RWMutex
}

// NewEventBus creates an EventBus instance.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: map[string]DataChannelSet{},
	}
}

var eventBus *EventBus

// Instance creates an single EventBus instance.
func Instance() *EventBus {
	if eventBus == nil {
		eventBus = NewEventBus()
	}
	return eventBus
}

// Publish publishes data in the topic provided.
func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if chans, found := eb.subscribers[topic]; found {
		// go func(data DataEvent, dataChannelSet DataChannelSet) {
		// 	for ch := range dataChannelSet {
		// 		ch <- data
		// 	}
		// }(DataEvent{Data: data, Topic: topic}, chans)

		for ch := range chans {
			ch <- DataEvent{Data: data, Topic: topic}
		}
	}
}

// Subscribe to the topic provided to receive data events.
func (eb *EventBus) Subscribe(topic string) *Subscription {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan DataEvent, 1)

	if prev, found := eb.subscribers[topic]; found {
		prev[ch] = true
	} else {
		eb.subscribers[topic] = DataChannelSet{ch: true}
	}

	return &Subscription{
		Event: ch,
		Unsubscribe: func() {
			eb.unsubscribe(topic, ch)
		},
	}
}

func (eb *EventBus) unsubscribe(topic string, ch chan DataEvent) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if chans, found := eb.subscribers[topic]; found {
		delete(chans, ch)
	}
}