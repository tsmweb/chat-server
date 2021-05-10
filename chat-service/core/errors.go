package core

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required id"}
var ErrFromValidateModel = &cerror.ErrValidateModel{Msg: "required from"}
var ErrReceiverValidateModel = &cerror.ErrValidateModel{Msg: "required to or group"}
var ErrDateValidateModel = &cerror.ErrValidateModel{Msg: "required date"}
var ErrContentTypeValidateModel = &cerror.ErrValidateModel{Msg: "required content_type"}
var ErrContentValidateModel = &cerror.ErrValidateModel{Msg: "required content"}
var ErrMessageSpoof = errors.New("internal not allowed")
var ErrClosedChannel = errors.New("closed channel")
