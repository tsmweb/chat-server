package server

import (
	"context"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/chat-service/server/user"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"github.com/tsmweb/go-helper-api/kafka"
	"net"
)

// Server registers the user's net.Conn connection and handles the data received and sent over the connection.
// It also produces and consumes Apache Kafka data to communicate with the cluster of services.
type Server struct {
	ctx      context.Context
	executor *executor.Executor
	poller   epoll.EPoll

	chUserIN       chan *UserConn
	chUserOUT      chan string
	chMessage      chan message.Message
	chGroupMessage chan message.Message
	chSendMessage  chan message.Message
	chError        chan ErrorEvent

	connReader     ConnReader
	connWriter     ConnWriter
	msgDecoder     message.Decoder
	consumeMessage kafka.Consumer

	handleMessage    HandleMessage
	handleOffMessage HandleMessage
	handleUserStatus HandleUserStatus
	handleError      HandleError
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
	handleError HandleError,
) *Server {
	server := &Server{
		ctx:              ctx,
		poller:           poll,
		chUserIN:         make(chan *UserConn),
		chUserOUT:        make(chan string),
		chMessage:        make(chan message.Message),
		chGroupMessage:   make(chan message.Message),
		chSendMessage:    make(chan message.Message),
		chError:          make(chan ErrorEvent),
		connReader:       connReader,
		connWriter:       connWriter,
		msgDecoder:       msgDecoder,
		consumeMessage:   consumeMessage,
		handleMessage:    handleMessage,
		handleOffMessage: handleOffMessage,
		handleUserStatus: handleUserStatus,
		handleError:      handleError,
	}

	server.run()

	return server
}

// Register registers the user's net.Conn connection and handles the data received and sent over the connection.
func (s *Server) Register(userID string, conn net.Conn) error {
	userConn := &UserConn{
		userID: userID,
		conn:   conn,
		reader: s.connReader,
		writer: s.connWriter,
	}

	observer, err := s.poller.ObservableRead(conn)
	if err != nil {
		s.chError <- *NewErrorEvent(userConn.userID, "Server.Register()", err.Error())
		return err
	}

	err = observer.Start(func(closed bool, errPoller error) {
		if closed || errPoller != nil {
			observer.Stop()
			s.chUserOUT <- userConn.userID
			if errPoller != nil {
				s.chError <- *NewErrorEvent(userConn.userID, "Server.Register():poller", errPoller.Error())
			}
			return
		}

		s.executor.Schedule(func(ctx context.Context) {
			msg, err := userConn.Receive() // receive message from userConn connection
			if err != nil {
				observer.Stop()
				s.chUserOUT <- userConn.userID
				s.chError <- *NewErrorEvent(userConn.userID, "UserConn.Receive()", err.Error())
				return
			}
			if msg != nil {
				s.chMessage <- *msg

				if err = userConn.WriteResponse(msg.ID, message.ContentTypeACK, message.AckMessage); err != nil {
					s.chError <- *NewErrorEvent(userConn.userID, "UserConn.WriteACK()", err.Error())
				}
			}
		})
	})

	if err != nil {
		s.chError <- *NewErrorEvent(userConn.userID, "Server.Register()", err.Error())
		return err
	}

	s.chUserIN <- userConn
	return nil
}

func (s *Server) run() {
	// Executor to perform background processing,
	// limiting resource consumption when executing a collection of jobs.
	s.executor = executor.New(config.GoPoolSize())

	go s.messageProcessor()
	go s.messageConsumer()
}

func (s *Server) stop() {
	s.executor.Shutdown()

	s.handleMessage.Close()
	s.handleOffMessage.Close()
	s.handleUserStatus.Close()
	s.handleError.Close()
}

func (s *Server) messageProcessor() {
	users := make(map[string]*UserConn) // all connected users

loop:
	for {
		select {
		case msg := <-s.chMessage:
			s.executor.Schedule(s.messageTask(msg))

		case msg := <-s.chSendMessage:
			userConn := users[msg.To]
			s.executor.Schedule(s.sendMessageTask(msg, userConn))

		case u := <-s.chUserIN:
			users[u.userID] = u
			s.executor.Schedule(s.userStatusTask(u.userID, user.Online))

		case userID := <-s.chUserOUT:
			delete(users, userID)
			s.executor.Schedule(s.userStatusTask(userID, user.Offline))

		case err := <-s.chError:
			s.executor.Schedule(s.errorTask(err))

		case <-s.ctx.Done():
			break loop
		}
	}

	s.stop()
}

func (s *Server) messageConsumer() {
	defer s.consumeMessage.Close()

	callbackFn := func(event *kafka.Event, err error) {
		if err != nil {
			s.chError <- *NewErrorEvent("", "Server.messageConsumer()", err.Error())
			return
		}

		var msg message.Message
		if err = s.msgDecoder.Unmarshal(event.Value, &msg); err != nil {
			s.chError <- *NewErrorEvent("", "Server.messageConsumer()", err.Error())
			return
		}

		s.chSendMessage <- msg
	}

	s.consumeMessage.Subscribe(s.ctx, callbackFn)
}

func (s *Server) messageTask(msg message.Message) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := s.handleMessage.Execute(ctx, msg); err != nil {
			s.chError <- *err
		}
	}
}

func (s *Server) sendMessageTask(msg message.Message, userConn *UserConn) func(ctx context.Context) {
	return func(ctx context.Context) {
		if userConn != nil {
			// write a message on user connection
			if err := userConn.WriteMessage(&msg); err != nil {
				s.sendOffMessage(ctx, msg)
			}
		} else {
			s.sendOffMessage(ctx, msg)
		}
	}
}

func (s *Server) sendOffMessage(ctx context.Context, msg message.Message) {
	var contentStatus = message.ContentTypeStatus
	if msg.ContentType == contentStatus.String() {
		return
	}

	if err := s.handleOffMessage.Execute(ctx, msg); err != nil {
		s.chError <- *err
	}
}

func (s *Server) userStatusTask(userID string, status user.Status) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := s.handleUserStatus.Execute(ctx, userID, status); err != nil {
			s.chError <- *err
		}
	}
}

func (s *Server) errorTask(errEvent ErrorEvent) func(ctx context.Context) {
	return func(ctx context.Context) {
		s.handleError.Execute(ctx, errEvent)
	}
}
