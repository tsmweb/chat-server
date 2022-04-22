package user

// Decoder is a byte slice decoder for User.
type Decoder interface {
	Unmarshal(in []byte, u *User) error
}

// The DecoderFunc type is an adapter to allow the use of ordinary functions as decoders
// of byte slice for User.
// If f is a function with the appropriate signature, DecoderFunc(f) is a Decoder that calls f.
type DecoderFunc func(in []byte, u *User) error

// Unmarshal calls f(in, m).
func (f DecoderFunc) Unmarshal(in []byte, u *User) error {
	return f(in, u)
}
