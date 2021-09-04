package epoll

import (
	"fmt"
	"github.com/tsmweb/easygo/netpoll"
	"net"
)

// epoll is an adapter for the netpoll package that implements the EPoll interface.
type epoll struct {
	poller netpoll.Poller
}

func NewEPoll(poller netpoll.Poller) EPoll {
	return &epoll{
		poller: poller,
	}
}

// ObservableRead returns an implementation of the Observer interface to observe
// the read events of the connection.
func (e *epoll) ObservableRead(conn net.Conn) (Observer, error) {
	desc, err := netpoll.HandleRead(conn)
	if err != nil {
		return nil, err
	}

	r := &reader{
		poller: e.poller,
		desc:   desc,
	}

	return r, nil
}

// reader implements the Observer interface to observe reading events.
type reader struct {
	poller netpoll.Poller
	desc   *netpoll.Desc
}

// Start adds the callback function in the observation list.
func (r *reader) Start(callbackFn func(closed bool, err error)) error {
	return r.poller.Start(r.desc, func(event netpoll.Event) {
		if event&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
			callbackFn(true, nil)
			return
		}
		if event&netpoll.EventErr != 0 {
			callbackFn(true, fmt.Errorf("poller event error"))
			return
		}
		if event&netpoll.EventPollerClosed != 0 {
			callbackFn(true, netpoll.ErrClosed)
			return
		}

		callbackFn(false, nil)
	})
}

// Stop removes from the observation list.
func (r *reader) Stop() {
	r.poller.Stop(r.desc)
}

// Poller OnWaitError will be called from goroutine, waiting for events.
func ProviderPollerConfig(callbackFn func(err error)) *netpoll.Config {
	return &netpoll.Config{
		OnWaitError: callbackFn,
	}
}
