package core

import (
	"log"
)

type MessageDispatcher interface {
	Send(msg *Message) error
}

type OfflineMessageDispatcher struct {
}

func NewOfflineMessageDispatcher() *OfflineMessageDispatcher {
	return &OfflineMessageDispatcher{}
}

func (omd *OfflineMessageDispatcher) Send(msg *Message) error {
	log.Printf("[>] %s sending an offline message to %s\n", msg.From, msg.To)
	return nil
}

type GroupMessageDispatcher struct {
}

func NewGroupMessageDispatcher() *GroupMessageDispatcher {
	return &GroupMessageDispatcher{}
}

func (gmd *GroupMessageDispatcher) Send(msg *Message) error {
	log.Printf("[>] %s sending message to the group %s\n", msg.From, msg.Group)
	return nil
}
