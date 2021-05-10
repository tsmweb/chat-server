package core

import (
	"fmt"
	"time"
)

type memoryRepository struct {
	usersOnline map[string]string
	messagesOffline map[string][]*Message
}

func NewMemoryRepository() Repository {
	return &memoryRepository{
		usersOnline: make(map[string]string),
		messagesOffline: make(map[string][]*Message),
	}
}

func (mr *memoryRepository) generatesMessages(userID string) {
	var msgs []*Message

	for i := 0; i < 10; i++ {
		msg, _ := NewMessage(
			"+5518988888888",
			userID,
			"",
			"Chat",
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
	mr.generatesMessages(userID)
	return mr.messagesOffline[userID], nil
}
