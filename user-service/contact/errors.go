package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrUserNotFound = errors.New("user not found")
var ErrContactNotFound = errors.New("contact not found")
var ErrContactAlreadyBlocked = errors.New("contact already blocked")
var ErrContactAlreadyExists = errors.New("contact already exists")
var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required ID"}
var ErrUserIDValidateModel = &cerror.ErrValidateModel{Msg: "required UserID"}
