package main

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/tsmweb/auth-service/app/login"
	"github.com/tsmweb/auth-service/app/user"
	"github.com/tsmweb/auth-service/config"
	"github.com/tsmweb/auth-service/infra/db"
	"github.com/tsmweb/auth-service/infra/repository"
	"github.com/tsmweb/auth-service/web/api/handler"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/kafka"
	"github.com/tsmweb/go-helper-api/middleware"
)

type Provider struct {
	ctx      context.Context
	jwt      auth.JWT
	mAuth    middleware.Auth
	dataBase db.Database
	kafka    kafka.Kafka
}

func CreateProvider(ctx context.Context) *Provider {
	return &Provider{
		ctx: ctx,
	}
}

func (p *Provider) UserRouter(mr *mux.Router) {
	repository := repository.NewUserRepositoryPostgres(p.DatabaseProvider())
	getUseCase := user.NewGetUseCase(repository)
	createUseCase := user.NewCreateUseCase(repository)
	updateUseCase := user.NewUpdateUseCase(repository)

	handler.MakeUserHandlers(
		mr,
		p.JwtProvider(),
		p.AuthProvider(),
		getUseCase,
		createUseCase,
		updateUseCase)
}

func (p *Provider) LoginRouter(mr *mux.Router) {
	repository := repository.NewLoginRepositoryPostgres(p.DatabaseProvider())
	loginUseCase := login.NewLoginUseCase(repository, p.JwtProvider())
	updateUseCase := login.NewUpdateUseCase(repository)

	handler.MakeLoginHandlers(
		mr,
		p.JwtProvider(),
		p.AuthProvider(),
		loginUseCase,
		updateUseCase)
}

func (p *Provider) NewKafkaProducer(topic string) kafka.Producer {
	return p.KafkaProvider().NewProducer(topic)
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

func (p *Provider) DatabaseProvider() db.Database {
	if p.dataBase == nil {
		p.dataBase = db.NewPostgresDatabase()
	}
	return p.dataBase
}

func (p *Provider) KafkaProvider() kafka.Kafka {
	if p.kafka == nil {
		p.kafka = kafka.New([]string{config.KafkaBootstrapServers()}, config.KafkaClientID())
	}
	return p.kafka
}
