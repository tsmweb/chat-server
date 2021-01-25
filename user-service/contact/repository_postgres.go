package contact

import "github.com/tsmweb/use-service/helper/database"

// repositoryPostgres implementation for Repository interface.
type repositoryPostgres struct {
	dataBase database.Database
}

// NewRepositoryPostgres creates a new instance of Repository.
func NewRepositoryPostgres(db database.Database) Repository {
	return &repositoryPostgres{dataBase: db}
}

// Get returns the contact by profileID and contactID.
func (r *repositoryPostgres) Get(profileID, contactID string) (*Contact, error) {
	return nil, nil
}

// GetAll returns all contacts by profileID.
func (r *repositoryPostgres) GetAll(profileID string) ([]*Contact, error) {
	return nil, nil
}

// ExistsProfile checks if the contact exists in the database.
func (r *repositoryPostgres) ExistsProfile(ID string) (bool, error) {
	return false, nil
}

// GetPresence returns the presence status of the contact.
func (r *repositoryPostgres) GetPresence(profileID, contactID string) (PresenceType, error) {
	return NotFound, nil
}

// Create creates a new contact in the database.
func (r *repositoryPostgres) Create(contact *Contact) error {
	return nil
}

// Update updates the contact data in the database.
func (r *repositoryPostgres) Update(contact *Contact) (int, error) {
	return 0, nil
}

// Delete deletes a contact from the database.
func (r *repositoryPostgres) Delete(contact *Contact) (int, error) {
	return 0, nil
}

// Block adds a contact to the blocked contacts database.
func (r *repositoryPostgres) Block(profileID, contactID string) (bool, error) {
	return false, nil
}

// Unblock removes a contact from the blocked contacts database.
func (r *repositoryPostgres) Unblock(profileID, contactID string) (bool, error) {
	return false, nil
}