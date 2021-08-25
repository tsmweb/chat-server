package main

import (
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/user-service/config"
	"github.com/tsmweb/user-service/contact"
	"github.com/tsmweb/user-service/group"
	"github.com/tsmweb/user-service/infra/db"
	"github.com/tsmweb/user-service/infra/repository"
	"github.com/tsmweb/user-service/web/api/handler"
)

type Provider struct {
	jwt      auth.JWT
	mAuth    middleware.Auth
	dataBase db.Database
}

func CreateProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ContactRouter(mr *mux.Router) {
	database := p.DatabaseProvider()
	repo := repository.NewContactRepositoryPostgres(database)
	getUseCase := contact.NewGetUseCase(repo)
	getAllUseCase := contact.NewGetAllUseCase(repo)
	getPresenceUseCase := contact.NewGetPresenceUseCase(repo)
	createUseCase := contact.NewCreateUseCase(repo)
	updateUseCase := contact.NewUpdateUseCase(repo)
	deleteUseCase := contact.NewDeleteUseCase(repo)
	blockUseCase := contact.NewBlockUseCase(repo)
	unblockUseCase := contact.NewUnblockUseCase(repo)

	handler.MakeContactRouters(
		mr,
		p.JwtProvider(),
		p.AuthProvider(),
		getUseCase,
		getAllUseCase,
		getPresenceUseCase,
		createUseCase,
		updateUseCase,
		deleteUseCase,
		blockUseCase,
		unblockUseCase)
}

func (p *Provider) GroupRouter(mr *mux.Router) {
	database := p.DatabaseProvider()
	repo := repository.NewGroupRepositoryPostgres(database)
	getUseCase := group.NewGetUseCase(repo)
	getAllUseCase := group.NewGetAllUseCase(repo)
	createUseCase := group.NewCreateUseCase(repo)
	updateUseCase := group.NewUpdateUseCase(repo)
	deleteUseCase := group.NewDeleteUseCase(repo)
	addMemberUseCase := group.NewAddMemberUseCase(repo)
	removeMemberUseCase := group.NewRemoveMemberUseCase(repo)
	setAdminUseCase := group.NewSetAdminUseCase(repo)

	handler.MakeGroupRouters(
		mr,
		p.JwtProvider(),
		p.AuthProvider(),
		getUseCase,
		getAllUseCase,
		createUseCase,
		updateUseCase,
		deleteUseCase,
		addMemberUseCase,
		removeMemberUseCase,
		setAdminUseCase)
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