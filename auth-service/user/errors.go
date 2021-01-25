package user

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required ID"}
var ErrNameValidateModel = &cerror.ErrValidateModel{Msg: "required Name"}
var ErrPasswordValidateModel = &cerror.ErrValidateModel{Msg: "required Password"}
