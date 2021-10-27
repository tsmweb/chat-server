package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/infra/db"
	"time"
)

const (
	groupMembersKey        = "group:members:%s"
	groupMembersExpiration = time.Minute * 15

	insertedStatusMessage  = "I"
	processedStatusMessage = "P"
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
	members, _ := r.cache.SMembers(ctx, _groupMembersKey)
	if members != nil {
		return members, nil
	}

	members, err := r.getAllGroupMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if members == nil {
		return nil, nil
	}

	if err = r.cache.SAdd(ctx, _groupMembersKey, members); err != nil {
		return nil, err
	}
	if err = r.cache.Expire(ctx, _groupMembersKey, groupMembersExpiration); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *messageRepository) getAllGroupMembers(ctx context.Context, groupID string) ([]string, error) {
	stmt, err := r.database.DB().PrepareContext(ctx, `
		SELECT user_id
		FROM group_member
		WHERE group_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []string

	for rows.Next() {
		var member string
		if err = rows.Scan(&member); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}

		members = append(members, member)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return members, nil
}

// GetAllMessages returns all offline messages by user ID.
func (r *messageRepository) GetAllMessages(ctx context.Context, userID string) ([]*message.Message, error) {
	if err := r.updateMessageStatusToProcessed(ctx, userID); err != nil {
		return nil, err
	}

	stmt, err := r.database.DB().PrepareContext(ctx, `
		SELECT msg_id, msg_from, msg_to, msg_group, msg_date, msg_content_type, msg_content
		FROM offline_message 
		WHERE msg_to = $1
		AND msg_status = $2`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID, processedStatusMessage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*message.Message

	for rows.Next() {
		var msg message.Message
		err = rows.Scan(
			&msg.ID,
			&msg.From,
			&msg.To,
			&msg.Group,
			&msg.Date,
			&msg.ContentType,
			&msg.Content)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}

		messages = append(messages, &msg)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messageRepository) updateMessageStatusToProcessed(ctx context.Context, userID string) error {
	txn, err := r.database.DB().Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.PrepareContext(ctx, `
		UPDATE offline_message
		SET msg_status = $1
		WHERE msg_to = $2
		AND msg_status = $3`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, processedStatusMessage, userID, insertedStatusMessage)
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

// AddMessage add a message to the database.
func (r *messageRepository) AddMessage(ctx context.Context, msg message.Message) error {
	txn, err := r.database.DB().Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.PrepareContext(ctx, `
		INSERT INTO offline_message(
			msg_id, msg_status, msg_from, msg_to, msg_group, msg_date, msg_content_type, msg_content)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx, msg.ID, insertedStatusMessage, msg.From, msg.To, msg.Group, msg.Date, msg.ContentType, msg.Content)
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

// DeleteAllMessages deletes all messages by userID.
func (r *messageRepository) DeleteAllMessages(ctx context.Context, userID string) error {
	txn, err := r.database.DB().Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.PrepareContext(ctx, `
		DELETE FROM offline_message
		WHERE msg_to = $1
		AND msg_status = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID, processedStatusMessage)
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
