package chat

import (
	"context"
	"fmt"
	"github.com/tsmweb/chat-service/chat/message"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/pkg/concurrent"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/util/connutil"
	"github.com/tsmweb/go-helper-api/kafka"
	"net"
)

// Chat registers the user's net.Conn connection and handles the data received and sent over the connection.
// It also produces and consumes Apache Kafka data to communicate with the cluster of services.
type Chat struct {
	poller     epoll.EPoll
	executor   concurrent.ExecutorService
	reader     connutil.Reader
	writer     connutil.Writer
	repository Repository
	kafka      kafka.Kafka

	chUserIN            chan *UserConn
	chUserOUT           chan *UserConn
	chPublishMessage    chan *message.Message
	chPublishOffMessage chan *message.Message
	chSendMessage       chan *message.Message
	chPublishError      chan *ErrorEvent

	msgProducer      kafka.Producer
	offMsgProducer   kafka.Producer
	userProducer     kafka.Producer
	errorProducer    kafka.Producer
}

// NewChat creates an instance of Chat.
func NewChat(
	p epoll.EPoll,
	e concurrent.ExecutorService,
	reader connutil.Reader,
	writer connutil.Writer,
	r Repository,
	k kafka.Kafka,
) *Chat {
	chat := &Chat{
		poller:              p,
		executor:            e,
		chUserIN:            make(chan *UserConn),
		chUserOUT:           make(chan *UserConn),
		chPublishMessage:    make(chan *message.Message),
		chPublishOffMessage: make(chan *message.Message),
		chSendMessage:       make(chan *message.Message),
		chPublishError:      make(chan *ErrorEvent),
		reader:              reader,
		writer:              writer,
		repository:          r,
		kafka:               k,
	}

	go chat.start()
	go chat.messageConsumer()

	return chat
}

func (c *Chat) start() {
	c.msgProducer = c.kafka.NewProducer(config.KafkaNewMessagesTopic())
	defer c.msgProducer.Close()

	c.offMsgProducer = c.kafka.NewProducer(config.KafkaOffMessagesTopic())
	defer c.offMsgProducer.Close()

	c.userProducer = c.kafka.NewProducer(config.KafkaUsersTopic())
	defer c.userProducer.Close()

	c.errorProducer = c.kafka.NewProducer(config.KafkaErrorsTopic())
	defer c.errorProducer.Close()

	users := make(map[string]*UserConn) // all connected users

	for {
		select {
		case msg := <-c.chPublishMessage:
			c.executor.Schedule(c.messageProducerJob(msg))

		case msg := <-c.chPublishOffMessage:
			c.executor.Schedule(c.offMessageProducerJob(msg))

		case msg := <-c.chSendMessage:
			user := users[msg.To]
			c.executor.Schedule(c.sendMessageJob(msg, user))

		case user := <-c.chUserIN:
			users[user.userID] = user
			c.executor.Schedule(c.userStatusProducerJob(user.userID, UserOnline))

		case user := <-c.chUserOUT:
			delete(users, user.userID)
			c.executor.Schedule(c.userStatusProducerJob(user.userID, UserOffline))

		case errEvent := <-c.chPublishError:
			c.executor.Schedule(c.errorProducerJob(errEvent))
		}
	}
}

func (c *Chat) messageConsumer() {
	consumer := c.kafka.NewConsumer(config.KafkaGroupID(), config.KafkaHostTopic())
	defer consumer.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			c.chPublishError <- NewErrorEvent("", "Chat.messageConsumer()", err.Error())
			return
		}

		msg := new(message.Message)
		if err = message.Unmarshal(event.Value, msg); err != nil {
			c.chPublishError <- NewErrorEvent("", "Chat.messageConsumer()", err.Error())
			return
		}

		c.chSendMessage <- msg
	}

	consumer.Subscribe(context.Background(), callbackFn)
}

