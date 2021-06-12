package status

// UserStatus type that represents the user's status as ONLINE and OFFLINE.
type UserStatus int

const (
	ONLINE  UserStatus = 0x1
	OFFLINE            = 0x2
)

func (us UserStatus) String() (str string) {
	name := func(status UserStatus, name string) bool {
		if us&status == 0 {
			return false
		}
		str = name
		return true
	}

	if name(ONLINE, "ONLINE") { return }
	if name(OFFLINE, "OFFLINE") { return }

	return
}
