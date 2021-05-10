//+build wireinject

package core

import (
	"github.com/google/wire"
)

var Inject = wire.NewSet(
	NewErrorDispatcher,
	NewPresenceDispatcher,
	NewGroupMessageDispatcher,
	NewOfflineMessageDispatcher,
	NewMemoryRepository,
	NewUserStatusHandler,
	NewMessageHandler,
	NewChat,
)
