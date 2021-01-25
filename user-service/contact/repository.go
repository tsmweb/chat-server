package contact

// Reader interface
type Reader interface {
	Get(profileID, contactID string) (*Contact, error)
	GetAll(profileID string) ([]*Contact, error)
	ExistsProfile(ID string) (bool, error)
	GetPresence(profileID, contactID string) (PresenceType, error)
}

// Writer contact writer
type Writer interface {
	Create(contact *Contact) error
	Update(contact *Contact) (int, error)
	Delete(contact *Contact) (int, error)
	Block(profileID, contactID string) (bool, error)
	Unblock(profileID, contactID string) (bool, error)
}

// Repository interface for contact data source.
type Repository interface {
	Reader
	Writer
}
