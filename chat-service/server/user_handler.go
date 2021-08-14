package server

import (
	"context"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/chat-service/server/user"
	"github.com/tsmweb/go-helper-api/kafka"
	"time"
)

// HandleUserStatus handles user status.
type HandleUserStatus interface {
	// Execute performs user status handling.
	Execute(ctx context.Context, userID string, status user.Status,
		chMessage, chSendMessage chan<- message.Message) *ErrorEvent

	// Close connections.
	Close()
}

type handleUserStatus struct {
	encoder    user.Encoder
	producer   kafka.Producer
	repository Repository
}

// NewHandleUserStatus implements the HandleUserStatus interface.
func NewHandleUserStatus(
	encoder user.Encoder,
	producer kafka.Producer,
	repository Repository,
) HandleUserStatus {
	return &handleUserStatus{
		encoder:    encoder,
		producer:   producer,
		repository: repository,
	}
}

// Execute performs user status handling as: register in the database,
// publish in topic kafka and notifies contacts.
func (h *handleUserStatus) Execute(ctx context.Context, userID string, status user.Status,
	chMessage, chSendMessage chan<- message.Message,
) *ErrorEvent {
	serverID := "OFF"
	if status == user.Online {
		serverID = config.HostID()
	}

	if err := h.setUserStatus(userID, serverID, status); err != nil {
		return err
	}

	if err := h.publishUserStatus(ctx, userID, serverID, status); err != nil {
		return err
	}

	if status == user.Online {
		if err := h.notifyContactStatusToUser(userID, status, chSendMessage); err != nil {
			return err
		}
	}

	if err := h.notifyUserStatusToContacts(userID, status, chMessage); err != nil {
		return err
	}

	return nil
}

// Close connection with kafka producer.
func (h *handleUserStatus) Close() {
	h.producer.Close()
}

// setUserStatus logs user status in the data store.
func (h *handleUserStatus) setUserStatus(userID, serverID string, status user.Status) *ErrorEvent {
	if status == user.Online {
		if err := h.repository.AddUserOnline(userID, serverID, time.Now().UTC()); err != nil {
			return NewErrorEvent(userID, "HandleUserStatus.setUserStatus()", err.Error())
		}
	} else {
		if err := h.repository.DeleteUserOnline(userID); err != nil {
			return NewErrorEvent(userID, "HandleUserStatus.setUserStatus()", err.Error())
		}
	}

	return nil
}

// publishUserStatus publish user status to kafka topic.
func (h *handleUserStatus) publishUserStatus(ctx context.Context, userID, serverID string,
	status user.Status,
) *ErrorEvent {
	u := user.NewUser(userID, status, serverID)
	upb, err := h.encoder.Marshal(u)
	if err != nil {
		return NewErrorEvent(userID, "HandleUserStatus.publishUserStatus()", err.Error())
	}

	if err = h.producer.Publish(ctx, []byte(userID), upb); err != nil {
		return NewErrorEvent(userID, "HandleUserStatus.publishUserStatus()", err.Error())
	}

	return nil
}

// notifyContactStatusToUser sends the contact's presence message to the user.
func (h *handleUserStatus) notifyContactStatusToUser(userID string, status user.Status,
	chSendMessage chan<- message.Message) *ErrorEvent {

	contacts, err := h.repository.GetUserContactsOnline(userID)
	if err != nil {
		return NewErrorEvent(userID, "Repository.GetUserContactsOnline()", err.Error())
	}

	for _, contact := range contacts {
		msgForUser, _ := message.NewMessage(contact, userID, "", message.ContentTypeStatus, status.String())
		chSendMessage <- *msgForUser
	}

	return nil
}

// notifyUserStatusToContacts sends user presence message to online contacts.
func (h *handleUserStatus) notifyUserStatusToContacts(userID string, status user.Status,
	chMessage chan<- message.Message) *ErrorEvent {

	contacts, err := h.repository.GetContactsWithUserOnline(userID)
	if err != nil {
		return NewErrorEvent(userID, "Repository.GetContactsWithUserOnline()", err.Error())
	}

	for _, contact := range contacts {
		msgForContact, _ := message.NewMessage(userID, contact, "", message.ContentTypeStatus, status.String())
		chMessage <- *msgForContact
	}

	return nil
}
