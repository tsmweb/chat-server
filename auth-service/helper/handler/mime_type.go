package handler

const (
	// MimeApplicationJSON represents the MimeType "application/json".
	MimeApplicationJSON = iota

	// MimeApplicationPDF represents the MimeType "application/pdf".
	MimeApplicationPDF

	// MimeApplicationOctetStream represents the MimeType "application/octet-stream".
	MimeApplicationOctetStream

	// MimeImageJPEG represents the MimeType "image/jpeg".
	MimeImageJPEG

	// MimeTextPlain represents the MimeType "text/plain".
	MimeTextPlain
)

var mimeTypeText = map[MimeType]string{
	MimeApplicationJSON:        "application/json",
	MimeApplicationPDF:         "application/pdf",
	MimeApplicationOctetStream: "application/octet-stream",
	MimeImageJPEG:              "image/jpeg",
	MimeTextPlain:              "text/plain",
}

// MimeType represents the mime type ("application/json", "text/plain", "image/jpeg", ...).
type MimeType int

// String return the name of the mimeType.
func (m MimeType) String() string {
	return mimeTypeText[m]
}

// MimeTypeText return the name of the mimeType.
func MimeTypeText(code MimeType) string {
	return code.String()
}
