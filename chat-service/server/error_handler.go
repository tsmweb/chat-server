package server

import (
	"context"
	"github.com/tsmweb/go-helper-api/kafka"
	"log"
)

// HandleError handles errors.
type HandleError interface {
	// Execute performs errors handling.
	Execute(ctx context.Context, err ErrorEvent)

	// Close connections.
	Close()
}

type handleError struct {
	encoder  ErrorEventEncoder
	producer kafka.Producer
}

// NewHandleError implements the HandleError interface.
func NewHandleError(
	encoder ErrorEventEncoder,
	producer kafka.Producer,
) HandleError {
	return &handleError{
		encoder:  encoder,
		producer: producer,
	}
}

// Execute performs errors handling and publish in topic kafka.
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

// Close connection with kafka userProducer.
func (h *handleError) Close() {
	h.producer.Close()
}
