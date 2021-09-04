package message

// Decoder is a byte slice decoder for Message.
type Decoder interface {
	Unmarshal(in []byte, m *Message) error
}

// The DecoderFunc type is an adapter to allow the use of ordinary functions as decoders of byte slice for Message.
// If f is a function with the appropriate signature, DecoderFunc(f) is a Decoder that calls f.
type DecoderFunc func(in []byte, m *Message) error

// Unmarshal calls f(in, m).
func (f DecoderFunc) Unmarshal(in []byte, m *Message) error {
	return f(in, m)
}
