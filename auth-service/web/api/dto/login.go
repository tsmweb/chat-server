package dto

import "github.com/tsmweb/auth-service/login"

// Login data.
type Login struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// ToEntity mapper Login to login.Login
func (p *Login) ToEntity() *login.Login {
	return &login.Login{
		ID: p.ID,
		Password: p.Password,
	}
}
