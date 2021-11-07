package main

import (
	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/config"
	"github.com/tsmweb/file-service/group"
	"github.com/tsmweb/file-service/infra/db"
	"github.com/tsmweb/file-service/infra/repository"
	"github.com/tsmweb/file-service/web/handler"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/middleware"
)

type Provider struct {
	jwt      auth.JWT
	mAuth    middleware.Auth
	dataBase db.Database
}

func CreateProvider() *Provider {
	return &Provider{}
}

func (p *Provider) UserRouter(mr *mux.Router) {
	handler.MakeUserHandlers(
		mr,
		p.JwtProvider(),
		p.AuthProvider())
}

func (p *Provider) GroupRouter(mr *mux.Router) {
	database := p.DatabaseProvider()
	repo := repository.NewGroupRepositoryPostgres(database)
	validateUseCase := group.NewValidateUseCase(repo)

	handler.MakeGroupHandlers(
		mr,
		p.JwtProvider(),
		p.AuthProvider(),
		validateUseCase)
}

func (p *Provider) JwtProvider() auth.JWT {
	if p.jwt == nil {
		p.jwt = auth.NewJWT(config.PathPrivateKey(), config.PathPublicKey())
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
