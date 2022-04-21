package adapter

import (
	"github.com/tsmweb/broker-service/broker"
	"github.com/tsmweb/broker-service/infrastructure/protobuf"
	"google.golang.org/protobuf/proto"
	"time"
)

// ErrorEventMarshal is a broker.ErrorEvent encoder for protobuf.Error.
func ErrorEventMarshal(e *broker.ErrorEvent) ([]byte, error) {
	epb := protobufFromErrorEvent(e)
	return proto.Marshal(epb)
}

// ErrorEventUnmarshal is a protobuf.Error decoder for broker.ErrorEvent.
func ErrorEventUnmarshal(in []byte, e *broker.ErrorEvent) error {
	epb := new(protobuf.Error)
	if err := proto.Unmarshal(in, epb); err != nil {
		return err
	}
	protobufToEventError(epb, e)
	return nil
}

func protobufFromErrorEvent(e *broker.ErrorEvent) *protobuf.Error {
	return &protobuf.Error{
		HostID: e.HostID,
		UserID: e.UserID,
		Title:  e.Title,
		Detail: e.Detail,
		Date:   e.Timestamp.Unix(),
	}
}

func protobufToEventError(epb *protobuf.Error, e *broker.ErrorEvent) {
	e.HostID = epb.GetHostID()
	e.UserID = epb.GetUserID()
	e.Title = epb.GetTitle()
	e.Detail = epb.GetDetail()
	e.Timestamp = time.Unix(epb.GetDate(), 0)
}
