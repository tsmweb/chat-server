package dto

import (
	"github.com/tsmweb/user-service/contact"
	"time"
)

// Contact data
type Contact struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastname"`
	UserID    string    `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// ToEntity mapper dto.Contact to contact.Contact
func (c *Contact) ToEntity() *contact.Contact {
	return &contact.Contact{
		ID:       c.ID,
		Name:     c.Name,
		LastName: c.LastName,
		UserID:   c.UserID,
	}
}

// FromEntity mapper contact.Contact to dto.Contact
func (c *Contact) FromEntity(entity *contact.Contact) {
	c.ID = entity.ID
	c.Name = entity.Name
	c.LastName = entity.LastName
	c.UserID = entity.UserID
	c.CreatedAt = entity.CreatedAt
	c.UpdatedAt = entity.UpdatedAt
}

// Presence data
type Presence struct {
	ID       string `json:"id"`
	Presence string `json:"presence"`
}

// EntityToContactDTO mapper []contact.Contact to []dto.Contact
func EntityToContactDTO(entities ...*contact.Contact) []*Contact {
	var contacts []*Contact

	for _, contact := range entities {
		c := &Contact{}
		c.FromEntity(contact)
		contacts = append(contacts, c)
	}

	return contacts
}
