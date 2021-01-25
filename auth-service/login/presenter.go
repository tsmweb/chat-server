package login

// Presenter data.
type Presenter struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// ToEntity mapper Presenter to Entity
func (p *Presenter) ToEntity() *Login {
	return &Login{
		ID: p.ID,
		Password: p.Password,
	}
}
