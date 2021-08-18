package test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/chat-service/adapter"
	"github.com/tsmweb/chat-service/server"
	"github.com/tsmweb/chat-service/server/message"
	"testing"
)

func TestHandleMessage_Execute(t *testing.T) {
	ctx := context.Background()
	encode := message.EncoderFunc(adapter.MessageMarshal)
	msg, _ := message.NewMessage("+5518977777777", "+5518966666666", "", message.ContentTypeText, "hello")

	t.Run("when message handler fails with error", func(t *testing.T) {
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		p.On("Close").
			Once()

		handler := server.NewHandleMessage(encode, p)
		err := handler.Execute(ctx, *msg)
		assert.Equal(t, msg.From, err.UserID)
	})

	t.Run("when message handler succeeds", func(t *testing.T) {
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		p.On("Close").
			Once()

		handler := server.NewHandleMessage(encode, p)
		err := handler.Execute(ctx, *msg)
		assert.Nil(t, err)
	})
}

func TestHandleGroupMessage_Execute(t *testing.T) {
	msg, _ := message.NewMessage("+5518977777777", "", "123456", message.ContentTypeText, "hello")

	t.Run("when group message handler fails with error", func(t *testing.T) {
		r := new(mockRepository)
		r.On("GetGroupMembers", mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		handler := server.NewHandleGroupMessage(r)
		err := handler.Execute(*msg, nil)
		assert.Equal(t, err.UserID, msg.From)
	})

	t.Run("when group message handler succeeds", func(t *testing.T) {
		memberID := "+5518955555555"

		r := new(mockRepository)
		r.On("GetGroupMembers", mock.Anything).
			Return([]string{memberID}, nil).
			Once()

		chMessage := make(chan message.Message, 1)

		handler := server.NewHandleGroupMessage(r)
		err := handler.Execute(*msg, chMessage)
		assert.Nil(t, err)

		m := <-chMessage
		assert.Equal(t, m.To, memberID)
	})
}
