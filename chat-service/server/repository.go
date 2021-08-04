package server

import (
	"time"
)

// Repository interface for user data source.
type Repository interface {
	AddUserOnline(userID string, host string, createAt time.Time) error
	DeleteUserOnline(userID string) error
	IsValidUser(fromID string, toID string) (bool, error)
	GetGroupMembers(groupID string) ([]string, error)
	GetUserContactsOnline(userID string) ([]string, error)
	GetContactsWithUserOnline(userID string) ([]string, error)
}