// Register registers the user's net.Conn connection and handles the data received and sent over the connection.
func (c *Chat) Register(userID string, conn net.Conn) error {
	user := &UserConn{
		userID: userID,
		conn:   conn,
		reader: c.reader,
		writer: c.writer,
	}

	observer, err := c.poller.ObservableRead(conn)
	if err != nil {
		c.chPublishError <- NewErrorEvent(user.userID, "Chat.Register()", err.Error())
		return err
	}

	err = observer.Start(func(closed bool, err error) {
		if closed || err != nil {
			observer.Stop()
			c.chUserOUT <- user
			if err != nil {
				c.chPublishError <- NewErrorEvent(user.userID, "Chat.Register()", err.Error())
			}
			return
		}

		c.executor.Schedule(func(ctx context.Context) {
			sendACK := func(msgID, content string) {
				if err := user.WriteACK(msgID, content); err != nil {
					c.chPublishError <- NewErrorEvent(user.userID, "UserConn.WriteACK()", err.Error())
				}
			}

			msg, err := user.Receive()
			if err != nil {
				observer.Stop()
				c.chUserOUT <- user
				c.chPublishError <- NewErrorEvent(user.userID, "UserConn.Receive()", err.Error())
				return
			}
			if msg != nil {
				if msg.IsGroupMessage() {
					c.handleGroupMessage(msg)
				} else {
					ok, err := c.repository.IsValidUser(msg.From, msg.To)
					if err != nil {
						c.chPublishError <- NewErrorEvent(user.userID, "Repository.IsValidUser()", err.Error())
						return
					}
					if !ok {
						sendACK(msg.ID, message.InvalidMessage)
						return
					}

					c.chPublishMessage <- msg
				}

				sendACK(msg.ID, "sent")
			}
		})
	})

	if err != nil {
		c.chPublishError <- NewErrorEvent(user.userID, "Chat.Register()", err.Error())
		return err
	}

	c.chUserIN <- user
	return nil
}

func (c *Chat) handleGroupMessage(msg *message.Message) {
	users, err := c.repository.GetGroupMembers(msg.Group)
	if err != nil {
		c.chPublishError <- NewErrorEvent(msg.From, "Chat.handleGroupMessage()", err.Error())
		return
	}

	for _, user := range users {
		msg.To = user
		c.chPublishMessage <- msg
	}
}

func (c *Chat) messageProducer(ctx context.Context, msg *message.Message, producer kafka.Producer) {
	mpb, err := message.Marshal(msg)
	if err != nil {
		c.chPublishError <- NewErrorEvent(msg.From, "Chat.messageProducer()", err.Error())
		return
	}

	if err = producer.Publish(ctx, []byte(msg.ID), mpb); err != nil {
		c.chPublishError <- NewErrorEvent(msg.From, "Chat.messageProducer()", err.Error())
	}
}

func (c *Chat) messageProducerJob(msg *message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		c.messageProducer(ctx, msg, c.msgProducer)
	}
}

func (c *Chat) offMessageProducerJob(msg *message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		c.messageProducer(ctx, msg, c.offMsgProducer)
	}
}

func (c *Chat) sendMessageJob(msg *message.Message, user *UserConn) func(ctx context.Context) {
	return func(ctx context.Context) {
		fmt.Println(msg)

		if user != nil {
			if err := user.WriteMessage(msg); err != nil {
				c.chPublishOffMessage <- msg
			}
		} else {
			c.chPublishOffMessage <- msg
		}
	}
}

func (c *Chat) userStatusProducerJob(userID string, status UserStatus) func(ctx context.Context) {
	return func(ctx context.Context) {
		var hostID []byte

		if status == UserOnline {
			hostID = []byte(config.HostID())
		} else {
			hostID = []byte("OFF")
		}

		if err := c.userProducer.Publish(ctx, []byte(userID), hostID); err != nil {
			c.chPublishError <- NewErrorEvent(userID, "Chat.userStatusProducerJob()", err.Error())
		}
	}
}

func (c *Chat) errorProducerJob(errEvent *ErrorEvent) func(ctx context.Context) {
	return func(ctx context.Context) {
		c.errorProducer.Publish(ctx, []byte(errEvent.HostID), errEvent.ToJSON())
	}
}
