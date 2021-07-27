package main

import (
	"context"
	"github.com/tsmweb/chat-service/adapter"
	"github.com/tsmweb/chat-service/chat"
	"github.com/tsmweb/chat-service/chat/message"
	"github.com/tsmweb/chat-service/chat/user"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/infra/db"
	"github.com/tsmweb/chat-service/pkg/epoll"
	"github.com/tsmweb/chat-service/web/api"
	"github.com/tsmweb/easygo/netpoll"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/kafka"
	"github.com/tsmweb/go-helper-api/middleware"
)

type Providers struct {
	ctx context.Context
	server       *chat.Server
	jwt        auth.JWT
	dataBase   db.Database
	repository chat.Repository
	kafka      kafka.Kafka
}

func CreateProvider(ctx context.Context) *Providers {
	return &Providers{
		ctx: ctx,
	}
}

func (p *Providers) ServerProvider() (*chat.Server, error) {
	if p.server == nil {
		epoll, err := p.EpollProvider()
		if err != nil {
			return nil, err
		}

		p.server = chat.NewServer(
			p.ctx,
			epoll,
			p.ConnReaderProvider(),
			p.ConnWriterProvider(),
			p.MessageEncoderProvider(),
			p.MessageDecoderProvider(),
			p.UserEncoderProvider(),
			p.RepositoryProvider(),
			p.KafkaProvider(),
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

func (p *Providers) ConnReaderProvider() chat.ConnReader {
	return chat.ConnReaderFunc(adapter.ReaderWS)
}

func (p *Providers) ConnWriterProvider() chat.ConnWriter {
	return chat.ConnWriterFunc(adapter.WriterWS)
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

func (p *Providers) RepositoryProvider() chat.Repository {
	if p.repository == nil {
		p.repository = chat.NewMemoryRepository()
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
				errEvent := chat.NewErrorEvent("", "PollerConfig.OnWaitError()", err.Error())
				errorProducer.Publish(p.ctx, []byte(errEvent.HostID), errEvent.ToJSON())
			}(err)
		},
	}
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