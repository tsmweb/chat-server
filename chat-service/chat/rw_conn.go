package chat

import (
	"io"
	"net"
)

// ConnReader is a net.Conn reader.
type ConnReader interface {
	Reader(conn net.Conn) (io.Reader, error)
}

// The ConnReaderFunc type is an adapter to allow the use of ordinary functions as readers of net.Conn.
// If f is a function with the appropriate signature, ConnReaderFunc(f) is a ConnReader that calls f.
type ConnReaderFunc func(conn net.Conn) (io.Reader, error)

// Reader calls f(conn).
func (f ConnReaderFunc) Reader(conn net.Conn) (io.Reader, error) {
	return f(conn)
}

// ConnWriter is a net.Conn writer.
type ConnWriter interface {
	Writer(conn net.Conn, data interface{}) error
}

// The ConnWriterFunc type is an adapter to allow the use of ordinary functions as writers of net.Conn.
// If f is a function with the appropriate signature, ConnWriterFunc(f) is a ConnWriter that calls f.
type ConnWriterFunc func(conn net.Conn, data interface{}) error

// Writer calls f(conn, data).
func (f ConnWriterFunc) Writer(conn net.Conn, data interface{}) error {
	return f(conn, data)
}
