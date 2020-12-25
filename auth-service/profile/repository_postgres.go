package profile

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/tsmweb/auth-service/helper/database"
	"github.com/tsmweb/helper-go/cerror"
)

// repositoryPostgres implementation for Repository interface.
type repositoryPostgres struct {
	dataBase database.Database
}

// NewRepositoryPostgres creates a new instance of Repository.
func NewRepositoryPostgres(db database.Database) Repository {
	return &repositoryPostgres{dataBase: db}
}

// Get returns the profile by id.
func (r *repositoryPostgres) Get(ID string) (Profile, error) {
	profile := Profile{}

	err := r.dataBase.DB().QueryRow("SELECT ID, name, lastname FROM profile WHERE ID = $1", ID).Scan(
		&profile.ID,
		&profile.Name,
		&profile.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return profile, cerror.ErrNotFound
		}

		return profile, err
	}

	return profile, nil
}

// Create new profile in the data base.
func (r *repositoryPostgres) Create(profile Profile) error {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return err
	}

	_, err = txn.Exec(`INSERT INTO profile(id, name, lastname) VALUES($1, $2, $3)`,
		profile.ID, profile.Name, profile.LastName)
	if err != nil {
		txn.Rollback()
		//"23505": "unique_violation"
		if err.(*pq.Error).Code == pq.ErrorCode("23505") {
			return cerror.ErrRecordAlreadyRegistered
		}

		return err
	}

	_, err = txn.Exec(`INSERT INTO login(client_id, password) VALUES($1, $2)`,
		profile.ID, profile.Password)
	if err != nil {
		txn.Rollback()
		return err
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		return err
	}

	return nil
}

// Update profile data in the data base.
func (r *repositoryPostgres) Update(profile Profile) error {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return err
	}

	result, err := txn.Exec(`
		UPDATE profile 
		SET name = $1, lastname = $2, update_at = CURRENT_TIMESTAMP 
		WHERE id = $3`,
		profile.Name, profile.LastName, profile.ID)
	if err != nil {
		txn.Rollback()
		return err
	}

	ra, _ := result.RowsAffected()
	if ra != 1 {
		txn.Rollback()
		return cerror.ErrNotFound
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		return err
	}

	return nil
}
