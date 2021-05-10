//+build wireinject

package api

import (
	"github.com/google/wire"
)

var Inject = wire.NewSet(
	NewController,
	NewRouter,
)
