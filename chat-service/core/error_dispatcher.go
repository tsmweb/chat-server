package core

import (
	"github.com/tsmweb/chat-service/common/ebus"
	"log"
)

type ErrorDispatcher struct {
	eBus ebus.EBus
}

func NewErrorDispatcher() *ErrorDispatcher {
	return &ErrorDispatcher{}
}

func (ed *ErrorDispatcher) Send(err error) {
	log.Printf("[!] send error: %s \n", err.Error())
}
