package core

import (
	"encoding/json"
	"github.com/tsmweb/chat-service/common/connutil"
	"github.com/tsmweb/chat-service/core/ctype"
	"io"
	"sync"
)

type User struct {
	id string

	io   sync.Mutex
	conn io.ReadWriteCloser

	reader connutil.Reader
	writer connutil.Writer
}

func (u *User) Receive() (*Message, error) {
	msg, err := u.readMessage()
	if err != nil {
		u.conn.Close()
		return nil, err
	}
	// Handled some control internal.
	if msg == nil {
		return nil, nil
	}

	msg.From = u.id

	if err = msg.Validate(); err != nil {
		return nil, u.WriteResponse(msg.ID, ctype.ERROR, err.Error())
	}

	return msg, nil
}

func (u *User) readMessage() (*Message, error) {
	u.io.Lock()
	defer u.io.Unlock()

	r, err := u.reader.Reader(u.conn)
	if err != nil {
		return nil, err
	}

	msg := &Message{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (u *User) WriteMessage(msg *Message) error {
	u.io.Lock()
	defer u.io.Unlock()

	return u.writer.Writer(u.conn, msg)
}

func (u *User) WriteResponse(msgID string, contentType ctype.ContentType, content string) error {
	u.io.Lock()
	defer u.io.Unlock()

	res := NewResponse(msgID, contentType, content)
	return u.writer.Writer(u.conn, res)
}

func (u *User) WriteACK(msgID string, content string) error {
	return u.WriteResponse(msgID, ctype.ACK, content)
}
