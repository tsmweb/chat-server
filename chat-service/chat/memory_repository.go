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
	usersOnline     map[string]string
	messagesOffline map[string][]*message.Message
	blockedUser     map[string]string
	groupMembers    map[string][]string
}

func NewMemoryRepository() Repository {
	mr := &memoryRepository{
		usersOnline:     make(map[string]string),
		messagesOffline: make(map[string][]*message.Message),
		blockedUser:     make(map[string]string),
		groupMembers:    make(map[string][]string),
	}

	mr.blockedUser[UserTest3] = UserTest1
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
		msg, _ := message.New(
			UserTest3,
			userID,
			"",
			message.TEXT,
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

func (mr *memoryRepository) GetMessagesOffline(userID string) ([]*message.Message, error) {
	//mr.generatesMessages(userID)
	return mr.messagesOffline[userID], nil
}

func (mr *memoryRepository) IsBlockedUser(userID string, blockedID string) (bool, error) {
	bID := mr.blockedUser[userID]
	return bID == blockedID, nil
}

func (mr *memoryRepository) GetGroupMembers(groupID string) ([]string, error) {
	return mr.groupMembers[groupID], nil
}
