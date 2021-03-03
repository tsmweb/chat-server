package user

import (
	"database/sql"
	"github.com/lib/pq"
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

// Get returns the user by id.
func (r *repositoryPostgres) Get(ID string) (*User, error) {
	stmt, err := r.dataBase.DB().Prepare(`
		SELECT ID, name, lastname, created_at, updated_at FROM "user" WHERE ID = $1`)
	if err != nil {
		return nil, err
	}

	var user User
	err = stmt.QueryRow(ID).Scan(
		&user.ID,
		&user.Name,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cerror.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

// Create new user in the data base.
func (r *repositoryPostgres) Create(user *User) error {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return err
	}

	_, err = txn.Exec(`INSERT INTO "user"(id, name, lastname, created_at) VALUES($1, $2, $3, $4)`,
		user.ID, user.Name, user.LastName, user.CreatedAt)
	if err != nil {
		txn.Rollback()
		//"23505": "unique_violation"
		if err.(*pq.Error).Code == pq.ErrorCode("23505") {
			return cerror.ErrRecordAlreadyRegistered
		}

		return err
	}

	_, err = txn.Exec(`INSERT INTO login(user_id, password, created_at) VALUES($1, $2, $3)`,
		user.ID, user.Password, user.CreatedAt)
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

// Update user data in the data base.
func (r *repositoryPostgres) Update(user *User) (int, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return -1, err
	}

	result, err := txn.Exec(`
		UPDATE "user" 
		SET name = $1, lastname = $2, updated_at = $3
		WHERE id = $4`,
		user.Name, user.LastName, user.UpdatedAt, user.ID)
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
