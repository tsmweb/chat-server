package core

import (
	"context"
	"fmt"
	"github.com/tsmweb/chat-service/common/concurrent"
	"github.com/tsmweb/chat-service/common/connutil"
	"github.com/tsmweb/chat-service/common/epoll"
	"net"
)

type Chat struct {
	poller    epoll.EPoll
	executor  concurrent.ExecutorService
	localhost string

	chUserIN  chan *User
	chUserOUT chan *User
	chMessage chan *Message

	reader connutil.Reader
	writer connutil.Writer

	userStatusHandler *UserStatusHandler
	messageHandler    *MessageHandler
	errorDispatcher   *ErrorDispatcher
}

func NewChat(
	p epoll.EPoll,
	e concurrent.ExecutorService,
	lh string,
	reader connutil.Reader,
	writer connutil.Writer,
	ush *UserStatusHandler,
	mh *MessageHandler,
	ed *ErrorDispatcher,
) *Chat {
	chat := &Chat{
		poller:            p,
		executor:          e,
		localhost:         lh,
		chUserIN:          make(chan *User),
		chUserOUT:         make(chan *User),
		chMessage:         make(chan *Message),
		reader:            reader,
		writer:            writer,
		userStatusHandler: ush,
		messageHandler:    mh,
		errorDispatcher:   ed,
	}

	go chat.start()

	return chat
}

func (c *Chat) Register(userID string, conn net.Conn) error {
	user := &User{
		id:     userID,
		conn:   conn,
		reader: c.reader,
		writer: c.writer,
	}

	observer, err := c.poller.ObservableRead(conn)
	if err != nil {
		c.sendError(fmt.Errorf("%s Chat.Register(): %s", user.id, err.Error()))
		return err
	}

	err = observer.Start(func(closed bool, err error) {
		if closed || err != nil {
			observer.Stop()
			c.chUserOUT <- user
			if err != nil {
				c.sendError(fmt.Errorf("%s epoll.Event(): %s", user.id, err.Error()))
			}
			return
		}

		c.executor.Schedule(func(ctx context.Context) {
			msg, err := user.Receive()
			if err != nil {
				observer.Stop()
				c.chUserOUT <- user
				c.sendError(fmt.Errorf("%s user.Receive(): %s", user.id, err.Error()))
			}
			if msg != nil {
				c.SendMessage(msg)

				if err := user.WriteACK(msg.ID); err != nil {
					c.sendError(fmt.Errorf("%s user.WriteACK(): %s", user.id, err.Error()))
				}
			}
		})
	})

	if err != nil {
		c.sendError(fmt.Errorf("%s Chat.Register(): %s", user.id, err.Error()))
		return err
	}

	c.chUserIN <- user
	return nil
}

func (c *Chat) start() {
	users := make(map[string]*User) // all connected users

	for {
		select {
		case msg := <-c.chMessage:
			user := users[msg.To]
			c.executor.Schedule(c.messageJob(msg, user))

		case user := <-c.chUserIN:
			users[user.id] = user
			c.executor.Schedule(c.userINJob(user))

		case user := <-c.chUserOUT:
			delete(users, user.id)
			c.executor.Schedule(c.userOUTJob(user))
		}
	}
}

// SendMessage implements the MessageDispatcher interface.
func (c *Chat) SendMessage(msg *Message) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrClosedChannel
		}
	}()

	c.chMessage <- msg
	return nil
}

func (c *Chat) sendError(err error) {
	c.errorDispatcher.Send(err)
}

func (c *Chat) messageJob(msg *Message, user *User) func(ctx context.Context) {
	return func(ctx context.Context) {
		handle := func(msg *Message) {
			if err := c.messageHandler.HandleMessage(msg); err != nil {
				c.sendError(err)
			}
		}

		if user != nil {
			if err := user.WriteMessage(msg); err != nil {
				handle(msg)
			}
		} else {
			handle(msg)
		}
	}
}

func (c *Chat) userINJob(user *User) func(ctx context.Context) {
	return func(ctx context.Context) {
		err := c.userStatusHandler.HandleStatus(user.id, c.localhost, ONLINE)
		if err != nil {
			c.sendError(fmt.Errorf("%s HandleStatus(): %s", user.id, err.Error()))
			return
		}

		if err = c.messageHandler.SendMessageOffline(user.id, c.chMessage); err != nil {
			c.sendError(fmt.Errorf("%s SendMessageOffline(): %s", user.id, err.Error()))
			return
		}
	}
}

func (c *Chat) userOUTJob(user *User) func(ctx context.Context) {
	return func(ctx context.Context) {
		err := c.userStatusHandler.HandleStatus(user.id, c.localhost, OFFLINE)
		if err != nil {
			c.sendError(fmt.Errorf("%s SetStatus(): %s", user.id, err.Error()))
			return
		}
	}
}
