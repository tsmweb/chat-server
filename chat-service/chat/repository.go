package chat

import (
	"github.com/tsmweb/chat-service/chat/message"
	"time"
)

// Repository interface for user data source.
type Repository interface {
	AddUserOnline(userID string, host string, createAt time.Time) error
	DeleteUserOnline(userID string) error
	GetUserOnline(userID string) (string, bool, error)
	GetMessagesOffline(userID string) ([]*message.Message, error)
	IsValidUser(fromID string, toID string) (bool, error)
	GetGroupMembers(groupID string) ([]string, error)
}
