package core

import (
	"fmt"
	"github.com/tsmweb/chat-service/core/ctype"
	"time"
)

const (
	userTest1 = "+5518911111111"
	userTest2 = "+5518922222222"
	userTest3 = "+5518933333333"
)

type memoryRepository struct {
	usersOnline map[string]string
	messagesOffline map[string][]*Message
	blockedUser map[string]string
}

func NewMemoryRepository() Repository {
	mr := &memoryRepository{
		usersOnline: make(map[string]string),
		messagesOffline: make(map[string][]*Message),
		blockedUser: make(map[string]string),
	}

	mr.blockedUser[userTest3] = userTest1
	return mr
}

func (mr *memoryRepository) generatesMessages(userID string) {
	var msgs []*Message

	for i := 0; i < 3; i++ {
		msg, _ := NewMessage(
			userTest3,
			userID,
			"",
			ctype.TEXT,
			fmt.Sprintf("%d: test internal", i))
		msgs = append(msgs, msg)
	}

	mr.messagesOffline[userID] = msgs
}

func (mr *memoryRepository) AddUserOnline(userID string, host string, createAt time.Time) error {
	mr.usersOnline[userID] = host
	return nil
}

func (mr *memoryRepository) DeleteUserOnline(userID string) error {
	delete(mr.usersOnline, userID)
	return nil
}

func (mr *memoryRepository) GetUserOnline(userID string) (string, bool, error) {
	host, ok := mr.usersOnline[userID]
	return host, ok, nil
}

func (mr *memoryRepository) GetMessagesOffline(userID string) ([]*Message, error) {
	//mr.generatesMessages(userID)
	return mr.messagesOffline[userID], nil
}

func (mr *memoryRepository) IsBlockedUser(userID string, blockedID string) (bool, error) {
	bID := mr.blockedUser[userID]
	return bID == blockedID, nil
}
