package server

import (
	"context"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/chat-service/server/user"
	"github.com/tsmweb/go-helper-api/kafka"
	"time"
)

type HandleUserStatus interface {
	Execute(ctx context.Context, userID string, status user.Status,
		chMessage, chSendMessage chan<- message.Message) *ErrorEvent
	Stop()
}

type handleUserStatus struct {
	encoder    user.Encoder
	producer   kafka.Producer
	repository Repository
}

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

	if err := h.notifyUserStatus(userID, status, chMessage, chSendMessage); err != nil {
		return err
	}

	return nil
}

func (h *handleUserStatus) Stop() {
	h.producer.Close()
}

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

func (h *handleUserStatus) notifyUserStatus(userID string, status user.Status,
	chMessage, chSendMessage chan<- message.Message) *ErrorEvent {
	// sends the contact's presence message to the user.
	contacts, err := h.repository.GetUserContactsOnline(userID)
	if err != nil {
		return NewErrorEvent(userID, "HandleUserStatus.notifyUserStatus()", err.Error())
	}

	for _, contact := range contacts {
		msgForUser, _ := message.NewMessage(contact, userID, "", message.ContentStatus, status.String())
		chSendMessage <- *msgForUser
	}

	// sends user presence message to online contacts.
	contacts, err = h.repository.GetContactsWithUserOnline(userID)
	if err != nil {
		return NewErrorEvent(userID, "HandleUserStatus.notifyUserStatus()", err.Error())
	}

	for _, contact := range contacts {
		msgForContact, _ := message.NewMessage(userID, contact, "", message.ContentStatus, status.String())
		chMessage <- *msgForContact
	}

	return nil
}
