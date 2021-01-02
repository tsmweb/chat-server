package profile

// ViewModel data
type ViewModel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Password string `json:"password,omitempty"`
}

// ToEntity mapper ViewModel to Entity
func (v *ViewModel) ToEntity() Profile {
	return Profile{
		ID:       v.ID,
		Name:     v.Name,
		LastName: v.LastName,
		Password: v.Password,
	}
}

// FromEntity mapper Entity to ViewModel
func (v *ViewModel) FromEntity(entity Profile) {
	v.ID = entity.ID
	v.Name = entity.Name
	v.LastName = entity.LastName
}
