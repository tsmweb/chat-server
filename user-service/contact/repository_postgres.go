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

// Get returns the contact by userID and contactID.
func (r *repositoryPostgres) Get(userID, contactID string) (*Contact, error) {
	return nil, nil
}

// GetAll returns all contacts by userID.
func (r *repositoryPostgres) GetAll(userID string) ([]*Contact, error) {
	return nil, nil
}

// ExistsUser checks if the contact exists in the database.
func (r *repositoryPostgres) ExistsUser(ID string) (bool, error) {
	return false, nil
}

// GetPresence returns the presence status of the contact.
func (r *repositoryPostgres) GetPresence(userID, contactID string) (PresenceType, error) {
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
func (r *repositoryPostgres) Delete(userID, contactID string) (int, error) {
	return 0, nil
}

// Block adds a contact to the blocked contacts database.
func (r *repositoryPostgres) Block(userID, contactID string) (bool, error) {
	return false, nil
}

// Unblock removes a contact from the blocked contacts database.
func (r *repositoryPostgres) Unblock(userID, contactID string) (bool, error) {
	return false, nil
}