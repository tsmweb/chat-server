package core

import (
	"context"
	"fmt"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"net"
)

type Chat struct {
	poller    netpoll.Poller
	executor  *executor.Executor
	localhost string

	chUserIN  chan *User
	chUserOUT chan *User
	chMessage chan *Message

	userStatusHandler *UserStatusHandler
	messageHandler    *MessageHandler
	errorDispatcher   *ErrorDispatcher
}

func NewChat(
	p netpoll.Poller,
	e *executor.Executor,
	lh string,
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
		userStatusHandler: ush,
		messageHandler:    mh,
		errorDispatcher:   ed,
	}

	go chat.start()

	return chat
}

func (c *Chat) Register(userID string, conn net.Conn) error {
	user := &User{
		id:   userID,
		conn: conn,
	}

	// Create netpoll event descriptor for conn.
	desc := netpoll.Must(netpoll.HandleRead(conn))

	// Subscribe to events about conn.
	err := c.poller.Start(desc, func(ev netpoll.Event) {
		if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
			c.poller.Stop(desc)
			c.chUserOUT <- user
			return
		}

		c.executor.Schedule(func(ctx context.Context) {
			msg, err := user.Receive()
			if err != nil {
				c.poller.Stop(desc)
				c.chUserOUT <- user
				c.sendError(fmt.Errorf("%s user.Receive(): %s", user.id, err.Error()))
			}
			if msg != nil {
				c.SendMessage(msg)
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

func (c *Chat) messageJob(msg *Message, user *User) executor.Job {
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

func (c *Chat) userINJob(user *User) executor.Job {
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

func (c *Chat) userOUTJob(user *User) executor.Job {
	return func(ctx context.Context) {
		err := c.userStatusHandler.HandleStatus(user.id, c.localhost, OFFLINE)
		if err != nil {
			c.sendError(fmt.Errorf("%s SetStatus(): %s", user.id, err.Error()))
			return
		}
	}
}
