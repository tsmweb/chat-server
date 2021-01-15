package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrProfileNotFound = errors.New("profile not found")
var ErrContactNotFound = errors.New("contact not found")
var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required ID"}
var ErrProfileIDValidateModel = &cerror.ErrValidateModel{Msg: "required ProfileID"}