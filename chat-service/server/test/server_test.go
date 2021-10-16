package test

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tsmweb/chat-service/adapter"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/server"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/chat-service/server/user"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/kafka"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	const (
		UserTest1 = "+5518911111111"
		UserTest2 = "+5518922222222"
		UserTest3 = "+5518933333333"
	)

	c := initServer(t)

	ln := runServerTest(t, c)
	defer ln.Close()

	conn1 := newConnUser(t, ln.Addr().String(), UserTest1)
	conn2 := newConnUser(t, ln.Addr().String(), UserTest2)
	defer func() {
		conn1.Close()
		conn2.Close()
		time.Sleep(time.Second)
	}()

	t.Run("when the message is not valid", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		msg, _ := message.NewMessage(UserTest1, UserTest2, "", message.ContentTypeText, "hello")
		msg.ContentType = ""

		if err := writerConn(conn1, msg); err != nil {
			t.Fatalf("error send message writerConn(): %v", err)
		}

		res := readMessageFromConn(t, conn1)
		t.Log(res)
		assert.Equal(t, res.ID, msg.ID)
		assert.Equal(t, res.Content, message.ErrContentTypeValidateModel.Error())
	})

	t.Run("when UserTest1 was blocked by UserTest3", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		msg, _ := message.NewMessage(UserTest1, UserTest3, "", message.ContentTypeText, "hello")

		if err := writerConn(conn1, msg); err != nil {
			t.Fatalf("error send message writerConn(): %v", err)
		}

		res := readMessageFromConn(t, conn1)
		t.Log(res)
		assert.Equal(t, res.ID, msg.ID)
		assert.Equal(t, res.Content, message.InvalidMessage)
	})

	t.Run("UserTest2 sends message to UserTest1", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		msg, _ := message.NewMessage(UserTest2, UserTest1, "", message.ContentTypeText, "hello test")

		if err := writerConn(conn2, msg); err != nil {
			t.Fatalf("error send message writerConn(): %v", err)
		}

		res := readMessageFromConn(t, conn2)
		t.Log(res)
		assert.Equal(t, res.ID, msg.ID)

		rMSG := readMessageFromConn(t, conn1)
		t.Log(rMSG)
		assert.Equal(t, msg.ID, rMSG.ID)
	})
}

func newConnUser(t *testing.T, addr, userID string) net.Conn {
	t.Helper()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	if err = setUserID(conn, userID); err != nil {
		t.Fatalf("error setUserID(): %v", err)
	}

	return conn
}

func readMessageFromConn(t *testing.T, conn net.Conn) *message.Message {
	t.Helper()
	reader := bufio.NewReader(conn)
	b, _, err := reader.ReadLine()
	if err != nil {
		t.Fatalf("error reader.ReadLine(): %v", err)
	}

	msg := new(message.Message)
	if err = json.Unmarshal(b, msg); err != nil {
		t.Fatalf("error json.Unmarshal(): %v", err)
	}
	return msg
}

func initServer(t *testing.T) *server.Server {
	t.Helper()
	config.Load("../")

	pollerConfig := epoll.ProviderPollerConfig(func(err error) {
		go func(err error) {
			t.Fatal("error netpoll.Config(): ", err)
		}(err)
	})
	poller, err := netpoll.New(pollerConfig)
	if err != nil {
		t.Fatal("error netpoll.New(): ", err)
	}
	ePoll := epoll.NewEPoll(poller)
	reader := server.ConnReaderFunc(readerConn)
	writer := server.ConnWriterFunc(writerConn)
	msgEncoder := message.EncoderFunc(adapter.MessageMarshal)
	msgDecoder := message.DecoderFunc(adapter.MessageUnmarshal)
	userEncoder := user.EncoderFunc(adapter.UserMarshal)
	eventErrorEncoder := server.ErrorEventEncoderFunc(adapter.ErrorEventMarshal)
	kaf := kafka.New([]string{config.KafkaBootstrapServers()}, config.KafkaClientID())
	consumeMessage := kaf.NewConsumer(config.KafkaGroupID(), config.KafkaHostTopic())
	handleMessage := server.NewHandleMessage(msgEncoder, kaf.NewProducer(config.KafkaNewMessagesTopic()))
	handleOffMessage := server.NewHandleMessage(msgEncoder, kaf.NewProducer(config.KafkaOffMessagesTopic()))
	handleUserStatus := server.NewHandleUserStatus(userEncoder,
		kaf.NewProducer(config.KafkaUsersTopic()),
		kaf.NewProducer(config.KafkaUsersPresenceTopic()))
	handleError := server.NewHandleError(eventErrorEncoder, kaf.NewProducer(config.KafkaErrorsTopic()))

	serv := server.NewServer(
		context.Background(),
		ePoll,
		reader,
		writer,
		msgDecoder,
		consumeMessage,
		handleMessage,
		handleOffMessage,
		handleUserStatus,
		handleError)

	return serv
}

func readerConn(conn net.Conn) (io.Reader, error) {
	return conn, nil
}

func writerConn(conn net.Conn, data interface{}) (err error) {
	encoder := json.NewEncoder(conn)
	if err = encoder.Encode(data); err != nil {
		return
	}
	return
}

func runServerTest(tb testing.TB, chat *server.Server) net.Listener {
	ln, err := net.Listen("tcp", "localhost:")
	if err != nil {
		tb.Fatal(err)
		return nil
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					// Server closed.
					return
				}
				tb.Fatal(err)
			}

			userID, err := getUserID(conn)
			if err != nil || userID == "" {
				conn.Close()
				continue
			}

			if err := chat.Register(userID, conn); err != nil {
				tb.Fatalf("failed to register user on server server: %v", err)
			}
		}
	}()
	return ln
}

func getUserID(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	b, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func setUserID(conn net.Conn, ID string) error {
	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString(ID + "\n"); err != nil {
		return err
	}
	writer.Flush()
	return nil
}
