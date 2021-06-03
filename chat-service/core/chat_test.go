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
	userID := "+5518977777777"

	executor := executor.New(100)
	defer executor.Shutdown()

	chat := initChat(t, executor)

	ln := runServer(t, userID, chat)
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	time.Sleep(10 * time.Millisecond)

	sMSG, _ := NewMessage("+5518966666666", userID, "", "chat", "hello")
	if err = chat.SendMessage(sMSG); err != nil {
		t.Fatalf("error chat.SendMessage(): %v", err)
	}

	rMSG := readMessageFromConn(t, conn)

	t.Log(rMSG)
	assert.Equal(t, sMSG.ID, rMSG.ID)
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

func runServer(tb testing.TB, userID string, chat *Chat) net.Listener {
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

			if err := chat.Register(userID, conn); err != nil {
				tb.Fatalf("failed to register user on chat server: %v", err)
			}
		}
	}()
	return ln
}
