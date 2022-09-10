package server

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/tsmweb/chat-service/common/service"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/chat-service/server/user"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"github.com/tsmweb/go-helper-api/kafka"
)

// Server registers the user's net.Conn connection and handles the data received and sent over
// the connection.
// It also produces and consumes Apache Kafka data to communicate with the cluster of services.
type Server struct {
	ctx             context.Context
	execMessageRecv *executor.Executor
	execMessageSent *executor.Executor
	poller          epoll.EPoll

	chUserIN       chan *UserConn
	chUserOUT      chan string
	chMessageRecv  chan *message.Message
	chMessageSent  chan *message.Message
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
		ctx:              ctx,
		poller:           poll,
		chUserIN:         make(chan *UserConn),
		chUserOUT:        make(chan string),
		chMessageRecv:    make(chan *message.Message),
		chMessageSent:    make(chan *message.Message),
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
		service.Error(userConn.userID, "Server.Register()", err)
		return err
	}

	err = observer.Start(func(closed bool, errPoller error) {
		if closed || errPoller != nil {
			observer.Stop()
			s.chUserOUT <- userConn.userID
			if errPoller != nil {
				service.Error(userConn.userID, "Server.Register():poller", errPoller)
			}
			return
		}

		s.execMessageRecv.Schedule(func(ctx context.Context) {
			msg, err := userConn.Receive() // receive message from userConn connection
			if err != nil {
				observer.Stop()
				s.chUserOUT <- userConn.userID
				return
			}
			if msg != nil {
				s.chMessageRecv <- msg

				if err = userConn.WriteResponse(msg.ID, message.ContentTypeACK, message.AckMessage); err != nil {
					service.Error(userConn.userID, "UserConn.WriteACK()", err)
				}
			}
		})
	})

	if err != nil {
		service.Error(userConn.userID, "Server.Register()", err)
		return err
	}

	s.chUserIN <- userConn
	return nil
}

func (s *Server) run() {
	// Executor to perform background processing,
	// limiting resource consumption when executing a collection of jobs.
	s.execMessageRecv = executor.New(config.GoPoolSize())
	s.execMessageSent = executor.New(config.GoPoolSize())

	go s.messageProcessor()
	go s.messageConsumer()
}

func (s *Server) stop() {
	s.execMessageRecv.Shutdown()
	s.execMessageSent.Shutdown()

	s.handleMessage.Close()
	s.handleOffMessage.Close()
	s.handleUserStatus.Close()
}

func (s *Server) messageProcessor() {
	users := make(map[string]*UserConn) // all connected users

loop:
	for {
		select {
		case msg := <-s.chMessageRecv:
			s.messageRecvTask(msg)

		case msg := <-s.chMessageSent:
			userConn := users[msg.To]
			s.messageSendTask(msg, userConn)

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
	defer s.consumeMessage.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil && err.Error() != "nil" {
			service.Error("", "Server.messageConsumer()", err)
			return
		}

		var msg message.Message
		if err = s.msgDecoder.Unmarshal(event.Value, &msg); err != nil {
			service.Error("", "Server.messageConsumer()", err)
			return
		}

		s.chMessageSent <- &msg
	}

	s.consumeMessage.Subscribe(s.ctx, callbackFn)
}

func (s *Server) messageRecvTask(msg *message.Message) {
	go func() {
		if err := s.handleMessage.Execute(s.ctx, msg); err != nil {
			service.Error(msg.From, "Server.handleMessage()", err)
		}
	}()
}

func (s *Server) messageSendTask(msg *message.Message, userConn *UserConn) {
	s.execMessageSent.Schedule(func(ctx context.Context) {
		if userConn != nil {
			// write a message on user connection
			if err := userConn.WriteMessage(msg); err != nil {
				s.sendOffMessage(ctx, msg)
			}
		} else {
			s.sendOffMessage(ctx, msg)
		}
	})
}

func (s *Server) sendOffMessage(ctx context.Context, msg *message.Message) {
	var contentStatus = message.ContentTypeStatus
	if msg.ContentType == contentStatus.String() {
		return
	}

	if err := s.handleOffMessage.Execute(ctx, msg); err != nil {
		service.Error(msg.From, "Server.handleOffMessage()", err)
	}
}

func (s *Server) userStatusTask(userID string, status user.Status) {
	go func() {
		if err := s.handleUserStatus.Execute(s.ctx, userID, status); err != nil {
			service.Error(userID, "Server.handleUserStatus()", err)
		}
	}()
}
