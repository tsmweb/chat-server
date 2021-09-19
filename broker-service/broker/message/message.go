package message

import (
	"context"
	"encoding/json"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"strconv"
	"strings"
	"time"
)

const (
	InvalidMessage = "invalid message"
	AckMessage     = "ack"
)

// ContentType represents the type of message content,
// such as ContentTypeACK, ContentTypeText, ContentTypeMedia, ContentTypeStatus,
// ContentTypeInfo and ContentTypeError.
type ContentType int

const (
	ContentTypeACK    ContentType = 0x1
	ContentTypeText   ContentType = 0x2
	ContentTypeMedia  ContentType = 0x4
	ContentTypeStatus ContentType = 0x8
	ContentTypeInfo   ContentType = 0x10
	ContentTypeError  ContentType = 0x20
)

func (ct ContentType) String() (str string) {
	name := func(contentType ContentType, name string) bool {
		if ct&contentType == 0 {
			return false
		}
		str = name
		return true
	}

	if name(ContentTypeACK, "ack") {
		return
	}
	if name(ContentTypeText, "text") {
		return
	}
	if name(ContentTypeMedia, "media") {
		return
	}
	if name(ContentTypeStatus, "status") {
		return
	}
	if name(ContentTypeInfo, "info") {
		return
	}
	if name(ContentTypeError, "error") {
		return
	}

	return
}

// Repository represents an abstraction of the data persistence layer.
type Repository interface {
	GetAllGroupMembers(ctx context.Context, groupID string) ([]string, error)
	GetAllMessages(ctx context.Context, userID string) ([]*Message, error)
	DeleteAllMessages(ctx context.Context, userID string) error
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

// New creates and returns a new Message instance.
func New(from string, to string, group string, contentType ContentType, content string) (*Message, error) {
	return newMessage(from, to, group, time.Now().UTC(), contentType.String(), content)
}

func newMessage(from string, to string, group string, date time.Time, contentType string, content string) (*Message, error) {
	msg := &Message{
		From:        from,
		To:          to,
		Group:       group,
		Date:        date,
		ContentType: contentType,
		Content:     content,
	}

	if err := msg.Validate(); err != nil {
		return nil, err
	}

	if err := msg.generateID(); err != nil {
		return nil, err
	}

	return msg, nil
}

// ReplicateTo replicate the message to another recipient.
func (m *Message) ReplicateTo(to string) (*Message, error) {
	return newMessage(m.From, to, m.Group, m.Date, m.ContentType, m.Content)
}

func (m *Message) generateID() error {
	id, err := hashutil.HashSHA1(m.From + m.To + m.Group + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

// Validate verifies that the required attributes of the message are present.
func (m *Message) Validate() error {
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

// Encoder is a Message encoder for byte slice.
type Encoder interface {
	Marshal(m *Message) ([]byte, error)
}

// The EncoderFunc type is an adapter to allow the use of ordinary functions as encoders of Message for byte slice.
// If f is a function with the appropriate signature, EncoderFunc(f) is a Encoder that calls f.
type EncoderFunc func(m *Message) ([]byte, error)

// Marshal calls f(m).
func (f EncoderFunc) Marshal(m *Message) ([]byte, error) {
	return f(m)
}

// Decoder is a byte slice decoder for Message.
type Decoder interface {
	Unmarshal(in []byte, m *Message) error
}

// The DecoderFunc type is an adapter to allow the use of ordinary functions as decoders of byte slice for Message.
// If f is a function with the appropriate signature, DecoderFunc(f) is a Decoder that calls f.
type DecoderFunc func(in []byte, m *Message) error

// Unmarshal calls f(in, m).
func (f DecoderFunc) Unmarshal(in []byte, m *Message) error {
	return f(in, m)
}