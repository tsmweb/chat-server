package core

import (
	"encoding/json"
	"github.com/tsmweb/chat-service/common/connutil"
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

	if err = msg.Validate(); err != nil {
		return nil, u.WriteResponse(msg.ID, ERROR, err.Error())
	}

	// Spoofed internal is discarded.
	if msg.From != u.id {
		return nil, u.WriteResponse(msg.ID, ERROR, ErrMessageSpoof.Error())
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

func (u *User) WriteResponse(msgID string, contentType ContentType, content string) error {
	u.io.Lock()
	defer u.io.Unlock()

	res := NewResponse(msgID, contentType, content)
	return u.writer.Writer(u.conn, res)
}

func (u *User) WriteACK(msgID string) error {
	return u.WriteResponse(msgID, ACK, "sent")
}
