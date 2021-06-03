package core

import (
	"time"
)

// Repository interface for user data source.
type Repository interface {
	AddUserOnline(userID string, host string, createAt time.Time) error
	DeleteUserOnline(userID string) error
	GetUserOnline(userID string) (string, bool, error)
	GetMessagesOffline(userID string) ([]*Message, error)
	IsBlockedUser(userID string, blockedID string) (bool, error)
}
