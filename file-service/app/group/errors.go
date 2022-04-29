package group

import "errors"

var (
	ErrGroupNotFound        = errors.New("group not found")
	ErrOperationNotAllowed  = errors.New("operation not allowed")
	ErrUnsupportedMediaType = errors.New("supported media type is image/jpeg")
	ErrFileTooBig           = errors.New("the uploaded file is too big")
)
