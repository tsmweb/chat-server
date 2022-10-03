package broker

import (
	"context"
	"time"

	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
)

// UserHandler handles user.
type UserHandler interface {
	// Execute performs user handling.
	Execute(ctx context.Context, usr user.User, chMessage chan<- message.Message) error
}

type userHandler struct {
	userRepository user.Repository
	msgRepository  message.Repository
}

// NewUserHandler implements the UserHandler interface.
func NewUserHandler(
	userRepository user.Repository,
	msgRepository message.Repository,
) UserHandler {
	return &userHandler{
		userRepository: userRepository,
		msgRepository:  msgRepository,
	}
}

// Execute performs user status handling as: register in the database and notifies contacts.
func (h *userHandler) Execute(
	ctx context.Context,
	usr user.User,
	chMessage chan<- message.Message,
) error {
	if err := h.setUserPresence(ctx, usr.ID, usr.Status, usr.ServerID); err != nil {
		return err
	}

	if usr.Status == user.Online.String() {
		if err := h.sendMessagesOffline(ctx, usr.ID, chMessage); err != nil {
			return err
		}

		if err := h.notifyPresenceOfContactsToUser(ctx, usr.ID, chMessage); err != nil {
			return err
		}
	}

	if err := h.notifyUserPresenceToContacts(ctx, usr.ID, usr.Status, chMessage); err != nil {
		return err
	}

	return nil
}

// setUserPresence logs user status in the data store.
func (h *userHandler) setUserPresence(
	ctx context.Context,
	userID string,
	userStatus string,
	serverID string,
) error {
	if userStatus == user.Online.String() {
		return h.userRepository.AddUserPresence(ctx, userID, serverID, time.Now().UTC())
	}
	return h.userRepository.RemoveUserPresence(ctx, userID)
}

// sendMessagesOffline send offline messages to user.
func (h *userHandler) sendMessagesOffline(
	ctx context.Context,
	userID string,
	chMessage chan<- message.Message,
) error {
	messages, err := h.msgRepository.GetAllMessages(ctx, userID)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		chMessage <- *msg
	}

	return h.msgRepository.DeleteAllMessages(ctx, userID)
}

// notifyPresenceOfContactsToUser sends presence message from contacts to user.
func (h *userHandler) notifyPresenceOfContactsToUser(
	ctx context.Context,
	userID string,
	chMessage chan<- message.Message,
) error {
	contacts, err := h.userRepository.GetAllContactsOnline(ctx, userID)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		msg, _ := message.New(contact, userID, "", message.ContentTypeStatus,
			user.Online.String())
		chMessage <- *msg
	}

	return nil
}

// notifyUserPresenceToContacts sends user presence message to online contacts.
func (h *userHandler) notifyUserPresenceToContacts(
	ctx context.Context,
	userID string,
	userStatus string,
	chMessage chan<- message.Message,
) error {
	contacts, err := h.userRepository.GetAllRelationshipsOnline(ctx, userID)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		msg, _ := message.New(userID, contact, "", message.ContentTypeStatus, userStatus)
		chMessage <- *msg
	}

	return nil
}
