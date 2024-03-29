package test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/chat-service/adapter"
	"github.com/tsmweb/chat-service/server"
	"github.com/tsmweb/chat-service/server/user"
)

func TestHandleUserStatus_Execute(t *testing.T) {
	userID := "+5518977777777"
	ctx := context.Background()
	encode := user.EncoderFunc(adapter.UserMarshal)

	t.Run("when publishUserStatus fails with error", func(t *testing.T) {
		userProducer := new(mockProducer)
		userProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		userProducer.On("Close").
			Once()

		userPresenceProducer := new(mockProducer)
		userPresenceProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		userPresenceProducer.On("Close").
			Once()

		handler := server.NewHandleUserStatus(encode, userProducer, userPresenceProducer)
		err := handler.Execute(ctx, userID, user.Online)
		assert.NotNil(t, err)
	})

	t.Run("when user handler succeeds", func(t *testing.T) {
		userProducer := new(mockProducer)
		userProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		userProducer.On("Close").
			Once()

		userPresenceProducer := new(mockProducer)
		userPresenceProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		userPresenceProducer.On("Close").
			Once()

		handler := server.NewHandleUserStatus(encode, userProducer, userPresenceProducer)
		err := handler.Execute(ctx, userID, user.Online)
		assert.Nil(t, err)
	})
}
