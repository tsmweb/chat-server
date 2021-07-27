package adapter

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"net"
)

// ReaderWS is a net.Conn websocket reader.
func ReaderWS(conn net.Conn) (io.Reader, error) {
	h, r, err := wsutil.NextReader(conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}
	if h.OpCode.IsControl() {
		return nil, wsutil.ControlFrameHandler(conn, ws.StateServerSide)(h, r)
	}

	return r, nil
}

// WriterWS is a net.Conn websocket writer.
func WriterWS(conn net.Conn, data interface{}) error {
	w := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return w.Flush()
}
