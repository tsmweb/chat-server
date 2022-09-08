package service

import (
	"fmt"
	"os"

	"github.com/tsmweb/go-helper-api/observability/event"
	"github.com/tsmweb/user-service/config"
)

func Info(id, title, detail string) {
	e := event.New(
		config.HostID(),
		id,
		title,
		event.Info,
		detail,
	)
	send(e)
}

func Log(id, title, detail string) {
	e := event.New(
		config.HostID(),
		id,
		title,
		event.Debug,
		detail,
	)
	send(e)
}

func Error(id, title string, err error) {
	e := event.New(
		config.HostID(),
		id,
		title,
		event.Error,
		err.Error(),
	)
	send(e)
}

func Warn(id, title, detail string) {
	e := event.New(
		config.HostID(),
		id,
		title,
		event.Warning,
		detail,
	)
	send(e)
}

func send(e *event.Event) {
	if err := event.Send(e); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[!] Error sending event: %v\n", err)
	}
}
