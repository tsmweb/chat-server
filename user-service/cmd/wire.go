//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/user-service/contact"
	"github.com/tsmweb/user-service/group"
	"github.com/tsmweb/user-service/helper/database"
	"github.com/tsmweb/user-service/helper/setting"
)

func InitContactRouter() *contact.Router {
	wire.Build(
		contact.NewRouter,
		middleware.NewAuth,
		contact.NewController,
		contact.NewService,
		contact.NewGetUseCase,
		contact.NewGetAllUseCase,
		contact.NewGetPresenceUseCase,
		contact.NewCreateUseCase,
		contact.NewUpdateUseCase,
		contact.NewDeleteUseCase,
		contact.NewBlockUseCase,
		contact.NewUnblockUseCase,
		contact.NewRepositoryPostgres,
		jwtProvider,
		dataBaseProvider)

	return &contact.Router{}
}

func InitGroupRouter() *group.Router {
	wire.Build(
		group.NewRouter,
		middleware.NewAuth,
		group.NewController,
		group.NewService,
		group.NewGetUseCase,
		group.NewGetAllUseCase,
		group.NewCreateUseCase,
		group.NewUpdateUseCase,
		group.NewDeleteUseCase,
		group.NewAddMemberUseCase,
		group.NewRemoveMemberUseCase,
		group.NewSetAdminUseCase,
		group.NewRepositoryPostgres,
		jwtProvider,
		dataBaseProvider)

	return &group.Router{}
}

/*
 * PROVIDERS
 */

// Data Base
var databaseInstance database.Database

func dataBaseProvider() database.Database {
	if databaseInstance == nil {
		databaseInstance = database.NewPostgresDatabase()
	}

	return databaseInstance
}

// Authentication JWT
var jwtInstance auth.JWT

func jwtProvider() auth.JWT {
	if jwtInstance == nil {
		jwtInstance = auth.NewJWT(setting.PathPrivateKey(), setting.PathPublicKey())
	}

	return jwtInstance
}
