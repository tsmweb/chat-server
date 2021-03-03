package group

import "github.com/tsmweb/go-helper-api/cerror"

var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required id"}
var ErrNameValidateModel = &cerror.ErrValidateModel{Msg: "required name"}
var ErrOwnerValidateModel = &cerror.ErrValidateModel{Msg: "required owner"}
var ErrGroupIDValidateModel = &cerror.ErrValidateModel{Msg: "required group_id"}
var ErrUserIDValidateModel = &cerror.ErrValidateModel{Msg: "required user_id"}
