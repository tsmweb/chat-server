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
	ctx      context.Context
	executor *executor.Executor

	chUser    chan user.User
	chMessage chan message.Message
	chError   chan ErrorEvent

	userDecoder user.Decoder
	msgDecoder  message.Decoder

	consumeUser    kafka.Consumer
	consumeMessage kafka.Consumer

	handleUser    HandleUser
	handleMessage HandleMessage
	handleError   HandleError
}

func (b *Broker) run() {
	// Executor to perform background processing,
	// limiting resource consumption when executing a collection of jobs.
	b.executor = executor.New(config.GoPoolSize())

	go b.messageProcessor()
	go b.usersConsumer()
	go b.messagesConsumer()
}

func (b *Broker) stop() {
	b.executor.Shutdown()

	b.handleError.Close()
}

func (b *Broker) messageProcessor() {
loop:
	for {
		select {
		case usr := <-b.chUser:
			b.executor.Schedule(b.userTask(usr))

		case msg := <-b.chMessage:
			b.executor.Schedule(b.messageTask(msg))

		case err := <-b.chError:
			b.executor.Schedule(b.errorTask(err))

		case <-b.ctx.Done():
			break loop
		}
	}

	b.stop()
}

func (b *Broker) usersConsumer() {
	defer b.consumeUser.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			b.chError <- *NewErrorEvent("", "Broker.usersConsumer()", err.Error())
			return
		}

		var u user.User
		if err := b.userDecoder.Unmarshal(event.Value, &u); err != nil {
			b.chError <- *NewErrorEvent("", "Broker.usersConsumer()", err.Error())
			return
		}

		b.chUser <- u
	}

	b.consumeUser.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) messagesConsumer() {
	defer b.consumeMessage.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			b.chError <- *NewErrorEvent("", "Broker.messagesConsumer()", err.Error())
			return
		}

		var m message.Message
		if err := b.msgDecoder.Unmarshal(event.Value, &m); err != nil {
			b.chError <- *NewErrorEvent("", "Broker.messagesConsumer()", err.Error())
			return
		}

		b.chMessage <- m
	}

	b.consumeMessage.Subscribe(b.ctx, callbackFn)
}

func (b *Broker) userTask(usr user.User) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.handleUser.Execute(ctx, usr, b.chMessage); err != nil {
			b.chError <- *err
		}
	}
}

func (b *Broker) messageTask(msg message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := b.handleMessage.Execute(ctx, msg); err != nil {
			b.chError <- *err
		}
	}
}

func (b *Broker) errorTask(errEvent ErrorEvent) func(ctx context.Context) {
	return func(ctx context.Context) {
		log.Println(string(errEvent.ToJSON()))
		b.handleError.Execute(ctx, errEvent)
	}
}
