package chat

import (
	"fmt"
	"github.com/tsmweb/chat-service/chat/message"
	"time"
)

const (
	Group1    = "123456"
	UserTest1 = "+5518911111111"
	UserTest2 = "+5518922222222"
	UserTest3 = "+5518933333333"
)

type memoryRepository struct {
	users           map[string]struct{}
	usersOnline     map[string]string
	messagesOffline map[string][]*message.Message
	blockedUser     map[string]string
	groups          map[string]struct{}
	groupMembers    map[string][]string
}

func NewMemoryRepository() Repository {
	mr := &memoryRepository{
		users:           make(map[string]struct{}),
		usersOnline:     make(map[string]string),
		messagesOffline: make(map[string][]*message.Message),
		blockedUser:     make(map[string]string),
		groups:          make(map[string]struct{}),
		groupMembers:    make(map[string][]string),
	}

	mr.users[UserTest1] = struct{}{}
	mr.users[UserTest2] = struct{}{}
	mr.users[UserTest3] = struct{}{}
	mr.blockedUser[UserTest3] = UserTest1
	mr.groups[Group1] = struct{}{}
	mr.groupMembers[Group1] = []string{
		UserTest1,
		UserTest2,
		UserTest3,
	}
	return mr
}

func (mr *memoryRepository) generatesMessages(userID string) {
	var msgs []*message.Message

	for i := 0; i < 3; i++ {
		msg, _ := message.NewMessage(
			UserTest3,
			userID,
			"",
			message.ContentText,
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

func (mr *memoryRepository) GetMessagesOffline(userID string) ([]*message.Message, error) {
	//mr.generatesMessages(userID)
	return mr.messagesOffline[userID], nil
}

func (mr *memoryRepository) IsValidUser(fromID string, toID string) (bool, error) {
	bID := mr.blockedUser[toID]
	isValid := !(bID == fromID)

	if isValid {
		_, ok := mr.users[toID]
		isValid = ok
	}

	return isValid, nil
}

func (mr *memoryRepository) GetGroupMembers(groupID string) ([]string, error) {
	return mr.groupMembers[groupID], nil
}
