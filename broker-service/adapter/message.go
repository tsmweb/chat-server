package adapter

import (
	"time"

	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/infra/protobuf"
	"google.golang.org/protobuf/proto"
)

// MessageMarshal is a message.Message encoder for protobuf.Message.
func MessageMarshal(m *message.Message) ([]byte, error) {
	mpb := protobufFromMessage(m)
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

func protobufFromMessage(m *message.Message) *protobuf.Message {
	return &protobuf.Message{
		Id:          m.ID,
		From:        m.From,
		To:          m.To,
		Group:       m.Group,
		Date:        m.Date.Unix(),
		ContentType: protobuf.ContentType(protobuf.ContentType_value[m.ContentType]),
		Content:     m.Content,
	}
}

func protobufToMessage(mpb *protobuf.Message, m *message.Message) {
	m.ID = mpb.GetId()
	m.From = mpb.GetFrom()
	m.To = mpb.GetTo()
	m.Group = mpb.GetGroup()
	m.Date = time.Unix(mpb.GetDate(), 0)
	m.ContentType = mpb.GetContentType().String()
	m.Content = mpb.GetContent()
}
