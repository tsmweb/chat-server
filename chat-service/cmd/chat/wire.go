//go:generate wire
//+build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/tsmweb/chat-service/chat"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/infrastructure/db"
	"github.com/tsmweb/chat-service/infrastructure/kafka"
	"github.com/tsmweb/chat-service/pkg/concurrent"
	"github.com/tsmweb/chat-service/pkg/ebus"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/pkg/topic"
	"github.com/tsmweb/chat-service/util/connutil"
	"github.com/tsmweb/chat-service/web/api"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/middleware"
)

func InitChatRouter(
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
		ProviderKafka,
		chat.NewMemoryRepository,
		chat.NewChat,
		api.NewController,
		api.NewRouter,
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
		jwtInstance = auth.NewJWT(config.PathPrivateKey(), config.PathPublicKey())
	}

	return jwtInstance
}

// Poller OnWaitError will be called from goroutine, waiting for events.
func ProviderPollerConfig(executor concurrent.ExecutorService) *netpoll.Config {
	return &netpoll.Config{
		OnWaitError: func(err error) {
			executor.Schedule(func(ctx context.Context) {
				ebus.Instance().Publish(topic.ErrorMessage, err.Error())
			})
		},
	}
}

// Kafka
func ProviderKafka() chat.Kafka {
	return kafka.New([]string{config.KafkaBootstrapServers()}, config.KafkaClientID())
}
