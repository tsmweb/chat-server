package broker

import (
	"context"
	"fmt"

	"github.com/tsmweb/broker-service/broker/group"
	"github.com/tsmweb/broker-service/broker/message"
)

// GroupEventHandler handles group events.
type GroupEventHandler interface {
	// Execute performs group event handling.
	Execute(ctx context.Context, evt group.Event) error
}

type groupEventHandler struct {
	msgRepository message.Repository
}

// NewGroupEventHandler implements the GroupEventHandler interface.
func NewGroupEventHandler(msgRepository message.Repository) GroupEventHandler {
	return &groupEventHandler{
		msgRepository: msgRepository,
	}
}

// Execute performs group event handling.
func (h *groupEventHandler) Execute(ctx context.Context, evt group.Event) error {
	if evt.Event == group.EventAddMember.String() {
		if err := h.msgRepository.AddGroupMemberToCache(ctx, evt.GroupID, evt.MemberID); err != nil {
			return fmt.Errorf("GroupEventHandler::msgRepository::AddGroupMemberToCache. Error: %v",
				err.Error())
		}
		return nil
	}

	if evt.Event == group.EventRemoveMember.String() {
		if err := h.msgRepository.RemoveGroupMemberFromCache(ctx, evt.GroupID, evt.MemberID); err != nil {
			return fmt.Errorf("GroupEventHandler.::msgRepository::RemoveGroupMemberFromCache. Error: %v",
				err.Error())
		}
		return nil
	}

	if evt.Event == group.EventDeleteGroup.String() {
		if err := h.msgRepository.RemoveGroupFromCache(ctx, evt.GroupID); err != nil {
			return fmt.Errorf("GroupEventHandler::msgRepository::RemoveGroupFromCache. Error: %v",
				err.Error())
		}
		return nil
	}

	return nil
}
