package repository

import (
	"context"
	"fmt"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/infra/db"
	"time"
)

const (
	groupMembersKey        = "group:members:%s"
	groupMembersExpiration = time.Minute * 15
)

// userRepository implementation for message.Repository interface.
type messageRepository struct {
	database db.Database
	cache    db.CacheDB
}

// NewMessageRepository creates a new instance of message.Repository.
func NewMessageRepository(database db.Database, cache db.CacheDB) message.Repository {
	return &messageRepository{
		database: database,
		cache:    cache,
	}
}

// GetAllGroupMembers returns all members of a group by groupID.
func (r *messageRepository) GetAllGroupMembers(ctx context.Context, groupID string) ([]string, error) {
	_groupMembersKey := fmt.Sprintf(groupMembersKey, groupID)
	members, err := r.cache.SMembers(ctx, _groupMembersKey)
	if err != nil {
		return nil, err
	}
	if members != nil {
		return members, nil
	}

	//TODO
	if err = r.cache.SAdd(ctx, _groupMembersKey, members); err != nil {
		return nil, err
	}
	if err = r.cache.Expire(ctx, _groupMembersKey, groupMembersExpiration); err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAllMessages returns all offline messages by user ID.
func (r *messageRepository) GetAllMessages(ctx context.Context, userID string) ([]*message.Message, error) {
	return nil, nil
}

// AddMessage add a message to the database.
func (r *messageRepository) AddMessage(ctx context.Context, msg message.Message) error {
	return nil
}

// DeleteAllMessages deletes all messages by userID.
func (r *messageRepository) DeleteAllMessages(ctx context.Context, userID string) error {
	return nil
}
