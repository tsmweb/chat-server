package chat

import (
	"context"
	"fmt"
	"github.com/tsmweb/chat-service/chat/message"
	"github.com/tsmweb/chat-service/chat/user"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"github.com/tsmweb/go-helper-api/kafka"
	"net"
)

// Server registers the user's net.Conn connection and handles the data received and sent over the connection.
// It also produces and consumes Apache Kafka data to communicate with the cluster of services.
type Server struct {
	ctx         context.Context
	executor    *executor.Executor
	poller      epoll.EPoll
	connReader  ConnReader
	connWriter  ConnWriter
	msgEncoder  message.Encoder
	msgDecoder  message.Decoder
	userEncoder user.Encoder
	repository  Repository
	kafka       kafka.Kafka

	chUserIN       chan *UserConn
	chUserOUT      chan *UserConn
	chMessage      chan *message.Message
	chOffMessage   chan *message.Message
	chGroupMessage chan *message.Message
	chSendMessage  chan *message.Message
	chError        chan *ErrorEvent

	msgProducer    kafka.Producer
	offMsgProducer kafka.Producer
	userProducer   kafka.Producer
	errorProducer  kafka.Producer
}

// NewServer creates an instance of Server.
func NewServer(
	ctx context.Context,
	p epoll.EPoll,
	connReader ConnReader,
	connWriter ConnWriter,
	msgEncoder message.Encoder,
	msgDecoder message.Decoder,
	userEncoder user.Encoder,
	r Repository,
	k kafka.Kafka,
) *Server {
	server := &Server{
		ctx:            ctx,
		poller:         p,
		chUserIN:       make(chan *UserConn),
		chUserOUT:      make(chan *UserConn),
		chMessage:      make(chan *message.Message),
		chOffMessage:   make(chan *message.Message),
		chGroupMessage: make(chan *message.Message),
		chSendMessage:  make(chan *message.Message),
		chError:        make(chan *ErrorEvent),
		connReader:     connReader,
		connWriter:     connWriter,
		msgEncoder:     msgEncoder,
		msgDecoder:     msgDecoder,
		userEncoder:    userEncoder,
		repository:     r,
		kafka:          k,
	}

	server.run()

	return server
}

func (s *Server) run() {
	// Executor to perform background processing,
	// limiting resource consumption when executing a collection of jobs.
	s.executor = executor.New(config.GoPoolSize())
	s.msgProducer = s.kafka.NewProducer(config.KafkaNewMessagesTopic())
	s.offMsgProducer = s.kafka.NewProducer(config.KafkaOffMessagesTopic())
	s.userProducer = s.kafka.NewProducer(config.KafkaUsersTopic())
	s.errorProducer = s.kafka.NewProducer(config.KafkaErrorsTopic())

	go s.messageProcessor()
	go s.messageConsumer()
}

func (s *Server) stop() {
	s.executor.Shutdown()
	s.msgProducer.Close()
	s.offMsgProducer.Close()
	s.userProducer.Close()
	s.errorProducer.Close()
}

func (s *Server) messageProcessor() {
	users := make(map[string]*UserConn) // all connected users

LOOP:
	for {
		select {
		case msg := <-s.chMessage:
			s.executor.Schedule(s.messageJob(msg))

		case msg := <-s.chOffMessage:
			s.executor.Schedule(s.offMessageJob(msg))

		case msg := <-s.chGroupMessage:
			s.executor.Schedule(s.groupMessageJob(msg))

		case msg := <-s.chSendMessage:
			user := users[msg.To]
			s.executor.Schedule(s.sendMessageJob(msg, user))

		case u := <-s.chUserIN:
			users[u.userID] = u
			s.executor.Schedule(s.userStatusJob(u.userID, user.Online))

		case u := <-s.chUserOUT:
			delete(users, u.userID)
			s.executor.Schedule(s.userStatusJob(u.userID, user.Offline))

		case errEvent := <-s.chError:
			s.executor.Schedule(s.errorJob(errEvent))

		case <-s.ctx.Done():
			s.stop()
			break LOOP
		}
	}
}

func (s *Server) messageConsumer() {
	consumer := s.kafka.NewConsumer(config.KafkaGroupID(), config.KafkaHostTopic())
	defer consumer.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			s.chError <- NewErrorEvent("", "Server.messageConsumer()", err.Error())
			return
		}

		msg := new(message.Message)
		if err = s.msgDecoder.Unmarshal(event.Value, msg); err != nil {
			s.chError <- NewErrorEvent("", "Server.messageConsumer()", err.Error())
			return
		}

		s.chSendMessage <- msg
	}

	consumer.Subscribe(s.ctx, callbackFn)
}

