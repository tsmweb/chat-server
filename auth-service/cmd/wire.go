//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/tsmweb/auth-service/helper/database"
	"github.com/tsmweb/auth-service/helper/setting"
	"github.com/tsmweb/auth-service/profile"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/middleware"
)

func InitProfileRouter() *profile.Router {
	wire.Build(
		profile.NewRouter,
		middleware.NewAuth,
		profile.NewController,
		profile.NewGetUseCase,
		profile.NewCreateUseCase,
		profile.NewUpdateUseCase,
		profile.NewRepositoryPostgres,
		jwtProvider,
		dataBaseProvider)

	return &profile.Router{}
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

