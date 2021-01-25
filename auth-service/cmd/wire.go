//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/tsmweb/auth-service/helper/database"
	"github.com/tsmweb/auth-service/helper/setting"
	"github.com/tsmweb/auth-service/login"
	"github.com/tsmweb/auth-service/user"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/middleware"
)

func InitUserRouter() *user.Router {
	wire.Build(
		user.NewRouter,
		middleware.NewAuth,
		user.NewController,
		user.NewService,
		user.NewGetUseCase,
		user.NewCreateUseCase,
		user.NewUpdateUseCase,
		user.NewRepositoryPostgres,
		jwtProvider,
		dataBaseProvider)

	return &user.Router{}
}

func InitLoginRouter() *login.Router {
	wire.Build(
		login.NewRoutes,
		middleware.NewAuth,
		login.NewController,
		login.NewService,
		login.NewLoginUseCase,
		login.NewUpdateUseCase,
		login.NewRepositoryPostgres,
		jwtProvider,
		dataBaseProvider)

	return &login.Router{}
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

