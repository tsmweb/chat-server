package login

import (
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required id"}
var ErrPasswordValidateModel = &cerror.ErrValidateModel{Msg: "required password"}
