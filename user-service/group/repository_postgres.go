package group

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/user-service/helper/database"
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

// Get returns group by groupID and userID.
func (r *repositoryPostgres) Get(ctx context.Context, groupID, userID string) (*Group, error) {
	stmt, err := r.dataBase.DB().PrepareContext(ctx, `
		SELECT g.id,
			g.owner_id,
			g.name,
			g.description,
			g.created_at,
			COALESCE(g.updated_at, g.created_at, g.updated_at) AS updated_at,
			COALESCE(g.updated_by, '', g.updated_by) AS updated_by
		FROM "group" g
		INNER JOIN group_member gm ON g.id = gm.group_id
		WHERE g.id = $1
		AND gm.user_id = $2`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var group Group
	err = stmt.QueryRowContext(ctx, groupID, userID).
		Scan(&group.ID,
			&group.Owner,
			&group.Name,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, cerror.ErrNotFound
		}
		return nil, err
	}

	members, err := r.getAllMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}
	group.Members = members

	return &group, nil
}

// GetAll returns groups by userID.
func (r *repositoryPostgres) GetAll(ctx context.Context, userID string) ([]*Group, error) {
	stmt, err := r.dataBase.DB().PrepareContext(ctx, `
		SELECT g.id,
			g.owner_id,
			g.name,
			g.description,
			g.created_at,
			COALESCE(g.updated_at, g.created_at, g.updated_at) AS updated_at,
			COALESCE(g.updated_by, '', g.updated_by) AS updated_by
		FROM "group" g
		INNER JOIN group_member gm ON g.id = gm.group_id
		WHERE gm.user_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	groups := make([]*Group, 0)

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group Group
		err = rows.Scan(&group.ID,
			&group.Owner,
			&group.Name,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.UpdatedBy)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, cerror.ErrNotFound
			}
			return nil, err
		}

		members, err := r.getAllMembers(ctx, group.ID)
		if err != nil {
			return nil, err
		}
		group.Members = members

		groups = append(groups, &group)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return groups, nil
}

// ExistsUser returns true if the user exists in the database.
func (r *repositoryPostgres) ExistsUser(ctx context.Context, userID string) (bool, error) {
	stmt, err := r.dataBase.DB().PrepareContext(ctx, `
		SELECT u.id 
		FROM "user" u 
		WHERE u.id = $1`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var ID string
	err = stmt.QueryRowContext(ctx, userID).Scan(&ID)
	if (err != nil) && (err != sql.ErrNoRows) {
		return false, err
	}

	return userID == ID, nil
}

// ExistsGroup returns true if the group exists in the database.
func (r *repositoryPostgres) ExistsGroup(ctx context.Context, groupID string) (bool, error) {
	stmt, err := r.dataBase.DB().PrepareContext(ctx, `
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

// IsGroupAdmin returns true if the user is a group administrator.
func (r *repositoryPostgres) IsGroupAdmin(ctx context.Context, groupID, userID string) (bool, error) {
	stmt, err := r.dataBase.DB().PrepareContext(ctx, `
		SELECT gm.admin 
		FROM group_member gm 
		WHERE gm.group_id = $1
		AND gm.user_id = $2`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var admin bool
	err = stmt.QueryRowContext(ctx, groupID, userID).Scan(&admin)
	if (err != nil) && (err != sql.ErrNoRows) {
		return false, err
	}

	return admin, nil
}

// IsGroupOwner returns true if the user owns the group.
func (r *repositoryPostgres) IsGroupOwner(ctx context.Context, groupID, userID string) (bool, error) {
	stmt, err := r.dataBase.DB().PrepareContext(ctx, `
		SELECT true 
		FROM "group" g 
		WHERE g.id = $1
		AND g.owner_id = $2`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var isOwner bool
	err = stmt.QueryRowContext(ctx, groupID, userID).Scan(&isOwner)
	if (err != nil) && (err != sql.ErrNoRows) {
		return false, err
	}

	return isOwner, nil
}

// Create creates a new group in the database.
func (r *repositoryPostgres) Create(ctx context.Context, group *Group) error {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return err
	}

	// insert group
	stmt, err := txn.PrepareContext(ctx, `
		INSERT INTO "group"(id, owner_id, name, description, created_at)
		VALUES($1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		group.ID, group.Owner, group.Name, group.Description, group.CreatedAt)
	if err != nil {
		txn.Rollback()
		// "23505": "unique_violation"
		if err.(*pq.Error).Code == "23505" {
			return cerror.ErrRecordAlreadyRegistered
		}

		return err
	}

	// insert member, the group owner is also a permanent member.
	member, err := NewMember(group.ID, group.Owner, true)
	if err != nil {
		return err
	}

	err = r.addMember(ctx, txn, member)
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

// Update updates the group data in the database.
func (r *repositoryPostgres) Update(ctx context.Context, group *Group) (bool, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return false, err
	}

	stmt, err := txn.PrepareContext(ctx, `
		UPDATE "group" 
		SET name = $1, 
		    description = $2, 
		    updated_at = $3,
			updated_by = $4
		WHERE id = $5`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		group.Name, group.Description, group.UpdatedAt, group.UpdatedBy, group.ID)
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

// Delete deletes a group from the database.
func (r *repositoryPostgres) Delete(ctx context.Context, groupID string) (bool, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return false, err
	}

	err = r.addAllMembersToNotify(ctx, txn, groupID)
	if err != nil {
		txn.Rollback()
		return false, err
	}

	err = r.removeAllMembers(ctx, txn, groupID)
	if err != nil {
		txn.Rollback()
		return false, err
	}

	stmt, err := txn.PrepareContext(ctx, `
		DELETE FROM "group"
		WHERE id = $1`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, groupID)
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

// AddMember add a member to the group.
func (r *repositoryPostgres) AddMember(ctx context.Context, member *Member) error {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return err
	}

	err = r.addMember(ctx, txn, member)
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

// SetAdmin elevates a member to administrator status.
func (r *repositoryPostgres) SetAdmin(ctx context.Context, member *Member) (bool, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return false, err
	}

	stmt, err := txn.PrepareContext(ctx, `
		UPDATE group_member
		SET admin = $1,
		    updated_at = $2,
			updated_by = $3
		WHERE group_id = $4
		AND user_id = $5`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		member.Admin, member.UpdatedAt, member.UpdatedBy, member.GroupID, member.UserID)
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

// RemoveMember removes a member from the group.
func (r *repositoryPostgres) RemoveMember(ctx context.Context, groupID, userID string) (bool, error) {
	txn, err := r.dataBase.DB().Begin()
	if err != nil {
		return false, err
	}

	stmt, err := txn.PrepareContext(ctx, `
		DELETE FROM group_member
		WHERE group_id = $1
		AND user_id = $2`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, groupID, userID)
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

func (r *repositoryPostgres) addMember(ctx context.Context, txn *sql.Tx, member *Member) error {
	stmt, err := txn.PrepareContext(ctx, `
		INSERT INTO group_member(group_id, user_id, admin, created_at)
		VALUES($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		member.GroupID, member.UserID, member.Admin, member.CreatedAt)

	return err
}

func (r *repositoryPostgres) removeAllMembers(ctx context.Context, txn *sql.Tx, groupID string) error {
	stmt, err := txn.PrepareContext(ctx, `
		DELETE FROM group_member
		WHERE group_id = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, groupID)

	return err
}

func (r *repositoryPostgres) getAllMembers(ctx context.Context, groupID string) ([]*Member, error) {
	stmt, err := r.dataBase.DB().PrepareContext(ctx, `
		SELECT gm.group_id,
			gm.user_id,
			gm.admin,
			gm.created_at,
			COALESCE(gm.updated_at, gm.created_at, gm.updated_at) AS updated_at,
			COALESCE(gm.updated_by, '', gm.updated_by) AS updated_by
		FROM group_member gm
		WHERE gm.group_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	members := make([]*Member, 0)

	rows, err := stmt.QueryContext(ctx, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member Member
		err = rows.Scan(
			&member.GroupID,
			&member.UserID,
			&member.Admin,
			&member.CreatedAt,
			&member.UpdatedAt,
			&member.UpdatedBy)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, cerror.ErrNotFound
			}
			return nil, err
		}

		members = append(members, &member)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return members, nil
}

func (r *repositoryPostgres) addAllMembersToNotify(ctx context.Context, txn *sql.Tx, groupID string) error {
	members, err := r.getAllMembers(ctx, groupID)
	if err != nil {
		return err
	}

	if len(members) == 1 {
		return nil
	}

	stmt, err := txn.PrepareContext(ctx, `
		INSERT INTO group_member_notify(group_id, user_id, created_at)
		VALUES($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, member := range members {
		_, err = stmt.ExecContext(ctx,
			member.GroupID, member.UserID, time.Now().UTC())
		if err != nil {
			return err
		}
	}

	return nil
}
