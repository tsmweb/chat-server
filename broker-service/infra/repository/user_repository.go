package repository

import (
	"context"
	"fmt"
	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/infra/db"
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
	blockedUserExpiration = time.Minute * 5
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
	if err := r.cache.Set(ctx, userID, serverID, 0); err != nil {
		return err
	}

	//TODO

	return nil
}

// RemoveUserPresence removes user presence from database.
func (r *userRepository) RemoveUserPresence(ctx context.Context, userID string) error {
	if err := r.cache.Del(ctx, userID); err != nil {
		return err
	}

	//TODO

	return nil
}

// GetUserServer returns the server the user is online.
func (r *userRepository) GetUserServer(ctx context.Context, userID string) (string, error) {
	return r.cache.Get(ctx, userID)
}

// IsValidUser returns true if the user is valid and false otherwise.
func (r *userRepository) IsValidUser(ctx context.Context, userID string) (bool, error) {
	_validUserKey := fmt.Sprintf(validUserKey, userID)
	isValid, err := r.cache.Get(ctx, _validUserKey)
	if err != nil {
		return false, err
	}
	if isValid == validUserTrue {
		return true, nil
	}
	if isValid == validUserFalse {
		return false, err
	}

	//TODO

	if err = r.cache.Set(ctx, _validUserKey, validUserTrue, validUserExpiration); err != nil {
		return false, err
	}

	return false, nil
}

// IsBlockedUser returns true if the message sending user was blocked and false otherwise.
func (r *userRepository) IsBlockedUser(ctx context.Context, fromID string, toID string) (bool, error) {
	_blockedUserKey := fmt.Sprintf(blockedUserKey, fromID, toID)
	isBlocked, err := r.cache.Get(ctx, _blockedUserKey)
	if err != nil {
		return false, err
	}
	if isBlocked == blockedUserTrue {
		return true, nil
	}
	if isBlocked == blockedUserFalse {
		return false, err
	}

	//TODO

	if err = r.cache.Set(ctx, _blockedUserKey, blockedUserFalse, blockedUserExpiration); err != nil {
		return false, err
	}

	return false, nil
}

// GetAllContactsOnline returns all online contacts by userID.
func (r *userRepository) GetAllContactsOnline(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}

// GetAllRelationshipsOnline returns all online users for which I am a contact.
func (r *userRepository) GetAllRelationshipsOnline(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}
