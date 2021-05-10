package core

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"sync"
)

type User struct {
	io sync.Mutex
	conn io.ReadWriteCloser

	id   string
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
		return nil, u.WriteError(msg.ID, err)
	}

	// Spoofed internal is discarded.
	if msg.From != u.id {
		return nil, u.WriteError(msg.ID, ErrMessageSpoof)
	}

	return msg, nil
}

func (u *User) readMessage() (*Message, error) {
	u.io.Lock()
	defer u.io.Unlock()

	h, r, err := wsutil.NextReader(u.conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}
	if h.OpCode.IsControl() {
		return nil, wsutil.ControlFrameHandler(u.conn, ws.StateServerSide)(h, r)
	}

	msg := &Message{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (u *User) WriteMessage(msg *Message) error {
	return u.write(msg)
}

func (u *User) WriteError(msgID string, err error) error {
	return u.write(Error{
		ID: msgID,
		Error: err.Error(),
	})
}

func (u *User) write(x interface{}) error {
	w := wsutil.NewWriter(u.conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	u.io.Lock()
	defer u.io.Unlock()

	if err := encoder.Encode(x); err != nil {
		return err
	}

	return w.Flush()
}
