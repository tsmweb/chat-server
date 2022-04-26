package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/tsmweb/auth-service/app/user"
	"github.com/tsmweb/auth-service/infra/db"
	"github.com/tsmweb/go-helper-api/cerror"
	"time"
)

// userRepositoryPostgres implementation for user.Repository interface.
type userRepositoryPostgres struct {
	dataBase db.Database
}

// NewUserRepositoryPostgres creates a new instance of user.Repository.
func NewUserRepositoryPostgres(db db.Database) user.Repository {
	return &userRepositoryPostgres{dataBase: db}
}

// Get returns the user by id.
func (r *userRepositoryPostgres) Get(ctx context.Context, ID string) (*user.User, error) {
	stmt, err := r.dataBase.DB().PrepareContext(ctx, `
		SELECT u.id, 
			u.name, 
			u.lastname, 
			u.created_at,
			COALESCE(u.updated_at, u.created_at, u.updated_at) AS updated_at
		FROM "user" u WHERE u.id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user user.User
	err = stmt.QueryRowContext(ctx, ID).
		Scan(&user.ID,
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
func (r *userRepositoryPostgres) Create(ctx context.Context, user *user.User) error {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.PrepareContext(ctx, `
		INSERT INTO "user"(id, name, lastname, created_at) 
		VALUES($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		user.ID, user.Name, user.LastName, user.CreatedAt)
	if err != nil {
		txn.Rollback()
		//"23505": "unique_violation"
		if err.(*pq.Error).Code == pq.ErrorCode("23505") {
			return cerror.ErrRecordAlreadyRegistered
		}

		return err
	}

	err = r.addLogin(ctx, txn, user.ID, user.Password, user.CreatedAt)
	if err != nil {
		txn.Rollback()
		return err
	}

	if err = txn.Commit(); err != nil {
		txn.Rollback()
		return err
	}

	return nil
}

// Update user data in the data base.
func (r *userRepositoryPostgres) Update(ctx context.Context, user *user.User) (bool, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return false, err
	}

	stmt, err := txn.PrepareContext(ctx, `
		UPDATE "user" 
		SET name = $1, lastname = $2, updated_at = $3
		WHERE id = $4`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		user.Name, user.LastName, user.UpdatedAt, user.ID)
	if err != nil {
		txn.Rollback()
		return false, err
	}

	ra, _ := result.RowsAffected()
	if ra != 1 {
		txn.Rollback()
		return false, nil
	}

	if err = txn.Commit(); err != nil {
		txn.Rollback()
		return false, err
	}

	return true, nil
}

func (r *userRepositoryPostgres) addLogin(ctx context.Context,
	txn *sql.Tx, userID string, password string, createdAt time.Time) error {
	stmt, err := txn.PrepareContext(ctx, `
		INSERT INTO login(user_id, password, created_at) 
		VALUES($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		userID, password, createdAt)

	return err
}
