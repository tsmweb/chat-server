package message

import (
	"github.com/tsmweb/chat-service/pkg/protobuf"
	"google.golang.org/protobuf/proto"
	"time"
)

// Marshal marshal to protobuf.Message.
func Marshal(msg *Message) ([]byte, error) {
	mpb := new(protobuf.Message)
	protobufFromMessage(msg, mpb)
	return proto.Marshal(mpb)
}

// Unmarshal unmarshal to Message.
func Unmarshal(in []byte, msg *Message) error {
	mpb := new(protobuf.Message)
	if err := proto.Unmarshal(in, mpb); err != nil {
		return err
	}
	protobufToMessage(mpb, msg)
	return nil
}

func protobufFromMessage(m *Message, mpb *protobuf.Message) {
	mpb.Id = m.ID
	mpb.From = m.From
	mpb.To = m.To
	mpb.Date = m.Date.Unix()
	mpb.ContentType = protobuf.ContentType(protobuf.ContentType_value[m.ContentType])
	mpb.Content = m.Content
}

func protobufToMessage(mpb *protobuf.Message, m *Message) {
	m.ID = mpb.GetId()
	m.From = mpb.GetFrom()
	m.To = mpb.GetTo()
	m.Date = time.Unix(mpb.GetDate(), 0)
	m.ContentType = mpb.GetContentType().String()
	m.Content = mpb.GetContent()
}
