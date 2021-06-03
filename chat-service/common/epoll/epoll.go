package epoll

import (
	"net"
)

// EPoll describes an object that implements logic of polling connections for
// i/o events such as availability of read() operations.
type EPoll interface {
	// ObservableRead returns an implementation of the Observer interface to observe
	// the read events of the connection.
	ObservableRead(conn net.Conn) (Observer, error)
}

// Observer describes an object that implements methods to start and stop watching events.
type Observer interface {
	// Start adds the callback function in the observation list.
	Start(callbackFn func(closed bool, err error)) error

	// Stop removes from the observation list.
	Stop()
}
