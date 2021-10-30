package broker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
	"testing"
	"time"
)

func TestUserHandler_Execute(t *testing.T) {
	ctx := context.Background()
	chMessage := make(chan message.Message, 1)

	go func(chMsg <-chan message.Message) {
		for msg := range chMsg {
			t.Log(msg.String())
		}
	}(chMessage)

	usr := user.New("+5518977777777", user.Online, "H01")
	msg1, _ := message.New("+5518911111111", "+5518977777777", "", message.ContentTypeText, "message test 1")
	msg2, _ := message.New("+5518922222222", "+5518977777777", "", message.ContentTypeText, "message test 2")
	msg3, _ := message.New("+5518933333333", "+5518977777777", "", message.ContentTypeACK, "read")
	msgs := []*message.Message {msg1, msg2, msg3}

	userRepo := new(mockUserRepository)
	userRepo.On("AddUserPresence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Once()
	userRepo.On("GetAllContactsOnline", mock.Anything, mock.Anything).
		Return([]string{"+5518988888888", "+5518999999999"}, nil).
		Once()
	userRepo.On("GetAllRelationshipsOnline", mock.Anything, mock.Anything).
		Return([]string{"+5518944444444", "+5518955555555"}, nil).
		Once()

	msgRepo := new(mockMessageRepository)
	msgRepo.On("GetAllMessages", mock.Anything, mock.Anything).
		Return(msgs, nil).
		Once()
	msgRepo.On("DeleteAllMessages", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	handler := NewUserHandler(userRepo, msgRepo)
	err := handler.Execute(ctx, *usr, chMessage)
	assert.Nil(t, err)

	time.Sleep(time.Millisecond * 100)
}
