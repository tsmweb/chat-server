package main

import (
	"context"

	"github.com/tsmweb/broker-service/adapter"
	"github.com/tsmweb/broker-service/broker"
	"github.com/tsmweb/broker-service/broker/group"
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
	cache    db.CacheDB
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
		groupEventDecoder := group.EventDecoderFunc(adapter.GroupEventUnmarshal)
		userEventDecoder := user.EventDecoderFunc(adapter.UserEventUnmarshal)

		userConsumer := p.KafkaProvider().NewConsumer(config.KafkaGroupID(),
			config.KafkaUsersTopic())
		userPresenceConsumer := p.KafkaProvider().NewConsumer(config.KafkaClientID(),
			config.KafkaUsersPresenceTopic())
		messageConsumer := p.KafkaProvider().NewConsumer(config.KafkaGroupID(),
			config.KafkaNewMessagesTopic())
		offMessageConsumer := p.KafkaProvider().NewConsumer(config.KafkaGroupID(),
			config.KafkaOffMessagesTopic())
		groupEventConsumer := p.KafkaProvider().NewConsumer(config.KafkaClientID(),
			config.KafkaGroupEventTopic())
		userEventConsumer := p.KafkaProvider().NewConsumer(config.KafkaClientID(),
			config.KafkaContactEventTopic())

		userRepository := repository.NewUserRepository(p.DatabaseProvider(), p.CacheDBProvider())
		messageRepository := repository.NewMessageRepository(p.DatabaseProvider(),
			p.CacheDBProvider())

		userHandler := broker.NewUserHandler(userRepository, messageRepository)
		userPresenceHandler := broker.NewUserPresenceHandler(userRepository)
		messageHandler := broker.NewMessageHandler(userRepository, messageRepository,
			p.KafkaProvider(), messageEncoder)
		offMessageHandler := broker.NewOfflineMessageHandler(messageRepository)
		groupEventHandler := broker.NewGroupEventHandler(messageRepository)
		userEventHandler := broker.NewUserEventHandler(userRepository)

		p.broker = broker.NewBroker(
			p.ctx,
			userDecoder,
			messageDecoder,
			groupEventDecoder,
			userEventDecoder,
			userConsumer,
			userPresenceConsumer,
			messageConsumer,
			offMessageConsumer,
			groupEventConsumer,
			userEventConsumer,
			userHandler,
			userPresenceHandler,
			messageHandler,
			offMessageHandler,
			groupEventHandler,
			userEventHandler,
		)
	}
	return p.broker
}

func (p *Provider) NewKafkaProducer(topic string) kafka.Producer {
	return p.KafkaProvider().NewProducer(topic)
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
