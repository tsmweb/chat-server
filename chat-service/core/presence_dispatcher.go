package core

import (
	"github.com/tsmweb/chat-service/common/ebus"
	"github.com/tsmweb/chat-service/core/status"
	"log"
)

type PresenceDispatcher struct {
	eBus ebus.EBus
}

func NewPresenceDispatcher() *PresenceDispatcher {
	return &PresenceDispatcher{}
}

func (pd *PresenceDispatcher) Send(userID string, status status.UserStatus) error {
	log.Printf("[>] send user presence %s: %s\n", userID, status.String())
	return nil
}