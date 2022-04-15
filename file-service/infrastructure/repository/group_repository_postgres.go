package repository

import (
	"context"
	"database/sql"
	"github.com/tsmweb/file-service/group"
	"github.com/tsmweb/file-service/infrastructure/db"
)

// groupRepositoryPostgres implementation for group.Repository interface.
type groupRepositoryPostgres struct {
	database db.Database
}

// NewGroupRepositoryPostgres creates a new instance of group.Repository.
func NewGroupRepositoryPostgres(db db.Database) group.Repository {
	return &groupRepositoryPostgres{database: db}
}

// ExistsGroup returns true if the group exists in the database.
func (r *groupRepositoryPostgres) ExistsGroup(ctx context.Context, groupID string) (bool, error) {
	stmt, err := r.database.DB().PrepareContext(ctx, `
		SELECT g.id 
		FROM "group" g 
		WHERE g.id = $1`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var ID string
	err = stmt.QueryRowContext(ctx, groupID).Scan(&ID)
	if (err != nil) && (err != sql.ErrNoRows) {
		return false, err
	}

	return groupID == ID, nil
}

// IsGroupMember returns true if the user is a member of the group.
func (r *groupRepositoryPostgres) IsGroupMember(ctx context.Context, groupID, userID string) (bool, error) {
	stmt, err := r.database.DB().PrepareContext(ctx, `
		SELECT gm.user_id
		FROM group_member gm 
		WHERE gm.group_id = $1
		AND gm.user_id = $2`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var id string
	if err = stmt.QueryRowContext(ctx, groupID, userID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return id == userID, nil
}

// IsGroupAdmin returns true if the user is a group administrator.
func (r *groupRepositoryPostgres) IsGroupAdmin(ctx context.Context, groupID, userID string) (bool, error) {
	stmt, err := r.database.DB().PrepareContext(ctx, `
		SELECT gm.admin 
		FROM group_member gm 
		WHERE gm.group_id = $1
		AND gm.user_id = $2`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	admin := false
	err = stmt.QueryRowContext(ctx, groupID, userID).Scan(&admin)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return admin, nil
}
