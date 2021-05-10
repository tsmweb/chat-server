package core

import "log"

type ErrorDispatcher struct {
}

func NewErrorDispatcher() *ErrorDispatcher {
	return &ErrorDispatcher{}
}

func (ed *ErrorDispatcher) Send(err error) {
	log.Printf("[!] send error: %s \n", err.Error())
}
