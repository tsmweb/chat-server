package chat

import (
	"encoding/json"
	"github.com/tsmweb/chat-service/chat/message"
	"github.com/tsmweb/chat-service/util/connutil"
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

func (u *User) Receive() (*message.Message, error) {
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
		return nil, u.WriteResponse(msg.ID, message.ERROR, err.Error())
	}

	return msg, nil
}

func (u *User) readMessage() (*message.Message, error) {
	u.io.Lock()
	defer u.io.Unlock()

	r, err := u.reader.Reader(u.conn)
	if err != nil {
		return nil, err
	}

	msg := &message.Message{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (u *User) WriteMessage(msg *message.Message) error {
	u.io.Lock()
	defer u.io.Unlock()

	return u.writer.Writer(u.conn, msg)
}

func (u *User) WriteResponse(msgID string, contentType message.ContentType, content string) error {
	u.io.Lock()
	defer u.io.Unlock()

	res := message.NewResponse(msgID, contentType, content)
	return u.writer.Writer(u.conn, res)
}

func (u *User) WriteACK(msgID string, content string) error {
	return u.WriteResponse(msgID, message.ACK, content)
}

// UserStatus type that represents the user's status as UserOnline and UserOffline.
type UserStatus int

const (
	UserOnline  UserStatus = 0x1
	UserOffline            = 0x2
)

func (us UserStatus) String() (str string) {
	name := func(status UserStatus, name string) bool {
		if us&status == 0 {
			return false
		}
		str = name
		return true
	}

	if name(UserOnline, "UserOnline") { return }
	if name(UserOffline, "UserOffline") { return }

	return
}

type UserPresence struct {
	ID     string
	Status UserStatus
}
