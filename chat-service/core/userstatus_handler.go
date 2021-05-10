package core

import (
	"log"
	"time"
)

type UserStatusHandler struct {
	repository         Repository
	presenceDispatcher *PresenceDispatcher
}

func NewUserStatusHandler(
	r Repository,
	pnd *PresenceDispatcher,
) *UserStatusHandler {
	return &UserStatusHandler{
		repository:         r,
		presenceDispatcher: pnd,
	}
}

func (ush *UserStatusHandler) HandleStatus(userID, host string, status UserStatus) error {
	err := ush.setStatus(userID, host, status)
	if err != nil {
		return err
	}

	return ush.presenceDispatcher.Send(userID, status)
}

func (ush *UserStatusHandler) setStatus(userID, host string, status UserStatus) error {
	log.Printf("[>] %s set status %s\n", userID, status.String())

	if status == ONLINE {
		if err := ush.repository.AddUserOnline(userID, host, time.Now().UTC()); err != nil {
			return err
		}
	} else {
		if err := ush.repository.DeleteUserOnline(userID); err != nil {
			return err
		}
	}

	return nil
}
