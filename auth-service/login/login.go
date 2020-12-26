package login

import "github.com/tsmweb/go-helper-api/cerror"

// Login data model.
type Login struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// Validate model login.
func (l Login) Validate() error {
	if l.ID == "" {
		return &cerror.ErrValidateModel{Msg: "Required ID"}
	}
	if l.Password == "" {
		return &cerror.ErrValidateModel{Msg: "Required Password"}
	}

	return nil
}
