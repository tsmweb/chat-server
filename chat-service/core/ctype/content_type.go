package ctype

// ContentType represents the type of message content, such as ACK, TEXT, MEDIA, STATUS and ERROR.
type ContentType int

const (
	ACK    ContentType = 0x1
	TEXT               = 0x2
	MEDIA              = 0x4
	STATUS             = 0x8
	ERROR              = 0x80
)

func (ct ContentType) String() (str string) {
	name := func(contentType ContentType, name string) bool {
		if ct&contentType == 0 {
			return false
		}
		str = name
		return true
	}

	if name(ACK, "ACK") { return }
	if name(TEXT, "TEXT") { return }
	if name(MEDIA, "MEDIA") { return }
	if name(STATUS, "STATUS") { return }
	if name(ERROR, "ERROR") { return }

	return
}
