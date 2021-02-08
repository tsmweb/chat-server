package contact

// Reader interface
type Reader interface {
	Get(userID, contactID string) (*Contact, error)
	GetAll(userID string) ([]*Contact, error)
	ExistsUser(ID string) (bool, error)
	GetPresence(userID, contactID string) (PresenceType, error)
}

// Writer contact writer
type Writer interface {
	Create(contact *Contact) error
	Update(contact *Contact) (int, error)
	Delete(userID, contactID string) (int, error)
	Block(userID, blockedUserID string) error
	Unblock(userID, blockedUserID string) (bool, error)
}

// Repository interface for contact data source.
type Repository interface {
	Reader
	Writer
}