// Register registers the user's net.Conn connection and handles the data received and sent over the connection.
func (s *Server) Register(userID string, conn net.Conn) error {
	user := &UserConn{
		userID: userID,
		conn:   conn,
		reader: s.connReader,
		writer: s.connWriter,
	}

	observer, err := s.poller.ObservableRead(conn)
	if err != nil {
		s.chError <- NewErrorEvent(user.userID, "Server.Register()", err.Error())
		return err
	}

	err = observer.Start(func(closed bool, err error) {
		if closed || err != nil {
			observer.Stop()
			s.chUserOUT <- user
			if err != nil {
				s.chError <- NewErrorEvent(user.userID, "Server.Register()", err.Error())
			}
			return
		}

		s.executor.Schedule(func(ctx context.Context) {
			sendACK := func(msgID, content string) {
				if err := user.WriteACK(msgID, content); err != nil {
					s.chError <- NewErrorEvent(user.userID, "UserConn.WriteACK()", err.Error())
				}
			}

			msg, err := user.Receive()
			if err != nil {
				observer.Stop()
				s.chUserOUT <- user
				s.chError <- NewErrorEvent(user.userID, "UserConn.Receive()", err.Error())
				return
			}
			if msg != nil {
				if msg.IsGroupMessage() {
					s.chGroupMessage <- msg
				} else {
					ok, err := s.repository.IsValidUser(msg.From, msg.To)
					if err != nil {
						s.chError <- NewErrorEvent(user.userID, "Repository.IsValidUser()", err.Error())
						return
					}
					if !ok {
						sendACK(msg.ID, message.InvalidMessage)
						return
					}

					s.chMessage <- msg
				}

				sendACK(msg.ID, "sent")
			}
		})
	})

	if err != nil {
		s.chError <- NewErrorEvent(user.userID, "Server.Register()", err.Error())
		return err
	}

	s.chUserIN <- user
	return nil
}

func (s *Server) messageJob(msg *message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		mpb, err := s.msgEncoder.Marshal(msg)
		if err != nil {
			s.chError <- NewErrorEvent(msg.From, "Server.messageJob()", err.Error())
			return
		}

		if err = s.msgProducer.Publish(ctx, []byte(msg.ID), mpb); err != nil {
			s.chError <- NewErrorEvent(msg.From, "Server.messageJob()", err.Error())
			return
		}
	}
}

func (s *Server) offMessageJob(msg *message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		mpb, err := s.msgEncoder.Marshal(msg)
		if err != nil {
			s.chError <- NewErrorEvent(msg.From, "Server.offMessageJob()", err.Error())
			return
		}

		if err = s.offMsgProducer.Publish(ctx, []byte(msg.ID), mpb); err != nil {
			s.chError <- NewErrorEvent(msg.From, "Server.offMessageJob()", err.Error())
			return
		}
	}
}

func (s *Server) groupMessageJob(msg *message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		users, err := s.repository.GetGroupMembers(msg.Group)
		if err != nil {
			s.chError <- NewErrorEvent(msg.From, "Server.groupMessageJob()", err.Error())
			return
		}

		for _, user := range users {
			msg.To = user
			s.chMessage <- msg
		}
	}
}

func (s *Server) sendMessageJob(msg *message.Message, user *UserConn) func(ctx context.Context) {
	return func(ctx context.Context) {
		fmt.Println(msg)

		if user != nil {
			if err := user.WriteMessage(msg); err != nil {
				s.chOffMessage <- msg
			}
		} else {
			s.chOffMessage <- msg
		}
	}
}

func (s *Server) userStatusJob(userID string, status user.Status) func(ctx context.Context) {
	return func(ctx context.Context) {
		serverID := "OFF"

		if status == user.Online {
			serverID = config.HostID()
		}

		u := user.NewUser(userID, status, serverID)
		upb, err := s.userEncoder.Marshal(u)
		if err != nil {
			s.chError <- NewErrorEvent(userID, "Server.userStatusJob()", err.Error())
			return
		}

		if err := s.userProducer.Publish(ctx, []byte(userID), upb); err != nil {
			s.chError <- NewErrorEvent(userID, "Server.userStatusJob()", err.Error())
			return
		}
	}
}

func (s *Server) errorJob(errEvent *ErrorEvent) func(ctx context.Context) {
	return func(ctx context.Context) {
		s.errorProducer.Publish(ctx, []byte(errEvent.HostID), errEvent.ToJSON())
	}
}
