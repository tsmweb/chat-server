package core

import (
	"encoding/json"
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	ID          string    `json:"id"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Group       string    `json:"group,omitempty"`
	Date        time.Time `json:"date"`
	ContentType string    `json:"content_type"`
	Content     string    `json:"content"`
}

func NewMessage(from, to, group, contentType, content string) (*Message, error) {
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
		ContentType: contentType,
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

type Error struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}
