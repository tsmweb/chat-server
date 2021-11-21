package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/infra/db"
	"strconv"
	"time"
)

const (
	validUserKey        = "valid:user:%s"
	validUserTrue       = "true"
	validUserFalse      = "false"
	validUserExpiration = time.Minute * 30

	blockedUserKey        = "blocked:user:%s:%s"
	blockedUserTrue       = "true"
	blockedUserFalse      = "false"
	blockedUserExpiration = time.Minute * 30
)

// userRepository implementation for user.Repository interface.
type userRepository struct {
	database db.Database
	cache    db.CacheDB
}

// NewUserRepository creates a new instance of user.Repository.
func NewUserRepository(database db.Database, cache db.CacheDB) user.Repository {
	return &userRepository{
		database: database,
		cache:    cache,
	}
}

// AddUserPresence adds the user's presence to the database.
func (r *userRepository) AddUserPresence(ctx context.Context, userID string, serverID string, createAt time.Time) error {
	txn, err := r.database.DB().Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.PrepareContext(ctx, `
		INSERT INTO online_user(user_id, server_id, created_at)
		VALUES($1, $2, $3)
		ON CONFLICT(user_id)
		DO UPDATE SET
			server_id = $2,
			created_at = $3`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID, serverID, createAt)
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

// RemoveUserPresence removes user presence from database.
func (r *userRepository) RemoveUserPresence(ctx context.Context, userID string) error {
	txn, err := r.database.DB().Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.PrepareContext(ctx, `
		DELETE FROM online_user
		WHERE user_id = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID)
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

// UpdateUserPresenceCache updates user presence in cache.
func (r *userRepository) UpdateUserPresenceCache(ctx context.Context, userID string, serverID string, status string) error {
	if user.Online.String() == status {
		return r.cache.Set(ctx, userID, serverID, 0)
	}
	return r.cache.Del(ctx, userID)
}

// GetUserServer returns the server the user is online.
func (r *userRepository) GetUserServer(ctx context.Context, userID string) (string, error) {
	return r.cache.Get(ctx, userID)
}

// IsValidUser returns true if the user is valid and false otherwise.
func (r *userRepository) IsValidUser(ctx context.Context, userID string) (bool, error) {
	_validUserKey := fmt.Sprintf(validUserKey, userID)
	isValidStr, err := r.cache.Get(ctx, _validUserKey)
	if err != nil {
		return false, err
	}
	if isValidStr == validUserTrue {
		return true, nil
	}
	if isValidStr == validUserFalse {
		return false, nil
	}

	isValid, err := r.isValidUser(ctx, userID)
	if err != nil {
		return false, err
	}

	if err = r.cache.Set(ctx, _validUserKey, strconv.FormatBool(isValid), validUserExpiration); err != nil {
		return false, err
	}

	return isValid, nil
}

func (r *userRepository) isValidUser(ctx context.Context, userID string) (bool, error) {
	stmt, err := r.database.DB().PrepareContext(ctx, `SELECT id FROM "user" WHERE id = $1`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var _userID string
	err = stmt.QueryRowContext(ctx, userID).Scan(&_userID)
	if (err != nil) && (err != sql.ErrNoRows) {
		return false, err
	}

	return userID == _userID, nil
}

// IsBlockedUser returns true if the message sending user was blocked and false otherwise.
func (r *userRepository) IsBlockedUser(ctx context.Context, userID string, blockedUserID string) (bool, error) {
	_blockedUserKey := fmt.Sprintf(blockedUserKey, userID, blockedUserID)
	isBlockedStr, err := r.cache.Get(ctx, _blockedUserKey)
	if err != nil {
		return false, err
	}
	if isBlockedStr == blockedUserTrue {
		return true, nil
	}
	if isBlockedStr == blockedUserFalse {
		return false, nil
	}

	isBlocked, err := r.isBlockedUser(ctx, userID, blockedUserID)
	if err != nil {
		return false, err
	}

	if err = r.cache.Set(ctx, _blockedUserKey, strconv.FormatBool(isBlocked), blockedUserExpiration); err != nil {
		return false, err
	}

	return isBlocked, nil
}

func (r *userRepository) isBlockedUser(ctx context.Context, userID string, blockedUserID string) (bool, error) {
	stmt, err := r.database.DB().PrepareContext(ctx, `
		SELECT blocked_user_id 
		FROM blocked_user 
		WHERE user_id = $1
		AND blocked_user_id = $2`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var _blockedUserID string
	err = stmt.QueryRowContext(ctx, userID, blockedUserID).Scan(&_blockedUserID)
	if (err != nil) && (err != sql.ErrNoRows) {
		return false, err
	}

	return _blockedUserID == blockedUserID, nil
}

// UpdateBlockedUserCache refresh blocked users cache.
func (r *userRepository) UpdateBlockedUserCache(ctx context.Context, userID string,
	blockedUserID string, blocked bool) error {
	_blockedUserKey := fmt.Sprintf(blockedUserKey, userID, blockedUserID)

	if r.cache.Key(ctx, _blockedUserKey) {
		if err := r.cache.Set(ctx, _blockedUserKey, strconv.FormatBool(blocked), blockedUserExpiration); err != nil {
			return err
		}
	}

	return nil
}

// GetAllContactsOnline returns all online contacts by userID.
func (r *userRepository) GetAllContactsOnline(ctx context.Context, userID string) ([]string, error) {
	stmt, err := r.database.DB().PrepareContext(ctx, `
		SELECT c.contact_id
		FROM contact c
		INNER JOIN online_user u ON u.user_id = c.contact_id
		WHERE c.user_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var contacts []string

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contactID string
		if err = rows.Scan(&contactID); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}

		contacts = append(contacts, contactID)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return contacts, nil
}

// GetAllRelationshipsOnline returns all online users for which I am a contact.
func (r *userRepository) GetAllRelationshipsOnline(ctx context.Context, userID string) ([]string, error) {
	stmt, err := r.database.DB().PrepareContext(ctx, `
		SELECT c.user_id
		FROM contact c
		INNER JOIN online_user ou ON ou.user_id = c.user_id
		WHERE c.contact_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var relationships []string

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var _userID string
		if err = rows.Scan(&_userID); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}

		relationships = append(relationships, _userID)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return relationships, nil
}
