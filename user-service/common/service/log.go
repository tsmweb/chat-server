package service

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	log.Printf("[%s] HostID: %s | UserID: %s | Title: %s | Detail: %s\n",
		strings.ToUpper(e.Type), e.Host, e.User, e.Title, e.Detail)

	if err := event.Send(e); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] sending event: %v\n", err)
	}
}
