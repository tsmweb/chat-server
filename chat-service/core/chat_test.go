package core

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsmweb/chat-service/common/concurrent"
	"github.com/tsmweb/chat-service/common/connutil"
	"github.com/tsmweb/chat-service/common/epoll"
	"github.com/tsmweb/chat-service/core/ctype"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestChat(t *testing.T) {
	executor := executor.New(100)
	defer executor.Shutdown()

	chat := initChat(t, executor)

	ln := runServerTest(t, chat)
	defer ln.Close()

	// user 1
	conn1 := newConnUser(t, ln.Addr().String(), userTest1)
	defer conn1.Close()

	t.Run("chat.SendMessage sends message to userTest1", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		sMSG, _ := NewMessage(userTest2, userTest1, "", ctype.TEXT, "hello")

		if err := chat.SendMessage(sMSG); err != nil {
			t.Fatalf("error chat.SendMessage(): %v", err)
		}

		rMSG := readMessageFromConn(t, conn1)

		t.Log(rMSG)
		assert.Equal(t, sMSG.ID, rMSG.ID)
	})
	
	t.Run("when the message is not valid", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		msg, _ := NewMessage(userTest1, userTest2, "", ctype.TEXT, "hello")
		msg.ContentType = ""

		if err := writerConn(conn1, msg); err != nil {
			t.Fatalf("error send message writerConn(): %v", err)
		}

		res := readMessageFromConn(t, conn1)
		t.Log(res)
		assert.Equal(t, res.ID, msg.ID)
		assert.Equal(t, res.Content, ErrContentTypeValidateModel.Error())
	})

	t.Run("when userTest1 was blocked by userTest3", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		msg, _ := NewMessage(userTest1, userTest3, "", ctype.TEXT, "hello")

		if err := writerConn(conn1, msg); err != nil {
			t.Fatalf("error send message writerConn(): %v", err)
		}

		res := readMessageFromConn(t, conn1)
		t.Log(res)
		assert.Equal(t, res.ID, msg.ID)
		assert.Equal(t, res.Content, fmt.Sprintf(BlockedMessage, msg.To))
	})

	t.Run("userTest2 sends message to userTest1", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		conn2 := newConnUser(t, ln.Addr().String(), userTest2)
		defer conn2.Close()

		sMSG, _ := NewMessage(userTest2, userTest1, "", ctype.TEXT, "hello test")

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

func readMessageFromConn(t *testing.T, conn net.Conn) *Message {
	t.Helper()
	reader := bufio.NewReader(conn)
	b, _, err := reader.ReadLine()
	if err != nil {
		t.Fatalf("error reader.ReadLine(): %v", err)
	}

	msg := new(Message)
	if err = json.Unmarshal(b, msg); err != nil {
		t.Fatalf("error json.Unmarshal(): %v", err)
	}
	return msg
}

func initChat(t *testing.T, executor concurrent.ExecutorService) *Chat {
	t.Helper()
	errorDispatcher := NewErrorDispatcher()
	config := epoll.ProviderPollerConfig(func(err error) {
		executor.Schedule(func(ctx context.Context) {
			t.Fatal("error netpoll.Config(): ", err)
		})
	})
	poller, err := netpoll.New(config)
	if err != nil {
		t.Fatal("error netpoll.New(): ", err)
	}
	ePoll := epoll.NewEPoll(poller)
	repository := NewMemoryRepository()
	reader := connutil.FuncReader(readerConn)
	writer := connutil.FuncWriter(writerConn)
	presenceDispatcher := NewPresenceDispatcher()
	userStatusHandler := NewUserStatusHandler(repository, presenceDispatcher)
	offlineMessageDispatcher := NewOfflineMessageDispatcher()
	groupMessageDispatcher := NewGroupMessageDispatcher()
	messageHandler := NewMessageHandler(repository, offlineMessageDispatcher, groupMessageDispatcher)
	chat := NewChat(
		ePoll,
		executor,
		"localhost",
		reader,
		writer,
		userStatusHandler,
		messageHandler,
		errorDispatcher)

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

func runServerTest(tb testing.TB, chat *Chat) net.Listener {
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