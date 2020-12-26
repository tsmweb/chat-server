package profile

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrProfileNotFound = errors.New("Presenter Not Found")
var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "Required ID"}
var ErrNameValidateModel = &cerror.ErrValidateModel{Msg: "Required Name"}
var ErrPasswordValidateModel = &cerror.ErrValidateModel{Msg: "Required Password"}
