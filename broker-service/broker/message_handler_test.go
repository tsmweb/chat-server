package broker

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/broker-service/broker/message"
	"testing"
)

func TestMessageHandler_Execute(t *testing.T) {
	ctx := context.Background()

	msg, _ := message.New("+5518911111111", "+5518977777777", "",
		message.ContentTypeText, "message test")
	msgGroup, _ := message.New("+5518911111111", "", "123456",
		message.ContentTypeText, "message group test")

	encoder := new(mockMessageEncoder)
	encoder.On("Marshal", mock.Anything).
		Return([]byte{}, nil)

	t.Run("when handling group messages fails", func(t *testing.T) {
		userRepo := new(mockUserRepository)
		msgRepo := new(mockMessageRepository)
		producer := new(mockProducer)
		producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		queue := new(mockKafka)
		queue.On("NewProducer", mock.Anything).
			Return(producer)
		handler := NewMessageHandler(userRepo, msgRepo, queue, encoder)

		msgRepo.On("GetAllGroupMembers", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		err := handler.Execute(ctx, *msgGroup)
		assert.NotNil(t, err)

		userRepo.On("GetUserServer", mock.Anything, mock.Anything).
			Return("", errors.New("error"))
		msgRepo.On("GetAllGroupMembers", mock.Anything, mock.Anything).
			Return([]string{"+5518911111111", "+5518977777777", "+5518988888888"}, nil)

		err = handler.Execute(ctx, *msgGroup)
		assert.NotNil(t, err)

		userRepo.On("GetUserServer", mock.Anything, mock.Anything).
			Return("H01", nil)
		producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error"))

		err = handler.Execute(ctx, *msgGroup)
		assert.NotNil(t, err)
	})

	t.Run("when handling group messages is successful", func(t *testing.T) {
		userRepo := new(mockUserRepository)
		userRepo.On("GetUserServer", mock.Anything, mock.Anything).
			Return("H01", nil)

		msgRepo := new(mockMessageRepository)
		msgRepo.On("GetAllGroupMembers", mock.Anything, mock.Anything).
			Return([]string{"+5518911111111", "+5518977777777", "+5518988888888"}, nil).
			Once()

		producer := new(mockProducer)
		producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		queue := new(mockKafka)
		queue.On("NewProducer", mock.Anything).
			Return(producer)

		handler := NewMessageHandler(userRepo, msgRepo, queue, encoder)
		err := handler.Execute(ctx, *msgGroup)
		assert.Nil(t, err)
	})

	t.Run("when message handling fails", func(t *testing.T) {
		msgRepo := new(mockMessageRepository)
		userRepo := new(mockUserRepository)
		producer := new(mockProducer)
		queue := new(mockKafka)
		queue.On("NewProducer", mock.Anything).
			Return(producer)
		handler := NewMessageHandler(userRepo, msgRepo, queue, encoder)

		userRepo.On("IsValidUser", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		err := handler.Execute(ctx, *msg)
		assert.NotNil(t, err)

		userRepo.On("IsValidUser", mock.Anything, mock.Anything).
			Return(true, nil)
		userRepo.On("IsBlockedUser", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		err = handler.Execute(ctx, *msg)
		assert.NotNil(t, err)

		userRepo.On("IsBlockedUser", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil)
		userRepo.On("GetUserServer", mock.Anything, mock.Anything).
			Return("", errors.New("error")).
			Once()

		err = handler.Execute(ctx, *msg)
		assert.NotNil(t, err)

		userRepo.On("IsBlockedUser", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil)
		userRepo.On("GetUserServer", mock.Anything, mock.Anything).
			Return("H01", nil)
		producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		err = handler.Execute(ctx, *msg)
		assert.NotNil(t, err)
	})

	t.Run("when message handling is successful", func(t *testing.T) {
		msgRepo := new(mockMessageRepository)
		userRepo := new(mockUserRepository)
		producer := new(mockProducer)
		producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		queue := new(mockKafka)
		queue.On("NewProducer", mock.Anything).
			Return(producer)
		handler := NewMessageHandler(userRepo, msgRepo, queue, encoder)

		userRepo.On("IsValidUser", mock.Anything, mock.Anything).
			Return(true, nil)
		userRepo.On("IsBlockedUser", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil)
		userRepo.On("GetUserServer", mock.Anything, mock.Anything).
			Return("H01", nil).
			Once()

		err := handler.Execute(ctx, *msg)
		assert.Nil(t, err)

		userRepo.On("GetUserServer", mock.Anything, mock.Anything).
			Return("", nil).
			Once()

		err = handler.Execute(ctx, *msg)
		assert.Nil(t, err)
	})
}
