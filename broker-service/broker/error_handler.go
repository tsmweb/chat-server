package broker

import (
	"context"
	"github.com/tsmweb/go-helper-api/kafka"
	"log"
)

// ErrorHandler handles errors.
type ErrorHandler interface {
	// Execute performs errors handling.
	Execute(ctx context.Context, err ErrorEvent)

	// Close connections.
	Close()
}

type errorHandler struct {
	encoder  ErrorEventEncoder
	producer kafka.Producer
}

// NewErrorHandler implements the ErrorHandler interface.
func NewErrorHandler(
	encoder ErrorEventEncoder,
	producer kafka.Producer,
) ErrorHandler {
	return &errorHandler{
		encoder:  encoder,
		producer: producer,
	}
}

// Execute performs errors handling and publish in topic kafka.
func (h *errorHandler) Execute(ctx context.Context, errEvent ErrorEvent) {
	epb, err := h.encoder.Marshal(&errEvent)
	if err != nil {
		log.Printf("[!] ErrorHandler.Execute() \n Error: %v \n Data: %v", err.Error(), errEvent)
		return
	}

	if err = h.producer.Publish(ctx, []byte(errEvent.HostID), epb); err != nil {
		log.Printf("[!] ErrorHandler.Execute() \n Error: %v \n Data: %v", err.Error(), errEvent)
	}
}

// Close connection with kafka producer.
func (h *errorHandler) Close() {
	h.producer.Close()
}
