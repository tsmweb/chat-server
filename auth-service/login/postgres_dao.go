package login

import (
	"database/sql"
	"github.com/tsmweb/auth-service/helper/database"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/util/hashutil"
)

// postgresDAO implementation for DAO interface.
type postgresDAO struct {
	dataBase database.Database
}

// NewPostgresDAO creates a new instance of DAO.
func NewPostgresDAO(db database.Database) DAO {
	return &postgresDAO{dataBase: db}
}

// Login returns if ID and password are valid.
func (p *postgresDAO) Login(login Login) (bool, error) {
	var hashedPassword string
	err := p.dataBase.DB().QueryRow(`SELECT password FROM login WHERE client_id = $1`,
		login.ID).Scan(&hashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, cerror.ErrNotFound
		}

		return false, err
	}

	return hashutil.VerifySHA1(hashedPassword, login.Password)
}

// Update login data in the data base.
func (p *postgresDAO) Update(login Login) error {
	hashedPassword, err := hashutil.HashSHA1(login.Password)
	if err != nil {
		return err
	}

	txn, err := p.dataBase.DB().Begin()
	if err != nil {
		return err
	}

	_, err = txn.Exec(`
		UPDATE login 
		SET password = $1, update_at = CURRENT_TIMESTAMP
		WHERE client_id = $2`,
		hashedPassword, login.ID)
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
