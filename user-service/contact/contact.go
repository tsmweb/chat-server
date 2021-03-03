package contact

import "time"

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
		CreatedAt: time.Now(),
	}

	err := c.Validate()
	if err != nil {
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
