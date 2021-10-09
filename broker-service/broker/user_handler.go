package broker

import (
	"context"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
	"time"
)

// UserHandler handles user.
type UserHandler interface {
	// Execute performs user handling.
	Execute(ctx context.Context, usr user.User, chMessage chan<- message.Message) *ErrorEvent
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
func (h *userHandler) Execute(ctx context.Context, usr user.User, chMessage chan<- message.Message) *ErrorEvent {
	if err := h.setUserStatus(ctx, usr); err != nil {
		return NewErrorEvent(usr.ID, "UserHandler.setUserStatus()", err.Error())
	}

	if usr.Status == user.Online.String() {
		if err := h.sendMessagesOffline(ctx, usr, chMessage); err != nil {
			return NewErrorEvent(usr.ID, "UserHandler.sendMessagesOffline()", err.Error())
		}

		if err := h.notifyPresenceOfContactsToUser(ctx, usr, chMessage); err != nil {
			return NewErrorEvent(usr.ID, "UserHandler.notifyPresenceOfContactsToUser()", err.Error())
		}
	}

	if err := h.notifyUserPresenceToContacts(ctx, usr, chMessage); err != nil {
		return NewErrorEvent(usr.ID, "UserHandler.notifyUserPresenceToContacts()", err.Error())
	}

	return nil
}

// setUserStatus logs user status in the data store.
func (h *userHandler) setUserStatus(ctx context.Context, usr user.User) error {
	if usr.Status == user.Online.String() {
		return h.userRepository.AddUserPresence(ctx, usr.ID, usr.ServerID, time.Now().UTC())
	}

	return h.userRepository.RemoveUserPresence(ctx, usr.ID)
}

// sendMessagesOffline send offline messages to user.
func (h *userHandler) sendMessagesOffline(ctx context.Context, usr user.User,
	chMessage chan<- message.Message) error {
	messages, err := h.msgRepository.GetAllMessages(ctx, usr.ID)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		chMessage <- *msg
	}

	return h.msgRepository.DeleteAllMessages(ctx, usr.ID)
}

// notifyPresenceOfContactsToUser sends presence message from contacts to user.
func (h *userHandler) notifyPresenceOfContactsToUser(ctx context.Context, usr user.User,
	chMessage chan<- message.Message) error {
	contacts, err := h.userRepository.GetAllContactsOnline(ctx, usr.ID)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		msg, _ := message.New(contact, usr.ID, "", message.ContentTypeStatus, user.Online.String())
		chMessage <- *msg
	}

	return nil
}

// notifyUserPresenceToContacts sends user presence message to online contacts.
func (h *userHandler) notifyUserPresenceToContacts(ctx context.Context, usr user.User,
	chMessage chan<- message.Message) error {
	contacts, err := h.userRepository.GetAllRelationshipsOnline(ctx, usr.ID)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		msg, _ := message.New(usr.ID, contact, "", message.ContentTypeStatus, usr.Status)
		chMessage <- *msg
	}

	return nil
}
