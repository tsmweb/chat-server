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
	BlockedMessage = "message blocked by %s"
)

// ContentType represents the type of message content, such as ACK, TEXT, MEDIA, STATUS and ERROR.
type ContentType int

const (
	ACK    ContentType = 0x1
	TEXT               = 0x2
	MEDIA              = 0x4
	STATUS             = 0x8
	ERROR              = 0x80
)

func (ct ContentType) String() (str string) {
	name := func(contentType ContentType, name string) bool {
		if ct&contentType == 0 {
			return false
		}
		str = name
		return true
	}

	if name(ACK, "ACK") {
		return
	}
	if name(TEXT, "TEXT") {
		return
	}
	if name(MEDIA, "MEDIA") {
		return
	}
	if name(STATUS, "STATUS") {
		return
	}
	if name(ERROR, "ERROR") {
		return
	}

	return
}

var (
	ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required id"}
	ErrFromValidateModel = &cerror.ErrValidateModel{Msg: "required from"}
	ErrReceiverValidateModel = &cerror.ErrValidateModel{Msg: "required to or group"}
	ErrDateValidateModel = &cerror.ErrValidateModel{Msg: "required date"}
	ErrContentTypeValidateModel = &cerror.ErrValidateModel{Msg: "required content_type"}
	ErrContentValidateModel = &cerror.ErrValidateModel{Msg: "required content"}
)


type Message struct {
	ID          string    `json:"id"`
	From        string    `json:"from,omitempty"`
	To          string    `json:"to,omitempty"`
	Group       string    `json:"group,omitempty"`
	Date        time.Time `json:"date"`
	ContentType string    `json:"content-type"`
	Content     string    `json:"content"`
	Host        string    `json:"host,omitempty"`
}

func NewResponse(msgID string, contentType ContentType, content string) *Message {
	return &Message{
		ID:          msgID,
		Date:        time.Now().UTC(),
		ContentType: contentType.String(),
		Content:     content,
	}
}

func New(from string, to string, group string, contentType ContentType, content string) (*Message, error) {
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

func (m Message) Validate() error {
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

func (m Message) IsGroupMessage() bool {
	return strings.TrimSpace(m.Group) != ""
}

func (m Message) ToJSON() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
}

func (m Message) String() string {
	return string(m.ToJSON())
}
