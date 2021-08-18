package user

import "time"

// Status type that represents the user's status as UserOnline and UserOffline.
type Status int

const (
	Online  Status = 0x1
	Offline Status = 0x2
)

func (s Status) String() (str string) {
	name := func(status Status, name string) bool {
		if s&status == 0 {
			return false
		}
		str = name
		return true
	}

	if name(Online, "online") {
		return
	}
	if name(Offline, "offline") {
		return
	}

	return
}

type User struct {
	ID       string
	Status   string
	ServerID string
	Date     time.Time
}

func NewUser(id string, status Status, serverID string) *User {
	return &User{
		ID:       id,
		Status:   status.String(),
		ServerID: serverID,
		Date:     time.Now().UTC(),
	}
}
