package server

import (
	"context"
	"github.com/tsmweb/go-helper-api/kafka"
	"log"
)

type HandleError interface {
	Execute(ctx context.Context, err ErrorEvent)
	Stop()
}

type handleError struct {
	encoder ErrorEventEncoder
	producer kafka.Producer
}

func NewHandleError(
	encoder ErrorEventEncoder,
	producer kafka.Producer,
) HandleError {
	return &handleError{
		encoder:  encoder,
		producer: producer,
	}
}

func (h *handleError) Execute(ctx context.Context, errEvent ErrorEvent) {
	epb, err := h.encoder.Marshal(&errEvent)
	if err != nil {
		log.Printf("[!] HandleError.Execute() \n Error: %v \n Data: %v", err.Error(), errEvent)
		return
	}

	if err := h.producer.Publish(ctx, []byte(errEvent.HostID), epb); err != nil {
		log.Printf("[!] HandleError.Execute() \n Error: %v \n Data: %v", err.Error(), errEvent)
	}
}

func (h *handleError) Stop() {
	h.producer.Close()
}