package contact

// Presenter data
type Presenter struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	UserID   string `json:"user_id,omitempty"`
}

// ToEntity mapper Presenter to Entity
func (v *Presenter) ToEntity() *Contact {
	return &Contact{
		ID:       v.ID,
		Name:     v.Name,
		LastName: v.LastName,
		UserID:   v.UserID,
	}
}

// FromEntity mapper Entity to Presenter
func (v *Presenter) FromEntity(entity *Contact) {
	v.ID = entity.ID
	v.Name = entity.Name
	v.LastName = entity.LastName
	v.UserID = entity.UserID
}

// Presence data
type Presence struct {
	ID       string `json:"id"`
	Presence string `json:"presence"`
}
