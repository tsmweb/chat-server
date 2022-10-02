package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/tsmweb/chat-service/common/service"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/chat-service/server/user"
	"github.com/tsmweb/go-helper-api/concurrent/gopool"
	"github.com/tsmweb/go-helper-api/kafka"
)

// Server registers the user's net.Conn connection and handles the data received and sent over
// the connection.
// It also produces and consumes Apache Kafka data to communicate with the cluster of services.
type Server struct {
	tag              string
	ctx              context.Context
	poller           epoll.EPoll
	poolUsers        *gopool.Pool
	poolSendMessages *gopool.Pool
	poolRecvMessages *gopool.Pool

	chUserIN       chan *UserConn
	chUserOUT      chan string
	chRecvMessage  chan message.Message
	connReader     ConnReader
	connWriter     ConnWriter
	msgDecoder     message.Decoder
	consumeMessage kafka.Consumer

	handleMessage    HandleMessage
	handleOffMessage HandleMessage
	handleUserStatus HandleUserStatus
}

// NewServer creates an instance of Server.
func NewServer(
	ctx context.Context,
	poll epoll.EPoll,
	connReader ConnReader,
	connWriter ConnWriter,
	msgDecoder message.Decoder,
	consumeMessage kafka.Consumer,
	handleMessage HandleMessage,
	handleOffMessage HandleMessage,
	handleUserStatus HandleUserStatus,
) *Server {
	server := &Server{
		tag:              "server::Server",
		ctx:              ctx,
		poller:           poll,
		chUserIN:         make(chan *UserConn),
		chUserOUT:        make(chan string),
		chRecvMessage:    make(chan message.Message),
		connReader:       connReader,
		connWriter:       connWriter,
		msgDecoder:       msgDecoder,
		consumeMessage:   consumeMessage,
		handleMessage:    handleMessage,
		handleOffMessage: handleOffMessage,
		handleUserStatus: handleUserStatus,
	}

	server.run()

	return server
}

// Register registers the user's net.Conn connection and handles the data received and sent over
// the connection.
func (s *Server) Register(userID string, conn net.Conn) error {
	userConn := &UserConn{
		userID: userID,
		conn:   conn,
		reader: s.connReader,
		writer: s.connWriter,
	}

	var fdConn net.Conn

	switch conn.(type) {
	case *tls.Conn:
		fdConn = conn.(*tls.Conn).NetConn()
	default:
		fdConn = conn
	}

	observer, err := s.poller.ObservableRead(fdConn)
	if err != nil {
		service.Error(userConn.userID, s.tag,
			fmt.Errorf("epoll::EPoll: %s", err.Error()))
		return err
	}

	err = observer.Start(func(closed bool, errPoller error) {
		if closed || errPoller != nil {
			observer.Stop()
			s.chUserOUT <- userConn.userID
			if errPoller != nil {
				service.Error(userConn.userID, s.tag,
					fmt.Errorf("epoll::Observer: %s", errPoller.Error()))
			}
			return
		}

		s.poolSendMessages.Schedule(func(ctx context.Context) {
			msg, err := userConn.Receive() // receive message from userConn connection
			if err != nil {
				observer.Stop()
				s.chUserOUT <- userConn.userID
				return
			}
			if msg != nil {
				if err := s.handleMessage.Execute(s.ctx, msg); err != nil {
					service.Error(msg.From, s.tag, err)

					if err = userConn.WriteResponse(
						msg.ID,
						message.ContentTypeError,
						"internal server error",
					); err != nil {
						service.Error(userConn.userID, s.tag,
							fmt.Errorf("server::UserConn: %s", err.Error()))
					}
					return
				}

				if err = userConn.WriteResponse(
					msg.ID,
					message.ContentTypeACK,
					message.AckMessage,
				); err != nil {
					service.Error(userConn.userID, s.tag,
						fmt.Errorf("server::UserConn: %s", err.Error()))
				}
			}
		})
	})

	if err != nil {
		service.Error(userConn.userID, s.tag, err)
		return err
	}

	s.chUserIN <- userConn
	return nil
}

func (s *Server) run() {
	// Executor to perform background processing,
	// limiting resource consumption when executing a collection of jobs.
	workerSize := config.GoPoolSize()
	queueSize := 1

	s.poolUsers = gopool.New(workerSize, queueSize)
	s.poolSendMessages = gopool.New(workerSize, queueSize)
	s.poolRecvMessages = gopool.New(workerSize, queueSize)

	go s.messageProcessor()
	go s.messageConsumer()
}

func (s *Server) stop() {
	s.poolUsers.Close()
	s.poolSendMessages.Close()
	s.poolRecvMessages.Close()

	s.handleMessage.Close()
	s.handleOffMessage.Close()
	s.handleUserStatus.Close()
}

func (s *Server) messageProcessor() {
	users := make(map[string]*UserConn) // all connected users

loop:
	for {
		select {
		case msg := <-s.chRecvMessage:
			userConn := users[msg.To]
			s.recvMessageTask(msg, userConn)

		case u := <-s.chUserIN:
			users[u.userID] = u
			s.userStatusTask(u.userID, user.Online)

		case userID := <-s.chUserOUT:
			delete(users, userID)
			s.userStatusTask(userID, user.Offline)

		case <-s.ctx.Done():
			break loop
		}
	}

	s.stop()
}

func (s *Server) messageConsumer() {
	defer func() {
		s.consumeMessage.Close()
		log.Println("[STOP] Server::consumeMessage")
	}()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil && err.Error() != "nil" {
			service.Error("", s.tag, fmt.Errorf("kafka::Consumer: %s", err.Error()))
			return
		}

		var msg message.Message
		if err = s.msgDecoder.Unmarshal(event.Value, &msg); err != nil {
			service.Error("", s.tag, fmt.Errorf("kafka::Consumer: %s", err.Error()))
			return
		}

		s.chRecvMessage <- msg
	}

	s.consumeMessage.Subscribe(s.ctx, callbackFn)
}

func (s *Server) recvMessageTask(msg message.Message, userConn *UserConn) {
	s.poolRecvMessages.Schedule(func(ctx context.Context) {
		if userConn != nil {
			// write a message on user connection
			if err := userConn.WriteMessage(&msg); err != nil {
				s.sendOffMessage(ctx, &msg)
			}
		} else {
			s.sendOffMessage(ctx, &msg)
		}
	})
}

func (s *Server) sendOffMessage(ctx context.Context, msg *message.Message) {
	var contentStatus = message.ContentTypeStatus
	if msg.ContentType == contentStatus.String() {
		return
	}

	if err := s.handleOffMessage.Execute(ctx, msg); err != nil {
		service.Error(msg.From, s.tag,
			fmt.Errorf("server::HandleOffMessage: %s", err.Error()))
	}
}

func (s *Server) userStatusTask(userID string, status user.Status) {
	s.poolUsers.Schedule(func(ctx context.Context) {
		if err := s.handleUserStatus.Execute(s.ctx, userID, status); err != nil {
			service.Error(userID, s.tag,
				fmt.Errorf("server::HandleUserStatus: %s", err.Error()))
		}
	})
}
