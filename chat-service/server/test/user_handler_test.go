package test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/chat-service/adapter"
	"github.com/tsmweb/chat-service/server"
	"github.com/tsmweb/chat-service/server/user"
	"testing"
)

func TestHandleUserStatus_Execute(t *testing.T) {
	userID := "+5518977777777"
	ctx := context.Background()
	encode := user.EncoderFunc(adapter.UserMarshal)

	t.Run("when publishUserStatus fails with error", func(t *testing.T) {
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		p.On("Close").
			Once()

		handler := server.NewHandleUserStatus(encode, p)
		err := handler.Execute(ctx, userID, user.Online)
		assert.Equal(t, userID, err.UserID)
	})

	t.Run("when user handler succeeds", func(t *testing.T) {
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		p.On("Close").
			Once()

		handler := server.NewHandleUserStatus(encode, p)
		err := handler.Execute(ctx, userID, user.Online)
		assert.Nil(t, err)
	})
}
