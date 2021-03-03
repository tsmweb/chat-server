package user

import "time"

// Presenter data
type Presenter struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Password string `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToEntity mapper Presenter to Entity
func (p *Presenter) ToEntity() *User {
	return &User{
		ID:        p.ID,
		Name:      p.Name,
		LastName:  p.LastName,
		Password:  p.Password,
	}
}

// FromEntity mapper Entity to Presenter
func (p *Presenter) FromEntity(entity *User) {
	p.ID = entity.ID
	p.Name = entity.Name
	p.LastName = entity.LastName
	p.CreatedAt = entity.CreatedAt
	p.UpdatedAt = entity.UpdatedAt
}
