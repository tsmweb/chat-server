package contact

import "time"

// Presenter data
type Presenter struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastname"`
	UserID    string    `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToEntity mapper Presenter to Entity
func (p *Presenter) ToEntity() *Contact {
	return &Contact{
		ID:       p.ID,
		Name:     p.Name,
		LastName: p.LastName,
		UserID:   p.UserID,
	}
}

// FromEntity mapper Entity to Presenter
func (p *Presenter) FromEntity(entity *Contact) {
	p.ID = entity.ID
	p.Name = entity.Name
	p.LastName = entity.LastName
	p.UserID = entity.UserID
	p.CreatedAt = entity.CreatedAt
	p.UpdatedAt = entity.UpdatedAt
}

// Presence data
type Presence struct {
	ID       string `json:"id"`
	Presence string `json:"presence"`
}

// EntityToPresenters mapper Entities to Presenters
func EntityToPresenters(entities ...*Contact) []*Presenter {
	var vms []*Presenter

	for _, contact := range entities {
		vm := &Presenter{}
		vm.FromEntity(contact)
		vms = append(vms, vm)
	}

	return vms
}
