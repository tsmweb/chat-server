package contact

// Reader interface
type Reader interface {
	Get(profileID, contactID string) (Contact, error)
	GetAll(profileID string) ([]Contact, error)
	ExistsProfile(ID string) (bool, error)
}

// Writer contact writer
type Writer interface {
	Create(contact Contact) error
	Update(contact Contact) error
	Delete(contact Contact) error
}

// Repository interface for contact data source.
type Repository interface {
	Reader
	Writer
}
