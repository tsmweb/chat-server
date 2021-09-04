package contact

const (
	// Online represents the contact online
	Online = iota

	// Offline represents the contact offline
	Offline

	// NotFound represents the contact not found
	NotFound
)

var presenceTypeText = map[PresenceType]string{
	Online:   "online",
	Offline:  "offline",
	NotFound: "not_found",
}

// PresenceType represents the presence type ("online", "offline").
type PresenceType int

// String return the name of the PresenceType.
func (p PresenceType) String() string {
	return presenceTypeText[p]
}

// PresenceTypeText return the name of the PresenceType.
func PresenceTypeText(code PresenceType) string {
	return code.String()
}
