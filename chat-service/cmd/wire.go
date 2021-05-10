//go:generate wire
//+build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/tsmweb/chat-service/api"
	"github.com/tsmweb/chat-service/core"
	"github.com/tsmweb/chat-service/helper/database"
	"github.com/tsmweb/chat-service/helper/setting"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"github.com/tsmweb/go-helper-api/middleware"
)

func InitChat(
	localhost string,
	exe *executor.Executor,
) (*api.Router, error) {
	wire.Build(
		ProviderJWT,
		middleware.NewAuth,
		ProviderPollerConfig,
		netpoll.New,
		core.Inject,
		api.Inject,
	)

	return nil, nil
}

/*
 * PROVIDERS
 */

// Data Base
var databaseInstance database.Database

func ProviderDataBase() database.Database {
	if databaseInstance == nil {
		databaseInstance = database.NewPostgresDatabase()
	}

	return databaseInstance
}

// Authentication JWT
var jwtInstance auth.JWT

func ProviderJWT() auth.JWT {
	if jwtInstance == nil {
		jwtInstance = auth.NewJWT(setting.PathPrivateKey(), setting.PathPublicKey())
	}

	return jwtInstance
}

// Poller OnWaitError will be called from goroutine, waiting for events.
func ProviderPollerConfig(exe *executor.Executor, dispatcher *core.ErrorDispatcher) *netpoll.Config {
	return &netpoll.Config{
		OnWaitError: func(err error) {
			exe.Schedule(func(ctx context.Context) {
				dispatcher.Send(err)
			})
		},
	}
}
