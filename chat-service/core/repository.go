package core

import "time"

type Repository interface {
	AddUserOnline(userID string, host string, createAt time.Time) error
	DeleteUserOnline(userID string) error
	GetUserOnline(userID string) (string, bool, error)
	GetMessagesOffline(userID string) ([]*Message, error)
}
