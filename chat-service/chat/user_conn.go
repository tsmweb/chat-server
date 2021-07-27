package chat

import (
	"encoding/json"
	"github.com/tsmweb/chat-service/chat/message"
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
		return nil, u.writeResponse(msg.ID, message.ContentError, err.Error())
	}

	return msg, nil
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
	if err := decoder.Decode(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// WriteMessage writes a message on the user's connection.
func (u *UserConn) WriteMessage(msg *message.Message) error {
	u.io.Lock()
	defer u.io.Unlock()

	return u.writer.Writer(u.conn, msg)
}

// WriteACK writes an acknowledgment message on the user's connection.
func (u *UserConn) WriteACK(msgID string, content string) error {
	return u.writeResponse(msgID, message.ContentACK, content)
}

func (u *UserConn) writeResponse(msgID string, contentType message.ContentType, content string) error {
	u.io.Lock()
	defer u.io.Unlock()

	res := message.NewResponse(msgID, contentType, content)
	return u.writer.Writer(u.conn, res)
}
