package contact

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/use-service/helper/database"
	"time"
)

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
	stmt, err := r.dataBase.DB().Prepare(`
			SELECT user_id, contact_id, name, lastname, created_at, updated_at 
			FROM contact 
			WHERE user_id = $1 
			  AND contact_id = $2`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var contact Contact
	err = stmt.QueryRow(userID, contactID).
		Scan(&contact.UserID,
			&contact.ID,
			&contact.Name,
			&contact.LastName,
			&contact.CreatedAt,
			&contact.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cerror.ErrNotFound
		}
		return nil, err
	}

	return &contact, nil
}

// GetAll returns all contacts by userID.
func (r *repositoryPostgres) GetAll(userID string) ([]*Contact, error) {
	stmt, err := r.dataBase.DB().Prepare(`
			SELECT user_id, contact_id, name, lastname, created_at, updated_at 
			FROM contact 
			WHERE user_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	contacts := make([]*Contact, 0)

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contact Contact
		err = rows.Scan(
			&contact.UserID,
			&contact.ID,
			&contact.Name,
			&contact.LastName,
			&contact.CreatedAt,
			&contact.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, cerror.ErrNotFound
			}
			return nil, err
		}

		contacts = append(contacts, &contact)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return contacts, nil
}

// ExistsUser checks if the contact exists in the database.
func (r *repositoryPostgres) ExistsUser(ID string) (bool, error) {
	stmt, err := r.dataBase.DB().Prepare(`SELECT id FROM "user" WHERE id = $1`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var userID string
	err = stmt.QueryRow(ID).Scan(&userID)
	if (err != nil) && (err != sql.ErrNoRows) {
		return false, err
	}

	return userID == ID, nil
}

// GetPresence returns the presence status of the contact.
func (r *repositoryPostgres) GetPresence(userID, contactID string) (PresenceType, error) {
	stmt, err := r.dataBase.DB().Prepare(`
			SELECT CASE WHEN o.user_id IS NULL THEN 'F' ELSE 'T' END as status
			FROM contact c
			LEFT JOIN online_user o ON c.contact_id = o.user_id
			WHERE c.user_id = $1
  			  AND c.contact_id = $2
  			  AND NOT EXISTS (
        		SELECT 1 FROM blocked_user b
        		WHERE c.user_id = b.blocked_user_id
          		  AND c.contact_id = b.user_id
    		  )`)
	if err != nil {
		return NotFound, err
	}
	defer stmt.Close()

	online := "N"
	err = stmt.QueryRow(userID, contactID).Scan(&online)
	if (err != nil) && (err != sql.ErrNoRows) {
		return NotFound, err
	}

	if online == "F" {
		return Offline, nil
	}
	if online == "T" {
		return Online, nil
	}
	return NotFound, nil
}

// Create creates a new contact in the database.
func (r *repositoryPostgres) Create(contact *Contact) error {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return err
	}

	_, err = txn.Exec(`
		INSERT INTO contact(user_id, contact_id, name, lastname, created_at) 
		VALUES($1, $2, $3, $4, $5)`,
		contact.UserID, contact.ID, contact.Name, contact.LastName, contact.CreatedAt)
	if err != nil {
		txn.Rollback()
		// "23505": "unique_violation"
		if err.(*pq.Error).Code == "23505" {
			return cerror.ErrRecordAlreadyRegistered
		}

		return err
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		return err
	}

	return nil
}

// Update updates the contact data in the database.
func (r *repositoryPostgres) Update(contact *Contact) (int, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return -1, err
	}

	result, err := txn.Exec(`
		UPDATE contact 
		SET name = $1, 
		    lastname = $2, 
		    updated_at = $3
		WHERE user_id = $4
		  AND contact_id = $5`,
		  contact.Name, contact.LastName, contact.UpdatedAt, contact.UserID, contact.ID)
	if err != nil {
		txn.Rollback()
		return -1, err
	}

	ra, _ := result.RowsAffected()
	if ra != 1 {
		txn.Rollback()
		return int(ra), nil
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		return -1, err
	}

	return int(ra), nil
}

// Delete deletes a contact from the database.
func (r *repositoryPostgres) Delete(userID, contactID string) (int, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return -1, err
	}

	result, err := txn.Exec(`
		DELETE FROM contact 
		WHERE user_id = $1 
		  AND contact_id = $2`,
		  userID, contactID)
	if err != nil {
		txn.Rollback()
		return -1, err
	}

	ra, _ := result.RowsAffected()
	if ra != 1 {
		txn.Rollback()
		return int(ra), nil
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		return -1, err
	}

	return int(ra), nil
}

// Block adds a contact to the blocked contacts database.
func (r *repositoryPostgres) Block(userID, blockedUserID string, createdAt time.Time) error {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return err
	}

	_, err = txn.Exec(`
		INSERT INTO blocked_user(user_id, blocked_user_id, created_at) 
		VALUES($1, $2, $3)`,
		userID, blockedUserID, createdAt)
	if err != nil {
		txn.Rollback()
		// "23505": "unique_violation"
		if err.(*pq.Error).Code == "23505" {
			return cerror.ErrRecordAlreadyRegistered
		}

		return err
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		return err
	}

	return nil
}

// Unblock removes a contact from the blocked contacts database.
func (r *repositoryPostgres) Unblock(userID, blockedUserID string) (bool, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return false, err
	}

	result, err := txn.Exec(`
		DELETE FROM blocked_user 
		WHERE user_id = $1
		  AND blocked_user_id = $2`,
		  userID, blockedUserID)
	if err != nil {
		txn.Rollback()
		return false, err
	}

	ra, _ := result.RowsAffected()
	if ra != 1 {
		txn.Rollback()
		return false, nil
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		return false, err
	}

	return true, nil
}