package user

// Encoder is a User encoder for byte slice.
type Encoder interface {
	Marshal(u *User) ([]byte, error)
}

// The EncoderFunc type is an adapter to allow the use of ordinary functions as encoders of User for byte slice.
// If f is a function with the appropriate signature, EncoderFunc(f) is a Encoder that calls f.
type EncoderFunc func (u *User) ([]byte, error)

// Marshal calls f(m).
func (f EncoderFunc) Marshal(u *User) ([]byte, error) {
	return f(u)
}
