package adapter

import (
	"time"

	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/common/service"
	"github.com/tsmweb/broker-service/infra/protobuf"
	"google.golang.org/protobuf/proto"
)

// UserEventMarshal is a user.Event encoder for protobuf.ContactEvent.
func UserEventMarshal(e *user.Event) ([]byte, error) {
	epb := protobufFromUserEvent(e)
	b, err := proto.Marshal(epb)
	if err != nil {
		return nil, service.FormatError("adapter::UserEventMarshal", err)
	}
	return b, nil
}

// UserEventUnmarshal is a protobuf.ContactEvent decoder for user.Event.
func UserEventUnmarshal(in []byte, e *user.Event) error {
	epb := new(protobuf.ContactEvent)
	if err := proto.Unmarshal(in, epb); err != nil {
		return service.FormatError("adapter::UserEventUnmarshal", err)
	}
	protobufToUserEvent(epb, e)
	return nil
}

func protobufFromUserEvent(e *user.Event) *protobuf.ContactEvent {
	return &protobuf.ContactEvent{
		UserId:    e.UserID,
		ContactId: e.ContactID,
		Event:     protobuf.ContactEventType(protobuf.ContactEventType_value[e.Event]),
		EventDate: e.EventDate.Unix(),
	}
}

func protobufToUserEvent(epb *protobuf.ContactEvent, e *user.Event) {
	e.UserID = epb.GetUserId()
	e.ContactID = epb.GetContactId()
	e.Event = epb.GetEvent().String()
	e.EventDate = time.Unix(epb.GetEventDate(), 0)
}
