package core

import (
	"encoding/json"
	"github.com/tsmweb/chat-service/core/ctype"
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"strconv"
	"strings"
	"time"
)

const (
	BlockedMessage = "message blocked by %s"
)

type Message struct {
	ID          string    `json:"id"`
	From        string    `json:"from,omitempty"`
	To          string    `json:"to,omitempty"`
	Group       string    `json:"group,omitempty"`
	Date        time.Time `json:"date"`
	ContentType string    `json:"content-type"`
	Content     string    `json:"content"`
}

func NewResponse(msgID string, contentType ctype.ContentType, content string) *Message {
	return &Message{
		ID: msgID,
		Date: time.Now().UTC(),
		ContentType: contentType.String(),
		Content: content,
	}
}

func NewMessage(from string, to string, group string, contentType ctype.ContentType, content string) (*Message, error) {
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

func (m Message) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(b)
}
