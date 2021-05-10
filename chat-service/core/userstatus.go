package core

const (
	ONLINE = iota
	OFFLINE
)

// UserStatus type that represents the user's status as ONLINE and OFFLINE.
type UserStatus int

func (us UserStatus) String() (str string) {
	name := func(status UserStatus, name string) {
		if us&status == 0 {
			return
		}
		str = name
	}

	name(ONLINE, "ONLINE")
	name(OFFLINE, "OFFLINE")

	return
}
