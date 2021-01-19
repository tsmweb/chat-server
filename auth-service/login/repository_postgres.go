package login

import (
	"database/sql"
	"github.com/tsmweb/auth-service/helper/database"
	"github.com/tsmweb/go-helper-api/cerror"
)

// repositoryPostgres implementation for Repository interface.
type repositoryPostgres struct {
	dataBase database.Database
}

// NewRepositoryPostgres creates a new instance of Repository.
func NewRepositoryPostgres(db database.Database) Repository {
	return &repositoryPostgres{dataBase: db}
}

// Login returns if ID and password are valid.
func (r *repositoryPostgres) Login(login *Login) (bool, error) {
	err := r.dataBase.DB().
		QueryRow(`SELECT 1 FROM login WHERE client_id = $1 AND password = $2`, login.ID, login.Password).
		Err()

	if err != nil {
		if err == sql.ErrNoRows {
			return false, cerror.ErrNotFound
		}

		return false, err
	}

	return true, nil
}

// Update login data in the data base.
func (r *repositoryPostgres) Update(login *Login) (int, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return -1, err
	}

	result, err := txn.Exec(`
		UPDATE login 
		SET password = $1, update_at = CURRENT_TIMESTAMP
		WHERE client_id = $2`,
		login.Password, login.ID)
	if err != nil {
		txn.Rollback()
		return -1, err
	}

	ra, _ := result.RowsAffected()
	if ra != 1 {
		txn.Rollback()
		return 0, nil
	}

	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		return -1, err
	}

	return int(ra), nil
}
