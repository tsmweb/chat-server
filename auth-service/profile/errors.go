package profile

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrProfileNotFound = errors.New("profile not found")
var ErrProfileAlreadyExists = errors.New("profile already exists")
var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required ID"}
var ErrNameValidateModel = &cerror.ErrValidateModel{Msg: "required Name"}
var ErrPasswordValidateModel = &cerror.ErrValidateModel{Msg: "required Password"}
