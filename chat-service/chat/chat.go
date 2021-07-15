package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tsmweb/chat-service/chat/message"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/pkg/concurrent"
	"github.com/tsmweb/chat-service/pkg/ebus"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/pkg/topic"
	"github.com/tsmweb/chat-service/util/connutil"
	"net"
)

type Chat struct {
	poller   epoll.EPoll
	executor concurrent.ExecutorService

	chUserIN         chan *User
	chUserOUT        chan *User
	chPublishMessage chan *message.Message
	chSendMessage    chan *message.Message
	chError          chan *ErrorEvent

	reader connutil.Reader
	writer connutil.Writer

	repository Repository
	kafka      Kafka
}

func NewChat(
	p epoll.EPoll,
	e concurrent.ExecutorService,
	reader connutil.Reader,
	writer connutil.Writer,
	r Repository,
	k Kafka,
) *Chat {
	chat := &Chat{
		poller:           p,
		executor:         e,
		chUserIN:         make(chan *User),
		chUserOUT:        make(chan *User),
		chPublishMessage: make(chan *message.Message),
		chSendMessage:    make(chan *message.Message),
		chError:          make(chan *ErrorEvent),
		reader:           reader,
		writer:           writer,
		repository:       r,
		kafka:            k,
	}

	go chat.start()
	go chat.consumerMessage()

	return chat
}

func (c *Chat) Register(userID string, conn net.Conn) error {
	user := &User{
		id:     userID,
		conn:   conn,
		reader: c.reader,
		writer: c.writer,
	}

	observer, err := c.poller.ObservableRead(conn)
	if err != nil {
		c.chError <- NewErrorEvent(user.id, "Chat.Register()", err.Error())
		return err
	}

	err = observer.Start(func(closed bool, err error) {
		if closed || err != nil {
			observer.Stop()
			c.chUserOUT <- user
			if err != nil {
				c.chError <- NewErrorEvent(user.id, "Chat.Register()", err.Error())
			}
			return
		}

		c.executor.Schedule(func(ctx context.Context) {
			sendACK := func(msgID, content string) {
				if err := user.WriteACK(msgID, content); err != nil {
					c.chError <- NewErrorEvent(user.id, "User.WriteACK()", err.Error())
				}
			}

			msg, err := user.Receive()
			if err != nil {
				observer.Stop()
				c.chUserOUT <- user
				c.chError <- NewErrorEvent(user.id, "User.Receive()", err.Error())
				return
			}
			if msg != nil {
				blocked, err := c.repository.IsBlockedUser(msg.To, msg.From)
				if err != nil {
					c.chError <- NewErrorEvent(user.id, "Repository.IsBlockedUser()", err.Error())
					return
				}
				if blocked {
					sendACK(msg.ID, fmt.Sprintf(message.BlockedMessage, msg.To))
					return
				}

				c.chPublishMessage <- msg
				sendACK(msg.ID, "sent")
			}
		})
	})

	if err != nil {
		c.chError <- NewErrorEvent(user.id, "Chat.Register()", err.Error())
		return err
	}

	c.chUserIN <- user
	return nil
}

func (c *Chat) start() {
	users := make(map[string]*User) // all connected users

	pubProducer := c.kafka.NewProducer(config.KafkaNewMessagesTopic())
	defer pubProducer.Close()

	errorProducer := c.kafka.NewProducer(config.KafkaErrorsTopic())
	defer errorProducer.Close()

	for {
		select {
		case msg := <-c.chPublishMessage:
			c.executor.Schedule(func(ctx context.Context) {
				pubProducer.Publish(ctx, []byte(msg.ID), msg.ToJSON())
			})

		case msg := <-c.chSendMessage:
			//user := users[msg.To]
			//c.executor.Schedule(c.messageJob(msg, user))
			c.executor.Schedule(func(ctx context.Context) {
				fmt.Println(msg.String())
			})

		case user := <-c.chUserIN:
			users[user.id] = user
			//bus.Publish(topic.UserStatus, UserPresence{user.id, UserOnline})

		case user := <-c.chUserOUT:
			delete(users, user.id)
			//bus.Publish(topic.UserStatus, UserPresence{user.id, UserOffline})

		case errEvent := <-c.chError:
			c.executor.Schedule(func(ctx context.Context) {
				errorProducer.Publish(ctx, []byte(errEvent.HostID), errEvent.ToJSON())
			})
		}
	}
}

func (c *Chat) consumerMessage() {
	consumer := c.kafka.NewConsumer(config.KafkaGroupID(), config.KafkaHostTopic())
	defer consumer.Close()

	callbackFn := func(event *KafkaEvent, err error) {
		if err != nil {
			c.chError <- NewErrorEvent("", "Chat.consumerMessage()", err.Error())
			return
		}

		//log.Println(string(event.Value))

		var msg message.Message
		if err = json.Unmarshal(event.Value, &msg); err != nil {
			c.chError <- NewErrorEvent("", "Chat.consumerMessage()", err.Error())
			return
		}
		c.chSendMessage <- &msg
	}

	consumer.Subscribe(context.Background(), callbackFn)
}

func (c *Chat) messageJob(msg *message.Message, user *User) func(ctx context.Context) {
	return func(ctx context.Context) {
		handle := func(msg *message.Message) {
			bus := ebus.Instance()

			if msg.IsGroupMessage() {
				bus.Publish(topic.GroupMessage, msg)
				return
			}

			host, ok, err := c.repository.GetUserOnline(msg.To)
			if err != nil {
				c.chError <- NewErrorEvent(msg.To, "Repository.GetUserOnline()", err.Error())
				return
			}
			if !ok {
				bus.Publish(topic.OfflineMessage, msg)
			}

			msg.Host = host
			bus.Publish(topic.SendMessageGRPC, msg)
		}

		if user != nil {
			if err := user.WriteMessage(msg); err != nil {
				handle(msg)
			}
		} else {
			handle(msg)
		}
	}
}
