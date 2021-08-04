package message

import (
	"encoding/json"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"strconv"
	"strings"
	"time"
)

const (
	InvalidMessage = "invalid message"
)

// ContentType represents the type of message content,
// such as ContentACK, ContentText, ContentMedia, ContentStatus and ContentError.
type ContentType int

const (
	ContentACK    ContentType = 0x1
	ContentText               = 0x2
	ContentMedia              = 0x4
	ContentStatus             = 0x8
	ContentError              = 0x80
)

func (ct ContentType) String() (str string) {
	name := func(contentType ContentType, name string) bool {
		if ct&contentType == 0 {
			return false
		}
		str = name
		return true
	}

	if name(ContentACK, "ack") {
		return
	}
	if name(ContentText, "text") {
		return
	}
	if name(ContentMedia, "media") {
		return
	}
	if name(ContentStatus, "status") {
		return
	}
	if name(ContentError, "error") {
		return
	}

	return
}

var (
	ErrIDValidateModel          = &cerror.ErrValidateModel{Msg: "required id"}
	ErrFromValidateModel        = &cerror.ErrValidateModel{Msg: "required from"}
	ErrReceiverValidateModel    = &cerror.ErrValidateModel{Msg: "required to or group"}
	ErrDateValidateModel        = &cerror.ErrValidateModel{Msg: "required date"}
	ErrContentTypeValidateModel = &cerror.ErrValidateModel{Msg: "required content_type"}
	ErrContentValidateModel     = &cerror.ErrValidateModel{Msg: "required content"}
)

// Message represents data sent and received by users.
type Message struct {
	ID          string    `json:"id"`
	From        string    `json:"from,omitempty"`
	To          string    `json:"to,omitempty"`
	Group       string    `json:"group,omitempty"`
	Date        time.Time `json:"date"`
	ContentType string    `json:"content_type"`
	Content     string    `json:"content"`
}

// NewResponse creates and returns a new Message instance.
func NewResponse(msgID string, contentType ContentType, content string) *Message {
	return &Message{
		ID:          msgID,
		Date:        time.Now().UTC(),
		ContentType: contentType.String(),
		Content:     content,
	}
}

// NewMessage creates and returns a new Message instance.
func NewMessage(from string, to string, group string, contentType ContentType, content string) (*Message, error) {
	msgID, err := hashutil.HashSHA1(from + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return nil, err
	}

	msg := &Message{
		ID:          msgID,
		From:        from,
		To:          to,
		Group:       group,
		Date:        time.Now().UTC(),
		ContentType: contentType.String(),
		Content:     content,
	}

	if err = msg.Validate(); err != nil {
		return nil, err
	}

	return msg, nil
}

// Validate verifies that the required attributes of the message are present.
func (m *Message) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return ErrIDValidateModel
	}
	if strings.TrimSpace(m.From) == "" {
		return ErrFromValidateModel
	}
	if strings.TrimSpace(m.To) == "" && strings.TrimSpace(m.Group) == "" {
		return ErrReceiverValidateModel
	}
	if m.Date.IsZero() {
		return ErrDateValidateModel
	}
	if strings.TrimSpace(m.ContentType) == "" {
		return ErrContentTypeValidateModel
	}
	if strings.TrimSpace(m.Content) == "" {
		return ErrContentValidateModel
	}
	return nil
}

// IsGroupMessage returns true if the message is addressed to a group of users.
func (m *Message) IsGroupMessage() bool {
	return strings.TrimSpace(m.Group) != ""
}

func (m *Message) String() string {
	mj, _ := json.Marshal(m)
	return string(mj)
}
