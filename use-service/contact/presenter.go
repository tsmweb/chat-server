package contact

// Presenter data
type Presenter struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
}

// ToEntity mapper Presenter to Entity
func (v *Presenter) ToEntity() *Contact {
	return &Contact{
		ID:       v.ID,
		Name:     v.Name,
		LastName: v.LastName,
	}
}

// FromEntity mapper Entity to Presenter
func (v *Presenter) FromEntity(entity *Contact) {
	v.ID = entity.ID
	v.Name = entity.Name
	v.LastName = entity.LastName
}
