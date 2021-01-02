package login

import (
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required ID"}
var ErrPasswordValidateModel = &cerror.ErrValidateModel{Msg: "required Password"}
