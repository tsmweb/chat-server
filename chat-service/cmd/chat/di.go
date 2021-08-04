package main

import (
	"context"
	"github.com/tsmweb/chat-service/adapter"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/infra/db"
	"github.com/tsmweb/chat-service/infra/repository"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/server"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/chat-service/server/user"
	"github.com/tsmweb/chat-service/web/api"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/kafka"
	"github.com/tsmweb/go-helper-api/middleware"
)

type Providers struct {
	ctx        context.Context
	server     *server.Server
	jwt        auth.JWT
	dataBase   db.Database
	repository server.Repository
	kafka      kafka.Kafka
}

func CreateProvider(ctx context.Context) *Providers {
	return &Providers{
		ctx: ctx,
	}
}

func (p *Providers) ServerProvider() (*server.Server, error) {
	if p.server == nil {
		epoll, err := p.EpollProvider()
		if err != nil {
			return nil, err
		}

		p.server = server.NewServer(
			p.ctx,
			epoll,
			p.ConnReaderProvider(),
			p.ConnWriterProvider(),
			p.MessageDecoderProvider(),
			p.RepositoryProvider(),
			p.KafkaProvider(),
			p.HandleMessageProvider(),
			p.HandleGroupMessageProvider(),
			p.HandleOffMessageProvider(),
			p.HandleUserStatusProvider(),
			p.HandleErrorProvider(),
		)
	}
	return p.server, nil
}

func (p *Providers) EpollProvider() (epoll.EPoll, error) {
	poller, err := netpoll.New(p.PollerConfigProvider())
	if err != nil {
		return nil, err
	}
	return epoll.NewEPoll(poller), nil
}

func (p *Providers) ConnReaderProvider() server.ConnReader {
	return server.ConnReaderFunc(adapter.ReaderWS)
}

func (p *Providers) ConnWriterProvider() server.ConnWriter {
	return server.ConnWriterFunc(adapter.WriterWS)
}

func (p *Providers) MessageEncoderProvider() message.Encoder {
	return message.EncoderFunc(adapter.MessageMarshal)
}

func (p *Providers) MessageDecoderProvider() message.Decoder {
	return message.DecoderFunc(adapter.MessageUnmarshal)
}

func (p *Providers) UserEncoderProvider() user.Encoder {
	return user.EncoderFunc(adapter.UserMarshal)
}

func (p *Providers) EventErrorEncoderProvider() server.ErrorEventEncoder {
	return server.ErrorEventEncoderFunc(adapter.ErrorEventMarshal)
}

func (p *Providers) RepositoryProvider() server.Repository {
	if p.repository == nil {
		p.repository = repository.NewMemoryRepository()
	}
	return p.repository
}

func (p *Providers) DataBaseProvider() db.Database {
	if p.dataBase == nil {
		p.dataBase = db.NewPostgresDatabase()
	}
	return p.dataBase
}

func (p *Providers) KafkaProvider() kafka.Kafka {
	if p.kafka == nil {
		p.kafka = kafka.New([]string{config.KafkaBootstrapServers()}, config.KafkaClientID())
	}
	return p.kafka
}

// Poller OnWaitError will be called from goroutine, waiting for events.
func (p *Providers) PollerConfigProvider() *netpoll.Config {
	errorProducer := p.KafkaProvider().NewProducer(config.KafkaErrorsTopic())

	return &netpoll.Config{
		OnWaitError: func(err error) {
			go func(err error) {
				errEvent := server.NewErrorEvent("", "PollerConfig.OnWaitError()", err.Error())
				errorProducer.Publish(p.ctx, []byte(errEvent.HostID), errEvent.ToJSON())
			}(err)
		},
	}
}

func (p *Providers) HandleMessageProvider() server.HandleMessage {
	return server.NewHandleMessage(
		p.MessageEncoderProvider(),
		p.KafkaProvider().NewProducer(config.KafkaNewMessagesTopic()),
	)
}

func (p *Providers) HandleGroupMessageProvider() server.HandleGroupMessage {
	return server.NewHandleGroupMessage(
		p.RepositoryProvider(),
	)
}

func (p *Providers) HandleOffMessageProvider() server.HandleMessage {
	return server.NewHandleMessage(
		p.MessageEncoderProvider(),
		p.KafkaProvider().NewProducer(config.KafkaOffMessagesTopic()),
	)
}

func (p *Providers) HandleUserStatusProvider() server.HandleUserStatus {
	return server.NewHandleUserStatus(
		p.UserEncoderProvider(),
		p.KafkaProvider().NewProducer(config.KafkaUsersTopic()),
		p.RepositoryProvider(),
	)
}

func (p *Providers) HandleErrorProvider() server.HandleError {
	return server.NewHandleError(
		p.EventErrorEncoderProvider(),
		p.KafkaProvider().NewProducer(config.KafkaErrorsTopic()),
	)
}

func (p *Providers) JwtProvider() auth.JWT {
	if p.jwt == nil {
		p.jwt = auth.NewJWT(config.PathPrivateKey(), config.PathPublicKey())
	}
	return p.jwt
}

func (p *Providers) RouterProvider() (*api.Router, error) {
	server, err := p.ServerProvider()
	if err != nil {
		return nil, err
	}

	jwt := p.JwtProvider()
	handler := api.HandleWS(jwt, server)
	router := api.NewRouter(middleware.NewAuth(jwt), handler)

	return router, nil
}