package profile

// Presenter data
type Presenter struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Password string `json:"password, omitempty"`
}

// ToEntity mapper Presenter to Entity
func (p *Presenter) ToEntity() Profile {
	return Profile{
		ID: p.ID,
		Name: p.Name,
		LastName: p.LastName,
		Password: p.Password,
	}
}

// FromEntity mapper Entity to Presenter
func (p *Presenter) FromEntity(entity Profile) {
	p.ID = entity.ID
	p.Name = entity.Name
	p.LastName = entity.LastName
}
