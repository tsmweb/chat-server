package core

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tsmweb/chat-service/common/concurrent"
	"github.com/tsmweb/chat-service/common/connutil"
	"github.com/tsmweb/chat-service/common/epoll"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestChat(t *testing.T) {
	user1 := "+5518977777777"
	user2 := "+5518966666666"

	executor := executor.New(100)
	defer executor.Shutdown()

	chat := initChat(t, executor)

	ln := runServerTest(t, chat)
	defer ln.Close()

	// user 1
	conn1 := newConnUser(t, ln.Addr().String(), user1)
	defer conn1.Close()

	t.Run("chat.SendMessage sends message to user1", func(t *testing.T) {
		time.Sleep(10 * time.Millisecond)

		sMSG, _ := NewMessage(user2, user1, "", TEXT, "hello")
		if err := chat.SendMessage(sMSG); err != nil {
			t.Fatalf("error chat.SendMessage(): %v", err)
		}

		rMSG := readMessageFromConn(t, conn1)

		t.Log(rMSG)
		assert.Equal(t, sMSG.ID, rMSG.ID)
	})

	t.Run("user2 sends message to user1", func(t *testing.T) {
		// user 2
		conn2 := newConnUser(t, ln.Addr().String(), user2)
		defer conn2.Close()

		sMSG, _ := NewMessage(user2, user1, "", TEXT, "hello test")
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