package message

// Encoder is a Message encoder for byte slice.
type Encoder interface {
	Marshal(m *Message) ([]byte, error)
}

// The EncoderFunc type is an adapter to allow the use of ordinary functions as encoders of Message for byte slice.
// If f is a function with the appropriate signature, EncoderFunc(f) is a Encoder that calls f.
type EncoderFunc func(m *Message) ([]byte, error)

// Marshal calls f(m).
func (f EncoderFunc) Marshal(m *Message) ([]byte, error) {
	return f(m)
}
