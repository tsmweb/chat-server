package chat_test

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tsmweb/chat-service/chat"
	"github.com/tsmweb/chat-service/chat/message"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/pkg/concurrent"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/util/connutil"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"github.com/tsmweb/go-helper-api/kafka"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestChat(t *testing.T) {
	executor := executor.New(100)
	defer executor.Shutdown()

	c := initChat(t, executor)

	ln := runServerTest(t, c)
	defer ln.Close()

	// user 1
	conn1 := newConnUser(t, ln.Addr().String(), chat.UserTest1)
	defer func() {
		conn1.Close()
		time.Sleep(time.Second)
	}()

	t.Run("when the message is not valid", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		msg, _ := message.NewMessage(chat.UserTest1, chat.UserTest2, "", message.ContentText, "hello")
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
		msg, _ := message.NewMessage(chat.UserTest1, chat.UserTest3, "", message.ContentText, "hello")

		if err := writerConn(conn1, msg); err != nil {
			t.Fatalf("error send message writerConn(): %v", err)
		}

		res := readMessageFromConn(t, conn1)
		t.Log(res)
		assert.Equal(t, res.ID, msg.ID)
		assert.Equal(t, res.Content, message.InvalidMessage)
	})

	t.Run("send message to UserTest1", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		msg, _ := message.NewMessage(chat.UserTest2, chat.UserTest1, "", message.ContentText, "hello")

		if err := writerConn(conn1, msg); err != nil {
			t.Fatalf("error send message writerConn(): %v", err)
		}

		res := readMessageFromConn(t, conn1)

		t.Log(res)
		assert.Equal(t, msg.ID, res.ID)
	})

	/*
	t.Run("UserTest2 sends message to UserTest1", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		conn2 := newConnUser(t, ln.Addr().String(), chat.UserTest2)
		defer conn2.Close()

		sMSG, _ := chat.New(chat.UserTest2, chat.UserTest1, "", chat.ContentText, "hello test")

		if err := writerConn(conn2, sMSG); err != nil {
			t.Fatalf("error send message writerConn(): %v", err)
		}

		res := readMessageFromConn(t, conn2)
		t.Log(res)
		assert.Equal(t, res.ID, sMSG.ID)

		rMSG := readMessageFromConn(t, conn1)
		t.Log(rMSG)
		assert.Equal(t, sMSG.ID, rMSG.ID)
	})
	*/
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

func initChat(t *testing.T, executor concurrent.ExecutorService) *chat.Chat {
	t.Helper()
	config.Load("../")

	pollerConfig := epoll.ProviderPollerConfig(func(err error) {
		executor.Schedule(func(ctx context.Context) {
			t.Fatal("error netpoll.Config(): ", err)
		})
	})
	poller, err := netpoll.New(pollerConfig)
	if err != nil {
		t.Fatal("error netpoll.New(): ", err)
	}
	ePoll := epoll.NewEPoll(poller)
	repository := chat.NewMemoryRepository()
	reader := connutil.FuncReader(readerConn)
	writer := connutil.FuncWriter(writerConn)
	Kaf := kafka.New([]string{config.KafkaBootstrapServers()}, config.KafkaClientID())

	chat := chat.NewChat(
		ePoll,
		executor,
		reader,
		writer,
		repository,
		Kaf)

	return chat
}

func readerConn(conn io.ReadWriter) (io.Reader, error) {
	return conn, nil
}

func writerConn(conn io.Writer, x interface{}) (err error) {
	encoder := json.NewEncoder(conn)
	if err = encoder.Encode(x); err != nil {
		return
	}
	return
}

func runServerTest(tb testing.TB, chat *chat.Chat) net.Listener {
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
				tb.Fatalf("failed to register user on chat server: %v", err)
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
