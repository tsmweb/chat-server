package broker

import (
	"context"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/config"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"github.com/tsmweb/go-helper-api/kafka"
	"log"
)

type Broker struct {
	ctx                    context.Context
	executor               *executor.Executor
	chUser                 chan user.User
	chUserPresence         chan user.User
	chMessage              chan message.Message
	chOfflineMessage       chan message.Message
	chError                chan ErrorEvent
	userDecoder            user.Decoder
	msgDecoder             message.Decoder
	userConsumer           kafka.Consumer
	userPresenceConsumer   kafka.Consumer
	messageConsumer        kafka.Consumer
	offlineMessageConsumer kafka.Consumer
	userHandler            UserHandler
	userPresenceHandler    UserPresenceHandler
	messageHandler         MessageHandler
	offlineMessageHandler  OfflineMessageHandler
	errorHandler           ErrorHandler
}

// NewBroker creates an instance of Broker.
func NewBroker(
	ctx context.Context,
	userDecoder user.Decoder,
	msgDecoder message.Decoder,
	userConsumer kafka.Consumer,
	userPresenceConsumer kafka.Consumer,
	messageConsumer kafka.Consumer,
	offlineMessageConsumer kafka.Consumer,
	userHandler UserHandler,
	userPresenceHandler UserPresenceHandler,
	messageHandler MessageHandler,
	offlineMessageHandler OfflineMessageHandler,
	errorHandler ErrorHandler,
) *Broker {
	broker := &Broker{
		ctx:                    ctx,
		chUser:                 make(chan user.User),
		chUserPresence:         make(chan user.User),
		chMessage:              make(chan message.Message),
		chOfflineMessage:       make(chan message.Message),
		chError:                make(chan ErrorEvent),
		userDecoder:            userDecoder,
		msgDecoder:             msgDecoder,
		userConsumer:           userConsumer,
		userPresenceConsumer:   userPresenceConsumer,
		messageConsumer:        messageConsumer,
		offlineMessageConsumer: offlineMessageConsumer,
		userHandler:            userHandler,
		userPresenceHandler:    userPresenceHandler,
		messageHandler:         messageHandler,
		offlineMessageHandler:  offlineMessageHandler,
		errorHandler:           errorHandler,
	}

	return broker
}

func (b *Broker) Start() {
	// Executor to perform background processing,
	// limiting resource consumption when executing a collection of jobs.
	b.executor = executor.New(config.GoPoolSize())

	go b.usersConsumer()
	go b.usersPresenceConsumer()
	go b.messagesConsumer()
	go b.offlineMessagesConsumer()
	b.messageProcessor()
}

func (b *Broker) stop() {
	b.executor.Shutdown()

	b.errorHandler.Close()
}

func (b *Broker) messageProcessor() {
loop:
	for {
		select {
		case usr := <-b.chUser:
			b.executor.Schedule(b.userTask(usr))

		case userPresence := <-b.chUserPresence:
			b.executor.Schedule(b.userPresenceTask(userPresence))

		case msg := <-b.chMessage:
			b.executor.Schedule(b.messageTask(msg))

		case msg := <-b.chOfflineMessage:
			b.executor.Schedule(b.offlineMessageTask(msg))

		case err := <-b.chError:
			b.executor.Schedule(b.errorTask(err))

		case <-b.ctx.Done():
			break loop
		}
	}

	b.stop()
}

func (b *Broker) usersConsumer() {
	defer b.userConsumer.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			b.chError <- *NewErrorEvent("", "Broker.usersConsumer()", err.Error())
			return
		}

		var usr user.User
		if err := b.userDecoder.Unmarshal(event.Value, &usr); err != nil {
			b.chError <- *NewErrorEvent("", "Broker.usersConsumer()", err.Error())
			return
		}

		b.chUser <- usr
	}

	b.userConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) usersPresenceConsumer() {
	defer b.userPresenceConsumer.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			b.chError <- *NewErrorEvent("", "Broker.usersPresenceConsumer()", err.Error())
			return
		}

		var usr user.User
		if err := b.userDecoder.Unmarshal(event.Value, &usr); err != nil {
			b.chError <- *NewErrorEvent("", "Broker.usersPresenceConsumer()", err.Error())
			return
		}

		b.chUserPresence <- usr
	}

	b.userPresenceConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) messagesConsumer() {
	defer b.messageConsumer.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			b.chError <- *NewErrorEvent("", "Broker.messagesConsumer()", err.Error())
			return
		}

		var msg message.Message
		if err := b.msgDecoder.Unmarshal(event.Value, &msg); err != nil {
			b.chError <- *NewErrorEvent("", "Broker.messagesConsumer()", err.Error())
			return
		}

		b.chMessage <- msg
	}

	b.messageConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) offlineMessagesConsumer() {
	defer b.offlineMessageConsumer.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			b.chError <- *NewErrorEvent("", "Broker.offlineMessagesConsumer()", err.Error())
			return
		}

		var msg message.Message
		if err := b.msgDecoder.Unmarshal(event.Value, &msg); err != nil {
			b.chError <- *NewErrorEvent("", "Broker.offlineMessagesConsumer()", err.Error())
			return
		}

		b.chOfflineMessage <- msg
	}

	b.offlineMessageConsumer.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) userTask(usr user.User) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.userHandler.Execute(ctx, usr, b.chMessage); err != nil {
			b.chError <- *err
		}
	}
}

func (b *Broker) userPresenceTask(usr user.User) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.userPresenceHandler.Execute(ctx, usr); err != nil {
			b.chError <- *err
		}
	}
}

func (b *Broker) messageTask(msg message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.messageHandler.Execute(ctx, msg); err != nil {
			b.chError <- *err
		}
	}
}

func (b *Broker) offlineMessageTask(msg message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.offlineMessageHandler.Execute(ctx, msg); err != nil {
			b.chError <- *err
		}
	}
}

func (b *Broker) errorTask(errEvent ErrorEvent) func(ctx context.Context) {
	return func(ctx context.Context) {
		log.Println(string(errEvent.ToJSON()))
		b.errorHandler.Execute(ctx, errEvent)
	}
}
