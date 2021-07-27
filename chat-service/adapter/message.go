package adapter

import (
	"github.com/tsmweb/chat-service/chat/message"
	"github.com/tsmweb/chat-service/infra/protobuf"
	"google.golang.org/protobuf/proto"
	"time"
)

// MessageMarshal is a message.Message encoder for protobuf.Message.
func MessageMarshal(m *message.Message) ([]byte, error) {
	mpb := new(protobuf.Message)
	protobufFromMessage(m, mpb)
	return proto.Marshal(mpb)
}

// MessageUnmarshal is a protobuf.Message decoder for message.Message.
func MessageUnmarshal(in []byte, m *message.Message) error {
	mpb := new(protobuf.Message)
	if err := proto.Unmarshal(in, mpb); err != nil {
		return err
	}
	protobufToMessage(mpb, m)
	return nil
}

func protobufFromMessage(m *message.Message, mpb *protobuf.Message) {
	mpb.Id = m.ID
	mpb.From = m.From
	mpb.To = m.To
	mpb.Date = m.Date.Unix()
	mpb.ContentType = protobuf.ContentType(protobuf.ContentType_value[m.ContentType])
	mpb.Content = m.Content
}

func protobufToMessage(mpb *protobuf.Message, m *message.Message) {
	m.ID = mpb.GetId()
	m.From = mpb.GetFrom()
	m.To = mpb.GetTo()
	m.Date = time.Unix(mpb.GetDate(), 0)
	m.ContentType = mpb.GetContentType().String()
	m.Content = mpb.GetContent()
}