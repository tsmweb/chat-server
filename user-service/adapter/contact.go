package adapter

import (
	"github.com/tsmweb/user-service/app/contact"
	"github.com/tsmweb/user-service/infra/protobuf"
	"google.golang.org/protobuf/proto"
	"time"
)

// ContactEventMarshal is a contact.Event encoder for protobuf.ContactEvent.
func ContactEventMarshal(e *contact.Event) ([]byte, error) {
	epb := protobufFromContactEvent(e)
	return proto.Marshal(epb)
}

// ContactEventUnmarshal is a protobuf.ContactEvent decoder for contact.Event.
func ContactEventUnmarshal(in []byte, e *contact.Event) error {
	epb := new(protobuf.ContactEvent)
	if err := proto.Unmarshal(in, epb); err != nil {
		return err
	}
	protobufToContactEvent(epb, e)
	return nil
}

func protobufFromContactEvent(e *contact.Event) *protobuf.ContactEvent {
	return &protobuf.ContactEvent{
		UserId:    e.UserID,
		ContactId: e.ContactID,
		Event:     protobuf.ContactEventType(protobuf.ContactEventType_value[e.Event]),
		EventDate: e.EventDate.Unix(),
	}
}

func protobufToContactEvent(epb *protobuf.ContactEvent, e *contact.Event) {
	e.UserID = epb.GetUserId()
	e.ContactID = epb.GetContactId()
	e.Event = epb.GetEvent().String()
	e.EventDate = time.Unix(epb.GetEventDate(), 0)
}
