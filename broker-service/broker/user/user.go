package user

import (
	"context"
	"time"
)

// Status type that represents the user's status as UserOnline and UserOffline.
type Status int

const (
	Online  Status = 0x1
	Offline Status = 0x2
)

func (s Status) String() (str string) {
	name := func(status Status, name string) bool {
		if s&status == 0 {
			return false
		}
		str = name
		return true
	}

	if name(Online, "online") {
		return
	}
	if name(Offline, "offline") {
		return
	}

	return
}

// Repository represents an abstraction of the data persistence layer.
type Repository interface {
	// AddUserPresence adds the user's presence to the database.
	AddUserPresence(ctx context.Context, userID string, serverID string, createAt time.Time) error

	// RemoveUserPresence removes user presence from database.
	RemoveUserPresence(ctx context.Context, userID string) error

	// GetUserServer returns the server the user is online.
	GetUserServer(ctx context.Context, userID string) (string, error)

	// IsValidUser returns true if the user is valid and false otherwise.
	IsValidUser(ctx context.Context, userID string) (bool, error)

	// IsBlockedUser returns true if the message sending user was blocked and false otherwise.
	IsBlockedUser(ctx context.Context, fromID string, toID string) (bool, error)

	// GetAllContactsOnline returns all online contacts by userID.
	GetAllContactsOnline(ctx context.Context, userID string) ([]string, error)

	// GetAllRelationshipsOnline returns all online users for which I am a contact.
	GetAllRelationshipsOnline(ctx context.Context, userID string) ([]string, error)
}

// User represents the status of the user's connection.
type User struct {
	ID       string
	Status   string
	ServerID string
	Date     time.Time
}

// New create and return an User instance.
func New(id string, status Status, serverID string) *User {
	return &User{
		ID:       id,
		Status:   status.String(),
		ServerID: serverID,
		Date:     time.Now().UTC(),
	}
}

// Encoder is a User encoder for byte slice.
type Encoder interface {
	Marshal(u *User) ([]byte, error)
}

// The EncoderFunc type is an adapter to allow the use of ordinary functions as encoders of User for byte slice.
// If f is a function with the appropriate signature, EncoderFunc(f) is a Encoder that calls f.
type EncoderFunc func(u *User) ([]byte, error)

// Marshal calls f(m).
func (f EncoderFunc) Marshal(u *User) ([]byte, error) {
	return f(u)
}

// Decoder is a byte slice decoder for User.
type Decoder interface {
	Unmarshal(in []byte, u *User) error
}

// The DecoderFunc type is an adapter to allow the use of ordinary functions as decoders of byte slice for User.
// If f is a function with the appropriate signature, DecoderFunc(f) is a Decoder that calls f.
type DecoderFunc func(in []byte, u *User) error

// Unmarshal calls f(in, m).
func (f DecoderFunc) Unmarshal(in []byte, u *User) error {
	return f(in, u)
}
