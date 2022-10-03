package broker

import (
	"context"

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
		return h.msgRepository.AddGroupMemberToCache(ctx, evt.GroupID, evt.MemberID)
	}

	if evt.Event == group.EventRemoveMember.String() {
		return h.msgRepository.RemoveGroupMemberFromCache(ctx, evt.GroupID, evt.MemberID)
	}

	if evt.Event == group.EventDeleteGroup.String() {
		return h.msgRepository.RemoveGroupFromCache(ctx, evt.GroupID)
	}

	return nil
}
