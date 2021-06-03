package core

import (
	"github.com/tsmweb/chat-service/common/ebus"
	"log"
)

type MessageDispatcher interface {
	Send(msg *Message) error
}

type OfflineMessageDispatcher struct {
	eBus ebus.EBus
}

func NewOfflineMessageDispatcher() *OfflineMessageDispatcher {
	return &OfflineMessageDispatcher{}
}

func (omd *OfflineMessageDispatcher) Send(msg *Message) error {
	log.Printf("[>] %s sending an offline message to %s\n", msg.From, msg.To)
	return nil
}

type GroupMessageDispatcher struct {
	eBus ebus.EBus
}

func NewGroupMessageDispatcher() *GroupMessageDispatcher {
	return &GroupMessageDispatcher{}
}

func (gmd *GroupMessageDispatcher) Send(msg *Message) error {
	log.Printf("[>] %s sending message to the group %s\n", msg.From, msg.Group)
	return nil
}
