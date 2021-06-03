//go:generate wire
//+build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/tsmweb/chat-service/api"
	"github.com/tsmweb/chat-service/common/concurrent"
	"github.com/tsmweb/chat-service/common/connutil"
	"github.com/tsmweb/chat-service/common/epoll"
	"github.com/tsmweb/chat-service/common/setting"
	"github.com/tsmweb/chat-service/core"
	"github.com/tsmweb/chat-service/infrastructure/db"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/middleware"
)

func InitChat(
	localhost string,
	executor concurrent.ExecutorService,
) (*api.Router, error) {
	wire.Build(
		ProviderJWT,
		middleware.NewAuth,
		ProviderPollerConfig,
		netpoll.New,
		epoll.NewEPoll,
		ProviderConnReader,
		ProviderConnWriter,
		core.Inject,
		api.Inject,
	)

	return nil, nil
}

/*
 * PROVIDERS
 */

func ProviderConnReader() connutil.Reader {
	return connutil.FuncReader(connutil.ReaderWS)
}

func ProviderConnWriter() connutil.Writer {
	return connutil.FuncWriter(connutil.WriterWS)
}

// Data Base
var databaseInstance db.Database

func ProviderDataBase() db.Database {
	if databaseInstance == nil {
		databaseInstance = db.NewPostgresDatabase()
	}

	return databaseInstance
}

// Authentication MockJWT
var jwtInstance auth.JWT

func ProviderJWT() auth.JWT {
	if jwtInstance == nil {
		jwtInstance = auth.NewJWT(setting.PathPrivateKey(), setting.PathPublicKey())
	}

	return jwtInstance
}

// Poller OnWaitError will be called from goroutine, waiting for events.
func ProviderPollerConfig(executor concurrent.ExecutorService, dispatcher *core.ErrorDispatcher) *netpoll.Config {
	return &netpoll.Config{
		OnWaitError: func(err error) {
			executor.Schedule(func(ctx context.Context) {
				dispatcher.Send(err)
			})
		},
	}
}
