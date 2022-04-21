package server

import (
	"encoding/json"
	"github.com/tsmweb/chat-service/server/message"
	"net"
	"sync"
)

// UserConn type that represents the user connection.
type UserConn struct {
	userID string

	io   sync.Mutex
	conn net.Conn

	reader ConnReader
	writer ConnWriter
}

// Receive read user connection data.
func (u *UserConn) Receive() (*message.Message, error) {
	msg, err := u.readMessage()
	if err != nil {
		u.conn.Close()
		return nil, err
	}
	// Handled some control internal.
	if msg == nil {
		return nil, nil
	}

	msg.From = u.userID

	if err = msg.Validate(); err != nil {
		return nil, u.WriteResponse(msg.ID, message.ContentTypeError, err.Error())
	}

	return msg, nil
}

// WriteMessage writes a message on the user's connection.
func (u *UserConn) WriteMessage(msg *message.Message) error {
	u.io.Lock()
	defer u.io.Unlock()

	return u.writer.Writer(u.conn, msg)
}

// WriteResponse write a response message on the user's connection.
func (u *UserConn) WriteResponse(msgID string, contentType message.ContentType, content string) error {
	u.io.Lock()
	defer u.io.Unlock()

	res := message.NewResponse(msgID, contentType, content)
	return u.writer.Writer(u.conn, res)
}

func (u *UserConn) readMessage() (*message.Message, error) {
	u.io.Lock()
	defer u.io.Unlock()

	r, err := u.reader.Reader(u.conn)
	if err != nil {
		return nil, err
	}

	msg := new(message.Message)
	decoder := json.NewDecoder(r)
	if err = decoder.Decode(msg); err != nil {
		return nil, err
	}

	return msg, nil
}
