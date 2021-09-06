package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/tsmweb/chat-service/adapter"
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

type Providers struct {
	ctx    context.Context
	server *server.Server
	jwt    auth.JWT
	mAuth  middleware.Auth
	kafka  kafka.Kafka
}

func CreateProvider(ctx context.Context) *Providers {
	return &Providers{
		ctx: ctx,
	}
}

func (p *Providers) ServerProvider() (*server.Server, error) {
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
		errorEncoder := server.ErrorEventEncoderFunc(adapter.ErrorEventMarshal)

		messageConsumer := p.KafkaProvider().NewConsumer(config.KafkaGroupID(), config.KafkaHostTopic())
		messageProducer := p.KafkaProvider().NewProducer(config.KafkaNewMessagesTopic())
		offMessageProducer := p.KafkaProvider().NewProducer(config.KafkaOffMessagesTopic())
		userProducer := p.KafkaProvider().NewProducer(config.KafkaUsersTopic())
		errorProducer := p.KafkaProvider().NewProducer(config.KafkaErrorsTopic())

		handleMessage := server.NewHandleMessage(messageEncoder, messageProducer)
		handleOffMessage := server.NewHandleMessage(messageEncoder, offMessageProducer)
		handleUserStatus := server.NewHandleUserStatus(userEncoder, userProducer)
		handleError := server.NewHandleError(errorEncoder, errorProducer)

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
			handleError,
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

func (p *Providers) KafkaProvider() kafka.Kafka {
	if p.kafka == nil {
		p.kafka = kafka.New([]string{config.KafkaBootstrapServers()}, config.KafkaClientID())
	}
	return p.kafka
}

func (p *Providers) JwtProvider() auth.JWT {
	if p.jwt == nil {
		p.jwt = auth.NewJWT(config.PathPrivateKey(), config.PathPublicKey())
	}
	return p.jwt
}

func (p *Providers) AuthProvider() middleware.Auth {
	if p.mAuth == nil {
		p.mAuth = middleware.NewAuth(p.JwtProvider())
	}
	return p.mAuth
}

func (p *Providers) ChatRouter(mr *mux.Router) error {
	serv, err := p.ServerProvider()
	if err != nil {
		return err
	}

	api.MakeChatRouter(
		mr,
		p.JwtProvider(),
		p.AuthProvider(),
		serv)

	return nil
}
