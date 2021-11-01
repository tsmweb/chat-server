package adapter

import (
	"github.com/tsmweb/user-service/group"
	"github.com/tsmweb/user-service/infra/protobuf"
	"google.golang.org/protobuf/proto"
	"time"
)

// GroupEventMarshal is a group.Event encoder for protobuf.GroupEvent.
func GroupEventMarshal(e *group.Event) ([]byte, error) {
	epb := protobufFromGroupEvent(e)
	return proto.Marshal(epb)
}

// GroupEventUnmarshal is a protobuf.GroupEvent decoder for group.Event.
func GroupEventUnmarshal(in []byte, e *group.Event) error {
	epb := new(protobuf.GroupEvent)
	if err := proto.Unmarshal(in, epb); err != nil {
		return err
	}
	protobufToGroupEvent(epb, e)
	return nil
}

func protobufFromGroupEvent(e *group.Event) *protobuf.GroupEvent {
	return &protobuf.GroupEvent{
		GroupId:   e.GroupID,
		MemberId:  e.MemberID,
		Event:     protobuf.GroupEventType(protobuf.GroupEventType_value[e.Event]),
		EventDate: e.EventDate.Unix(),
	}
}

func protobufToGroupEvent(epb *protobuf.GroupEvent, e *group.Event) {
	e.GroupID = epb.GetGroupId()
	e.MemberID = epb.GetMemberId()
	e.Event = epb.GetEvent().String()
	e.EventDate = time.Unix(epb.GetEventDate(), 0)
}
