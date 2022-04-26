package contact

import (
	"context"
	"github.com/tsmweb/go-helper-api/cerror"
	"time"
)

var (
	ErrIDValidateModel     = &cerror.ErrValidateModel{Msg: "required id"}
	ErrUserIDValidateModel = &cerror.ErrValidateModel{Msg: "required user_id"}
)

// Contact data model
type Contact struct {
	ID        string
	Name      string
	LastName  string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewContact create a new Contact
func NewContact(ID, name, lastname, userID string) (*Contact, error) {
	c := &Contact{
		ID:        ID,
		Name:      name,
		LastName:  lastname,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}

	if err := c.Validate(); err != nil {
		return c, err
	}

	return c, nil
}

// Validate model Contact
func (c Contact) Validate() error {
	if c.ID == "" {
		return ErrIDValidateModel
	}
	if c.UserID == "" {
		return ErrUserIDValidateModel
	}

	return nil
}

// Repository interface for contact data source.
type Repository interface {
	Get(ctx context.Context, userID, contactID string) (*Contact, error)
	GetAll(ctx context.Context, userID string) ([]*Contact, error)
	ExistsUser(ctx context.Context, ID string) (bool, error)
	GetPresence(ctx context.Context, userID, contactID string) (PresenceType, error)
	Create(ctx context.Context, contact *Contact) error
	Update(ctx context.Context, contact *Contact) (bool, error)
	Delete(ctx context.Context, userID, contactID string) (bool, error)
	Block(ctx context.Context, userID, blockedUserID string, createdAt time.Time) error
	Unblock(ctx context.Context, userID, blockedUserID string) (bool, error)
}
