package user

import "errors"

var (
	ErrUnsupportedMediaType = errors.New("supported media type is image/jpeg")
	ErrFileTooBig           = errors.New("the uploaded file is too big")
)
