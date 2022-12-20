package main

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/tsmweb/chat-service/adapter"
	"github.com/tsmweb/chat-service/common/service"
	"github.com/tsmweb/chat-service/config"
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

type Provider struct {
	ctx    context.Context
	server *server.Server
	jwt    auth.JWT
	mAuth  middleware.Auth
	kafka  kafka.Kafka
}

func CreateProvider(ctx context.Context) *Provider {
	return &Provider{
		ctx: ctx,
	}
}

func (p *Provider) ServerProvider() (*server.Server, error) {
	if p.server == nil {
		poll, err := p.EpollProvider()
		if err != nil {
			return nil, err
		}

		connReader := server.ConnReaderFunc(adapter.ReaderWS)
		connWriter := server.ConnWriterFunc(adapter.WriterWS)

		messageDecoder := message.DecoderFunc(adapter.MessageUnmarshal)
		messageEncoder := message.EncoderFunc(adapter.MessageMarshal)
		userEncoder := user.EncoderFunc(adapter.UserMarshal)

		messageConsumer := p.KafkaProvider().NewConsumer(config.KafkaClientID(),
			config.KafkaHostTopic())
		messageProducer := p.KafkaProvider().NewProducer(config.KafkaNewMessagesTopic())
		offMessageProducer := p.KafkaProvider().NewProducer(config.KafkaOffMessagesTopic())
		userProducer := p.KafkaProvider().NewProducer(config.KafkaUsersTopic())
		userPresenceProducer := p.KafkaProvider().NewProducer(config.KafkaUsersPresenceTopic())

		handleMessage := server.NewHandleMessage(messageEncoder, messageProducer)
		handleOffMessage := server.NewHandleMessage(messageEncoder, offMessageProducer)
		handleUserStatus := server.NewHandleUserStatus(userEncoder, userProducer,
			userPresenceProducer)

		p.server = server.NewServer(
			p.ctx,
			poll,
			connReader,
			connWriter,
			messageDecoder,
			messageConsumer,
			handleMessage,
			handleOffMessage,
			handleUserStatus,
		)
	}
	return p.server, nil
}

func (p *Provider) EpollProvider() (epoll.EPoll, error) {
	poller, err := netpoll.New(p.PollerConfigProvider())
	if err != nil {
		return nil, err
	}
	return epoll.NewEPoll(poller), nil
}

// PollerConfigProvider OnWaitError will be called from goroutine, waiting for events.
func (p *Provider) PollerConfigProvider() *netpoll.Config {
	return &netpoll.Config{
		OnWaitError: func(err error) {
			go func(err error) {
				service.Error("", "PollerConfig.OnWaitError()", err)
			}(err)
		},
	}
}

func (p *Provider) NewKafkaProducer(topic string) kafka.Producer {
	return p.KafkaProvider().NewProducer(topic)
}

func (p *Provider) KafkaProvider() kafka.Kafka {
	if p.kafka == nil {
		p.kafka = kafka.New([]string{config.KafkaBootstrapServers()}, config.KafkaClientID())
	}
	return p.kafka
}

func (p *Provider) JwtProvider() auth.JWT {
	if p.jwt == nil {
		p.jwt = auth.NewJWT(config.KeySecureFile(), config.PubSecureFile())
	}
	return p.jwt
}

func (p *Provider) AuthProvider() middleware.Auth {
	if p.mAuth == nil {
		p.mAuth = middleware.NewAuth(p.JwtProvider())
	}
	return p.mAuth
}

func (p *Provider) ChatRouter(mr *mux.Router) error {
	serv, err := p.ServerProvider()
	if err != nil {
		return err
	}

	api.MakeChatRouter(
		mr,
		p.JwtProvider(),
		p.AuthProvider(),
		serv,
	)

	return nil
}
