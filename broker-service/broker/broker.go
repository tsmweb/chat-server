package broker

import (
	"context"
	"log"
	"sync"

	"github.com/tsmweb/broker-service/broker/group"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/common/service"
	"github.com/tsmweb/broker-service/config"
	"github.com/tsmweb/go-helper-api/concurrent/gopool"
	"github.com/tsmweb/go-helper-api/kafka"
)

type Broker struct {
	ctx                    context.Context
	chUser                 chan user.User
	chUserPresence         chan user.User
	chMessage              chan message.Message
	chOfflineMessage       chan message.Message
	chUserMessage          chan message.Message
	chGroupEvent           chan group.Event
	chUserEvent            chan user.Event
	userDecoder            user.Decoder
	msgDecoder             message.Decoder
	groupEventDecoder      group.EventDecoder
	userEventDecoder       user.EventDecoder
	userConsumer           kafka.Consumer
	userPresenceConsumer   kafka.Consumer
	messageConsumer        kafka.Consumer
	offlineMessageConsumer kafka.Consumer
	groupEventConsumer     kafka.Consumer
	userEventConsumer      kafka.Consumer
	userHandler            UserHandler
	userPresenceHandler    UserPresenceHandler
	messageHandler         MessageHandler
	offlineMessageHandler  OfflineMessageHandler
	groupEventHandler      GroupEventHandler
	userEventHandler       UserEventHandler
}

// NewBroker creates an instance of Broker.
func NewBroker(
	ctx context.Context,
	userDecoder user.Decoder,
	msgDecoder message.Decoder,
	groupEventDecoder group.EventDecoder,
	userEventDecoder user.EventDecoder,
	userConsumer kafka.Consumer,
	userPresenceConsumer kafka.Consumer,
	messageConsumer kafka.Consumer,
	offlineMessageConsumer kafka.Consumer,
	groupEventConsumer kafka.Consumer,
	userEventConsumer kafka.Consumer,
	userHandler UserHandler,
	userPresenceHandler UserPresenceHandler,
	messageHandler MessageHandler,
	offlineMessageHandler OfflineMessageHandler,
	groupEventHandler GroupEventHandler,
	userEventHandler UserEventHandler,
) *Broker {
	broker := &Broker{
		ctx:                    ctx,
		chUser:                 make(chan user.User),
		chUserPresence:         make(chan user.User),
		chMessage:              make(chan message.Message),
		chOfflineMessage:       make(chan message.Message),
		chUserMessage:          make(chan message.Message),
		chGroupEvent:           make(chan group.Event),
		chUserEvent:            make(chan user.Event),
		userDecoder:            userDecoder,
		msgDecoder:             msgDecoder,
		groupEventDecoder:      groupEventDecoder,
		userEventDecoder:       userEventDecoder,
		userConsumer:           userConsumer,
		userPresenceConsumer:   userPresenceConsumer,
		messageConsumer:        messageConsumer,
		offlineMessageConsumer: offlineMessageConsumer,
		groupEventConsumer:     groupEventConsumer,
		userEventConsumer:      userEventConsumer,
		userHandler:            userHandler,
		userPresenceHandler:    userPresenceHandler,
		messageHandler:         messageHandler,
		offlineMessageHandler:  offlineMessageHandler,
		groupEventHandler:      groupEventHandler,
		userEventHandler:       userEventHandler,
	}

	return broker
}

func (b *Broker) Start() {
	go b.usersConsumer()
	go b.usersPresenceConsumer()
	go b.messagesConsumer()
	go b.offlineMessagesConsumer()
	go b.groupEventsConsumer()
	go b.userEventsConsumer()

	b.messageProcessor()
}

func (b *Broker) messageProcessor() {
	// gopool.Pool to perform background processing,
	// limiting resource consumption when executing a collection of tasks.
	workerSize := config.GoPoolSize()
	queueSize := 1

	poolUsers := gopool.New(workerSize, queueSize)
	defer poolUsers.Close()
	poolMessages := gopool.New(workerSize, queueSize)
	defer poolMessages.Close()
	poolEvents := gopool.New(workerSize, queueSize)
	defer poolEvents.Close()

	var wg sync.WaitGroup

	// User
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			log.Println("[STOP] Broker::chUser")
		}()

		var userWG sync.WaitGroup

		for u := range b.chUser {
			userWG.Add(1)
			if err := poolUsers.Schedule(b.userTask(u, &userWG)); err != nil {
				log.Printf("[ERROR] Broker::poolUsers: %v\n", err)
			}
		}

		userWG.Wait()
		close(b.chUserMessage)
	}()

	// User Presence
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			log.Println("[STOP] Broker::chUserPresence")
		}()

		for u := range b.chUserPresence {
			if err := poolUsers.Schedule(b.userPresenceTask(u)); err != nil {
				log.Printf("[ERROR] Broker::poolUsers: %v\n", err)
			}
		}
	}()

	// User Message
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			log.Println("[STOP] Broker::chUserMessage")
		}()

		for m := range b.chUserMessage {
			if err := poolMessages.Schedule(b.messageTask(m)); err != nil {
				log.Printf("[ERROR] Broker::poolMessages: %v\n", err)
			}
		}
	}()

	// Message
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			log.Println("[STOP] Broker::chMessage")
		}()

		for m := range b.chMessage {
			if err := poolMessages.Schedule(b.messageTask(m)); err != nil {
				log.Printf("[ERROR] Broker::poolMessages: %v\n", err)
			}
		}
	}()

	// Offline Message
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			log.Println("[STOP] Broker::chOfflineMessage")
		}()

		for m := range b.chOfflineMessage {
			if err := poolMessages.Schedule(b.offlineMessageTask(m)); err != nil {
				log.Printf("[ERROR] Broker::poolMessages: %v\n", err)
			}
		}
	}()

	// Group Events
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			log.Println("[STOP] Broker::chGroupEvent")
		}()

		for e := range b.chGroupEvent {
			if err := poolEvents.Schedule(b.groupEventTask(e)); err != nil {
				log.Printf("[ERROR] Broker::poolEvents: %v\n", err)
			}
		}
	}()

	// User Events
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			log.Println("[STOP] Broker::chUserEvent")
		}()

		for e := range b.chUserEvent {
			if err := poolEvents.Schedule(b.userEventTask(e)); err != nil {
				log.Printf("[ERROR] Broker::poolEvents: %v\n", err)
			}
		}
	}()

	wg.Wait()
}

