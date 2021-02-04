package user

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required id"}
var ErrNameValidateModel = &cerror.ErrValidateModel{Msg: "required name"}
var ErrPasswordValidateModel = &cerror.ErrValidateModel{Msg: "required password"}
