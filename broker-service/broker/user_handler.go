package broker

import (
	"context"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
	"time"
)

// HandleUser handles user.
type HandleUser interface {
	// Execute performs user handling.
	Execute(ctx context.Context, usr user.User, chMessage chan<- message.Message) *ErrorEvent
}

type handleUser struct {
	userRepository user.Repository
	msgRepository  message.Repository
}

// NewHandleUser implements the HandleUser interface.
func NewHandleUser(
	userRepository user.Repository,
	msgRepository message.Repository,
) HandleUser {
	return &handleUser{
		userRepository: userRepository,
		msgRepository:  msgRepository,
	}
}

// Execute performs user status handling as: register in the database,
// publish in topic kafka and notifies contacts.
func (h *handleUser) Execute(ctx context.Context, usr user.User, chMessage chan<- message.Message) *ErrorEvent {
	if err := h.setUserStatus(ctx, usr); err != nil {
		return NewErrorEvent(usr.ID, "HandleUser.setUserStatus()", err.Error())
	}

	if usr.Status == user.Online.String() {
		if err := h.sendMessagesOffline(ctx, usr, chMessage); err != nil {
			return NewErrorEvent(usr.ID, "HandleUser.sendMessagesOffline()", err.Error())
		}

		if err := h.notifyPresenceOfContactsToUser(ctx, usr, chMessage); err != nil {
			return NewErrorEvent(usr.ID, "HandleUser.notifyPresenceOfContactsToUser()", err.Error())
		}
	}

	if err := h.notifyUserPresenceToContacts(ctx, usr, chMessage); err != nil {
		return NewErrorEvent(usr.ID, "HandleUser.notifyUserPresenceToContacts()", err.Error())
	}

	return nil
}

// setUserStatus logs user status in the data store.
func (h *handleUser) setUserStatus(ctx context.Context, usr user.User) error {
	if usr.Status == user.Online.String() {
		return h.userRepository.AddUser(ctx, usr.ID, usr.ServerID, time.Now().UTC())
	}

	return h.userRepository.DeleteUser(ctx, usr.ID, usr.ServerID)
}

// sendMessagesOffline send offline messages to user.
func (h *handleUser) sendMessagesOffline(ctx context.Context, usr user.User,
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
func (h *handleUser) notifyPresenceOfContactsToUser(ctx context.Context, usr user.User,
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
func (h *handleUser) notifyUserPresenceToContacts(ctx context.Context, usr user.User,
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