func (b *Broker) usersConsumer() {
	defer func() {
		b.userConsumer.Close()
		close(b.chUser)
		log.Println("[STOP] Broker::usersConsumer")
	}()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			service.Error("", "Broker::usersConsumer", err)
			return
		}

		var usr user.User
		if err = b.userDecoder.Unmarshal(event.Value, &usr); err != nil {
			service.Error(string(event.Key), "Broker::usersConsumer::userDecoder", err)
			return
		}

		b.chUser <- usr
	}

	b.userConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) usersPresenceConsumer() {
	defer func() {
		b.userPresenceConsumer.Close()
		close(b.chUserPresence)
		log.Println("[STOP] Broker::usersPresenceConsumer")
	}()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			service.Error("", "Broker::usersPresenceConsumer", err)
			return
		}

		var usr user.User
		if err = b.userDecoder.Unmarshal(event.Value, &usr); err != nil {
			service.Error(string(event.Key), "Broker::usersPresenceConsumer::userDecoder", err)
			return
		}

		b.chUserPresence <- usr
	}

	b.userPresenceConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) messagesConsumer() {
	defer func() {
		b.messageConsumer.Close()
		close(b.chMessage)
		log.Println("[STOP] Broker::messagesConsumer")
	}()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			service.Error("", "Broker::messagesConsumer", err)
			return
		}

		var msg message.Message
		if err = b.msgDecoder.Unmarshal(event.Value, &msg); err != nil {
			service.Error(string(event.Key), "Broker::messagesConsumer::msgDecoder", err)
			return
		}

		b.chMessage <- msg
	}

	b.messageConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) offlineMessagesConsumer() {
	defer func() {
		b.offlineMessageConsumer.Close()
		close(b.chOfflineMessage)
		log.Println("[STOP] Broker::offlineMessagesConsumer")
	}()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			service.Error("", "Broker::offlineMessagesConsumer", err)
			return
		}

		var msg message.Message
		if err = b.msgDecoder.Unmarshal(event.Value, &msg); err != nil {
			service.Error(string(event.Key), "Broker::offlineMessagesConsumer::msgDecoder", err)
			return
		}

		b.chOfflineMessage <- msg
	}

	b.offlineMessageConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) groupEventsConsumer() {
	defer func() {
		b.groupEventConsumer.Close()
		close(b.chGroupEvent)
		log.Println("[STOP] Broker::groupEventsConsumer")
	}()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			service.Error("", "Broker::groupEventsConsumer", err)
			return
		}

		var groupEvent group.Event
		if err = b.groupEventDecoder.Unmarshal(event.Value, &groupEvent); err != nil {
			service.Error(string(event.Key), "Broker::groupEventsConsumer::groupEventDecoder", err)
			return
		}

		b.chGroupEvent <- groupEvent
	}

	b.groupEventConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) userEventsConsumer() {
	defer func() {
		b.userEventConsumer.Close()
		close(b.chUserEvent)
		log.Println("[STOP] Broker::userEventsConsumer")
	}()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			service.Error("", "Broker::userEventsConsumer", err)
			return
		}

		var userEvent user.Event
		if err = b.userEventDecoder.Unmarshal(event.Value, &userEvent); err != nil {
			service.Error(string(event.Key), "Broker::userEventsConsumer::userEventDecoder", err)
			return
		}

		b.chUserEvent <- userEvent
	}

	b.userEventConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) userTask(usr user.User, wg *sync.WaitGroup) func(ctx context.Context) {
	return func(ctx context.Context) {
		defer wg.Done()

		if err := b.userHandler.Execute(ctx, usr, b.chUserMessage); err != nil {
			service.Error(usr.ID, "Broker::userTask::userHandler", err)
		}
	}
}

func (b *Broker) userPresenceTask(usr user.User) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.userPresenceHandler.Execute(ctx, usr); err != nil {
			service.Error(usr.ID, "Broker::userPresenceTask::userPresenceHandler", err)
		}
	}
}

func (b *Broker) messageTask(msg message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.messageHandler.Execute(ctx, msg); err != nil {
			service.Error(msg.ID, "Broker::messageTask::messageHandler", err)
		}
	}
}

func (b *Broker) offlineMessageTask(msg message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.offlineMessageHandler.Execute(ctx, msg); err != nil {
			service.Error(msg.ID, "Broker::offlineMessageTask::offlineMessageHandler", err)
		}
	}
}

func (b *Broker) groupEventTask(evt group.Event) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.groupEventHandler.Execute(ctx, evt); err != nil {
			service.Error(evt.GroupID, "Broker::groupEventTask::groupEventHandler", err)
		}
	}
}

func (b *Broker) userEventTask(evt user.Event) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.userEventHandler.Execute(ctx, evt); err != nil {
			service.Error(evt.UserID, "Broker::userEventTask::userEventHandler", err)
		}
	}
}
