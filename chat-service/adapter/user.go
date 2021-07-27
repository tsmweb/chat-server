package adapter

import (
	"github.com/tsmweb/chat-service/chat/user"
	"github.com/tsmweb/chat-service/infra/protobuf"
	"google.golang.org/protobuf/proto"
	"time"
)

// UserMarshal is a user.User encoder for protobuf.User.
func UserMarshal(u *user.User) ([]byte, error) {
	upb := new(protobuf.User)
	protobufFromUser(u, upb)
	return proto.Marshal(upb)
}

// UserUnmarshal is a protobuf.User decoder for user.User.
func UserUnmarshal(in []byte, u *user.User) error {
	upb := new(protobuf.User)
	if err := proto.Unmarshal(in, upb); err != nil {
		return err
	}
	protobufToUser(upb, u)
	return nil
}

func protobufFromUser(u *user.User, upb *protobuf.User) {
	upb.Id = u.ID
	upb.Status = protobuf.UserStatus(protobuf.UserStatus_value[u.Status])
	upb.ServerID = u.ServerID
	upb.Date = u.Date.Unix()
}

func protobufToUser(upb *protobuf.User, u *user.User) {
	u.ID = upb.GetId()
	u.Status = upb.GetStatus().String()
	u.ServerID = upb.GetServerID()
	u.Date = time.Unix(upb.GetDate(), 0)
}
