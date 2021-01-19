package login

// ViewModel data.
type ViewModel struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// ToEntity mapper ViewModel to Entity
func (p *ViewModel) ToEntity() *Login {
	return &Login{
		ID: p.ID,
		Password: p.Password,
	}
}
