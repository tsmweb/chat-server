package user

// Presenter data
type Presenter struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Password string `json:"password,omitempty"`
}

// ToEntity mapper Presenter to Entity
func (v *Presenter) ToEntity() *User {
	return &User{
		ID:       v.ID,
		Name:     v.Name,
		LastName: v.LastName,
		Password: v.Password,
	}
}

// FromEntity mapper Entity to Presenter
func (v *Presenter) FromEntity(entity *User) {
	v.ID = entity.ID
	v.Name = entity.Name
	v.LastName = entity.LastName
}
