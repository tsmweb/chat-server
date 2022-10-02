package adapter

import (
	"time"

	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/common/service"
	"github.com/tsmweb/broker-service/infra/protobuf"
	"google.golang.org/protobuf/proto"
)

// UserMarshal is a user.User encoder for protobuf.User.
func UserMarshal(u *user.User) ([]byte, error) {
	upb := protobufFromUser(u)
	b, err := proto.Marshal(upb)
	if err != nil {
		return nil, service.FormatError("adapter::UserMarshal", err)
	}
	return b, nil
}

// UserUnmarshal is a protobuf.User decoder for user.User.
func UserUnmarshal(in []byte, u *user.User) error {
	upb := new(protobuf.User)
	if err := proto.Unmarshal(in, upb); err != nil {
		return service.FormatError("adapter::UserUnmarshal", err)
	}
	protobufToUser(upb, u)
	return nil
}

func protobufFromUser(u *user.User) *protobuf.User {
	return &protobuf.User{
		Id:       u.ID,
		Status:   protobuf.UserStatus(protobuf.UserStatus_value[u.Status]),
		ServerID: u.ServerID,
		Date:     u.Date.Unix(),
	}
}

func protobufToUser(upb *protobuf.User, u *user.User) {
	u.ID = upb.GetId()
	u.Status = upb.GetStatus().String()
	u.ServerID = upb.GetServerID()
	u.Date = time.Unix(upb.GetDate(), 0)
}
