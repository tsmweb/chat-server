package contact

import (
	"context"
	"time"
)

// Reader interface
type Reader interface {
	Get(ctx context.Context, userID, contactID string) (*Contact, error)
	GetAll(ctx context.Context, userID string) ([]*Contact, error)
	ExistsUser(ctx context.Context, ID string) (bool, error)
	GetPresence(ctx context.Context, userID, contactID string) (PresenceType, error)
}

// Writer contact writer
type Writer interface {
	Create(ctx context.Context, contact *Contact) error
	Update(ctx context.Context, contact *Contact) (bool, error)
	Delete(ctx context.Context, userID, contactID string) (bool, error)
	Block(ctx context.Context, userID, blockedUserID string, createdAt time.Time) error
	Unblock(ctx context.Context, userID, blockedUserID string) (bool, error)
}

// Repository interface for contact data source.
type Repository interface {
	Reader
	Writer
}
