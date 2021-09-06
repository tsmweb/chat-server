package adapter

import (
	"github.com/tsmweb/broker-service/broker/group"
	"github.com/tsmweb/broker-service/infra/protobuf"
	"google.golang.org/protobuf/proto"
	"time"
)

// EventMarshal is a group.Event encoder for protobuf.GroupEvent.
func EventMarshal(e *group.Event) ([]byte, error) {
	epb := protobufFromEvent(e)
	return proto.Marshal(epb)
}

// EventUnmarshal is a protobuf.GroupEvent decoder for group.Event.
func EventUnmarshal(in []byte, e *group.Event) error {
	epb := new(protobuf.GroupEvent)
	if err := proto.Unmarshal(in, epb); err != nil {
		return err
	}
	protobufToEvent(epb, e)
	return nil
}

func protobufFromEvent(e *group.Event) *protobuf.GroupEvent {
	return &protobuf.GroupEvent{
		GroupId:   e.GroupID,
		MemberId:  e.MemberID,
		Event:     protobuf.EventType(protobuf.EventType_value[e.Event]),
		EventDate: e.EventDate.Unix(),
	}
}

func protobufToEvent(epb *protobuf.GroupEvent, e *group.Event) {
	e.GroupID = epb.GetGroupId()
	e.MemberID = epb.GetMemberId()
	e.Event = epb.GetEvent().String()
	e.EventDate = time.Unix(epb.GetEventDate(), 0)
}
