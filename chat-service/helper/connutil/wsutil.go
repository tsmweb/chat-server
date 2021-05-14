package connutil

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
)

func Read(conn io.ReadWriter) (io.Reader, error) {
	h, r, err := wsutil.NextReader(conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}
	if h.OpCode.IsControl() {
		return nil, wsutil.ControlFrameHandler(conn, ws.StateServerSide)(h, r)
	}

	return r, nil
}

func Write(conn io.Writer, x interface{}) error {
	w := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(x); err != nil {
		return err
	}

	return w.Flush()
}