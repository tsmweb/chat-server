package adapter

import (
	"github.com/tsmweb/chat-service/infra/protobuf"
	"github.com/tsmweb/chat-service/server"
	"google.golang.org/protobuf/proto"
	"time"
)

// ErrorEventMarshal is a server.ErrorEvent encoder for protobuf.Error.
func ErrorEventMarshal(e *server.ErrorEvent) ([]byte, error) {
	epb := protobufFromErrorEvent(e)
	return proto.Marshal(epb)
}

// ErrorEventUnmarshal is a protobuf.Error decoder for server.ErrorEvent.
func ErrorEventUnmarshal(in []byte, e *server.ErrorEvent) error {
	epb := new(protobuf.Error)
	if err := proto.Unmarshal(in, epb); err != nil {
		return err
	}
	protobufToEventError(epb, e)
	return nil
}

func protobufFromErrorEvent(e *server.ErrorEvent) *protobuf.Error {
	return &protobuf.Error{
		HostID: e.HostID,
		UserID: e.UserID,
		Title:  e.Title,
		Detail: e.Detail,
		Date:   e.Timestamp.Unix(),
	}
}

func protobufToEventError(epb *protobuf.Error, e *server.ErrorEvent) {
	e.HostID = epb.GetHostID()
	e.UserID = epb.GetUserID()
	e.Title = epb.GetTitle()
	e.Detail = epb.GetDetail()
	e.Timestamp = time.Unix(epb.GetDate(), 0)
}
