package main

import (
	"context"
	"github.com/tsmweb/broker-service/adapter"
	"github.com/tsmweb/broker-service/broker"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/config"
	"github.com/tsmweb/broker-service/infra/db"
	"github.com/tsmweb/broker-service/infra/repository"
	"github.com/tsmweb/go-helper-api/kafka"
)

type Provider struct {
	ctx      context.Context
	database db.Database
	cache db.CacheDB
	kafka    kafka.Kafka
	broker   *broker.Broker
}

func CreateProvider(ctx context.Context) *Provider {
	return &Provider{
		ctx: ctx,
	}
}

func (p *Provider) BrokerProvider() *broker.Broker {
	if p.broker == nil {
		userDecoder := user.DecoderFunc(adapter.UserUnmarshal)
		messageEncoder := message.EncoderFunc(adapter.MessageMarshal)
		messageDecoder := message.DecoderFunc(adapter.MessageUnmarshal)
		errorEncoder := broker.ErrorEventEncoderFunc(adapter.ErrorEventMarshal)

		userConsumer := p.KafkaProvider().NewConsumer(config.KafkaGroupID(), config.KafkaUsersTopic())
		userPresenceConsumer := p.KafkaProvider().NewConsumer("", config.KafkaUsersPresenceTopic())
		messageConsumer := p.KafkaProvider().NewConsumer(config.KafkaGroupID(), config.KafkaNewMessagesTopic())
		offMessageConsumer := p.KafkaProvider().NewConsumer(config.KafkaGroupID(), config.KafkaOffMessagesTopic())
		errorProducer := p.KafkaProvider().NewProducer(config.KafkaErrorsTopic())

		userRepository := repository.NewUserRepository(p.DatabaseProvider(), p.CacheDBProvider())
		messageRepository := repository.NewMessageRepository(p.DatabaseProvider(), p.CacheDBProvider())

		userHandler := broker.NewUserHandler(userRepository, messageRepository)
		userPresenceHandler := broker.NewUserPresenceHandler(userRepository)
		messageHandler := broker.NewMessageHandler(userRepository, messageRepository, p.KafkaProvider(), messageEncoder)
		offMessageHandler := broker.NewOfflineMessageHandler(messageRepository)
		errorHandler := broker.NewErrorHandler(errorEncoder, errorProducer)

		p.broker = broker.NewBroker(
			p.ctx,
			userDecoder,
			messageDecoder,
			userConsumer,
			userPresenceConsumer,
			messageConsumer,
			offMessageConsumer,
			userHandler,
			userPresenceHandler,
			messageHandler,
			offMessageHandler,
			errorHandler,
		)
	}
	return p.broker
}

func (p *Provider) DatabaseProvider() db.Database {
	if p.database == nil {
		p.database = db.NewPostgresDatabase()
	}
	return p.database
}

func (p *Provider) CacheDBProvider() db.CacheDB {
	if p.cache == nil {
		p.cache = db.NewRedisCacheDB(config.RedisHost(), config.RedisPassword())
	}
	return p.cache
}

func (p *Provider) KafkaProvider() kafka.Kafka {
	if p.kafka == nil {
		p.kafka = kafka.New([]string{config.KafkaBootstrapServers()}, config.KafkaClientID())
	}
	return p.kafka
}
