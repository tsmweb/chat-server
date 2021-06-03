package connutil

import "io"

type Reader interface {
	Reader(conn io.ReadWriter) (io.Reader, error)
}

type FuncReader func(conn io.ReadWriter) (io.Reader, error)

func (fr FuncReader) Reader(conn io.ReadWriter) (io.Reader, error) {
	return fr(conn)
}

type Writer interface {
	Writer(conn io.Writer, x interface{}) error
}

type FuncWriter func(conn io.Writer, x interface{}) error

func (fw FuncWriter) Writer(conn io.Writer, x interface{}) error {
	return fw(conn, x)
}

