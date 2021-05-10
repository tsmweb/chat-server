package core

import (
	"log"
)

type PresenceDispatcher struct {
}

func NewPresenceDispatcher() *PresenceDispatcher {
	return &PresenceDispatcher{}
}

func (pd *PresenceDispatcher) Send(userID string, status UserStatus) error {
	log.Printf("[>] send user presence %s: %s\n", userID, status.String())
	return nil
}