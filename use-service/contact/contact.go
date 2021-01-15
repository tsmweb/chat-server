package contact

// Contact data model
type Contact struct {
	ID string
	Name string
	LastName string
	ProfileID string
}

// NewContact create a new Contact
func NewContact(ID string, name string, lastname string, profileID string) (Contact, error) {
	c := Contact{
		ID: ID,
		Name: name,
		LastName: lastname,
		ProfileID: profileID,
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
	if c.ProfileID == "" {
		return ErrProfileIDValidateModel
	}

	return nil
}