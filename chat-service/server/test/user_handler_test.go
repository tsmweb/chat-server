package test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/chat-service/adapter"
	"github.com/tsmweb/chat-service/server"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/chat-service/server/user"
	"testing"
)

func TestHandleUserStatus_Execute(t *testing.T) {
	userID := "+5518977777777"
	imContactID := "+5518955555555"
	myContactID := "+5518966666666"
	ctx := context.Background()
	encode := user.EncoderFunc(adapter.UserMarshal)

	t.Run("when setUserStatus fails with error", func(t *testing.T) {
		p := new(mockProducer)

		r := new(mockRepository)
		r.On("AddUserOnline", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error AddUserOnline")).
			Once()
		r.On("DeleteUserOnline", mock.Anything).
			Return(errors.New("error DeleteUserOnline")).
			Once()

		handler := server.NewHandleUserStatus(encode, p, r)
		err := handler.Execute(ctx, userID, user.Online, nil, nil)
		assert.Equal(t, userID, err.UserID)

		err = handler.Execute(ctx, userID, user.Offline, nil, nil)
		assert.Equal(t, userID, err.UserID)
	})

	t.Run("when publishUserStatus fails with error", func(t *testing.T) {
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()
		p.On("Close").
			Once()

		r := new(mockRepository)
		r.On("AddUserOnline", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		handler := server.NewHandleUserStatus(encode, p, r)
		err := handler.Execute(ctx, userID, user.Online, nil, nil)
		assert.Equal(t, userID, err.UserID)
	})

	t.Run("when notifyContactStatusToUser fails with error", func(t *testing.T) {
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		p.On("Close").
			Once()

		r := new(mockRepository)
		r.On("AddUserOnline", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		r.On("GetUserContactsOnline", mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		handler := server.NewHandleUserStatus(encode, p, r)
		err := handler.Execute(ctx, userID, user.Online, nil, nil)
		assert.Equal(t, userID, err.UserID)
	})

	t.Run("when notifyUserStatusToContacts fails with error", func(t *testing.T) {
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		p.On("Close").
			Once()

		r := new(mockRepository)
		r.On("DeleteUserOnline", mock.Anything).
			Return(nil).
			Once()
		r.On("GetContactsWithUserOnline", mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		handler := server.NewHandleUserStatus(encode, p, r)
		err := handler.Execute(ctx, userID, user.Offline, nil, nil)
		assert.Equal(t, userID, err.UserID)
	})

	t.Run("when user handler succeeds", func(t *testing.T) {
		chMessage := make(chan message.Message, 1)
		chSendMessage := make(chan message.Message, 1)

		r := new(mockRepository)
		r.On("AddUserOnline", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		r.On("GetUserContactsOnline", mock.Anything).
			Return([]string{myContactID}, nil).
			Once()
		r.On("GetContactsWithUserOnline", mock.Anything).
			Return([]string{imContactID}, nil).
			Once()

		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		p.On("Close").
			Once()

		handler := server.NewHandleUserStatus(encode, p, r)
		err := handler.Execute(ctx, userID, user.Online, chMessage, chSendMessage)
		assert.Nil(t, err)

		m := <-chSendMessage
		assert.Equal(t, m.From, myContactID)
		//t.Log(m)

		m = <-chMessage
		assert.Equal(t, m.To, imContactID)
		//t.Log(m)
	})
}
