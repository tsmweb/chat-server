package test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/chat-service/adapter"
	"github.com/tsmweb/chat-service/server"
	"github.com/tsmweb/chat-service/server/message"
)

func TestHandleMessage_Execute(t *testing.T) {
	ctx := context.Background()
	encode := message.EncoderFunc(adapter.MessageMarshal)
	msg, _ := message.NewMessage("+5518977777777", "+5518966666666", "",
		message.ContentTypeText, "hello")

	t.Run("when message handler fails with error", func(t *testing.T) {
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		p.On("Close").
			Once()

		handler := server.NewHandleMessage(encode, p)
		err := handler.Execute(ctx, msg)
		assert.NotNil(t, err)
	})

	t.Run("when message handler succeeds", func(t *testing.T) {
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		p.On("Close").
			Once()

		handler := server.NewHandleMessage(encode, p)
		err := handler.Execute(ctx, msg)
		assert.Nil(t, err)
	})
}
