package repository

import (
	"context"
	"database/sql"
	"github.com/tsmweb/auth-service/infra/db"
	"github.com/tsmweb/auth-service/login"
)

// loginRepositoryPostgres implementation for login.Repository interface.
type loginRepositoryPostgres struct {
	dataBase db.Database
}

// NewLoginRepositoryPostgres creates a new instance of login.Repository.
func NewLoginRepositoryPostgres(db db.Database) login.Repository {
	return &loginRepositoryPostgres{dataBase: db}
}

// Login returns if ID and password are valid.
func (r *loginRepositoryPostgres) Login(ctx context.Context, login *login.Login) (bool, error) {
	ok := false

	stmt, err := r.dataBase.DB().PrepareContext(ctx,`
		SELECT true FROM login 
		WHERE user_id = $1 
		AND password = $2`)
	if err != nil {
		return ok, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, login.ID, login.Password).Scan(&ok)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return ok, nil
}

// Update login data in the data base.
func (r *loginRepositoryPostgres) Update(ctx context.Context, login *login.Login) (bool, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return false, err
	}

	stmt, err := txn.PrepareContext(ctx, `
		UPDATE login 
		SET password = $1, updated_at = $2
		WHERE user_id = $3`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		login.Password, login.UpdatedAt, login.ID)
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
