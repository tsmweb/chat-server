package broker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/config"
	"github.com/tsmweb/go-helper-api/kafka"
)

// MessageHandler handles messages.
type MessageHandler interface {
	// Execute performs message handling.
	Execute(ctx context.Context, msg message.Message) error
}

type messageHandler struct {
	userRepository user.Repository
	msgRepository  message.Repository
	queue          kafka.Kafka
	encoder        message.Encoder
}

// NewMessageHandler implements the MessageHandler interface.
func NewMessageHandler(
	userRepository user.Repository,
	msgRepository message.Repository,
	queue kafka.Kafka,
	encoder message.Encoder,
) MessageHandler {
	return &messageHandler{
		userRepository: userRepository,
		msgRepository:  msgRepository,
		queue:          queue,
		encoder:        encoder,
	}
}

// Execute performs message handling.
func (h *messageHandler) Execute(ctx context.Context, msg message.Message) error {
	// check if it's a group message
	if msg.IsGroupMessage() {
		return h.processGroupMessage(ctx, &msg)
	}

	// Checks if the addressee is a valid user.
	ok, err := h.isValidUser(ctx, &msg)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Check if the sender has been blocked by the addressee.
	ok, err = h.isBlockedUser(ctx, &msg)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	return h.sendMessage(ctx, &msg)
}

func (h *messageHandler) processGroupMessage(ctx context.Context, msg *message.Message) error {
	members, err := h.msgRepository.GetAllGroupMembers(ctx, msg.Group)
	if err != nil {
		return fmt.Errorf("MessageHandler::processGroupMessage::msgRepository::GetAllGroupMembers. Error: %v",
			err.Error())
	}
	if len(members) < 1 {
		msgResponse := message.NewResponse(
			msg.ID,
			msg.From,
			msg.Group,
			message.ContentTypeError,
			message.ErrGroupIsInvalid.Error(),
		)
		return h.sendMessage(ctx, msgResponse)
	}
	if len(members) == 1 {
		return nil
	}

	var errEvents []string

	for _, member := range members {
		if member == msg.From { // sender
			continue
		}
		m, _ := msg.ReplicateTo(member)
		if err := h.sendMessage(ctx, m); err != nil {
			errEvents = append(errEvents, err.Error())
		}
	}

	if len(errEvents) > 0 {
		return fmt.Errorf("MessageHandler::processGroupMessage::sendMessage. Error: %v",
			errors.New(strings.Join(errEvents, "|")))
	}

	return nil
}

// isValidUser checks if the addressee is a valid user.
func (h *messageHandler) isValidUser(ctx context.Context, msg *message.Message) (bool, error) {
	ok, err := h.userRepository.IsValidUser(ctx, msg.To)
	if err != nil {
		return false, fmt.Errorf("MessageHandler::isValidUser::userRepository::IsValidUser. Error: %v",
			err.Error())
	}
	if !ok {
		msgResponse := message.NewResponse(
			msg.ID,
			msg.From,
			"",
			message.ContentTypeError,
			message.ErrMessageAddresseeIsInvalid.Error(),
		)
		return false, h.sendMessage(ctx, msgResponse)
	}

	return true, nil
}

// isBlockedUser check if the sender has been blocked by the addressee.
func (h *messageHandler) isBlockedUser(ctx context.Context, msg *message.Message) (bool, error) {
	ok, err := h.userRepository.IsBlockedUser(ctx, msg.To, msg.From)
	if err != nil {
		return false, fmt.Errorf("MessageHandler::isBlockedUser::userRepository::IsBlockedUser. Error: %v",
			err.Error())
	}
	if ok {
		msgResponse := message.NewResponse(
			msg.ID,
			msg.From,
			"",
			message.ContentTypeError,
			message.ErrMessageSendingBlocked.Error(),
		)
		return true, h.sendMessage(ctx, msgResponse)
	}

	return false, nil
}

func (h *messageHandler) sendMessage(ctx context.Context, msg *message.Message) error {
	serverID, err := h.userRepository.GetUserServer(ctx, msg.To)
	if err != nil {
		return fmt.Errorf("MessageHandler::sendMessage::userRepository::GetUserServer. Error: %v",
			err.Error())
	}

	if strings.TrimSpace(serverID) != "" { // online
		err = h.dispatchMessagesToHosts(ctx, serverID, msg)
		if err != nil {
			return fmt.Errorf("MessageHandler::sendMessage::dispatchMessagesToHosts. Error: %v",
				err.Error())
		}
	} else if msg.ContentType != message.ContentTypeStatus.String() { // offline
		err = h.dispatchToOfflineMessages(ctx, msg)
		if err != nil {
			return fmt.Errorf("MessageHandler::sendMessage::dispatchToOfflineMessages. Error: %v",
				err.Error())
		}
	}

	return nil
}

func (h *messageHandler) dispatchMessagesToHosts(
	ctx context.Context,
	serverID string,
	msg *message.Message,
) error {
	producer := h.queue.NewProducer(config.KafkaHostTopic(serverID))
	defer producer.Close()
	return h.dispatchMessages(ctx, producer, msg)
}

func (h *messageHandler) dispatchToOfflineMessages(
	ctx context.Context,
	msg *message.Message,
) error {
	producer := h.queue.NewProducer(config.KafkaOffMessagesTopic())
	defer producer.Close()
	return h.dispatchMessages(ctx, producer, msg)
}

func (h *messageHandler) dispatchMessages(
	ctx context.Context,
	producer kafka.Producer,
	msg *message.Message,
) error {
	mpb, err := h.encoder.Marshal(msg)
	if err != nil {
		return err
	}

	if err = producer.Publish(ctx, []byte(msg.ID), mpb); err != nil {
		return err
	}

	return nil
}
