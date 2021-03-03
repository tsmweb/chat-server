package group

// Presenter data
type Presenter struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Owner string `json:"owner"`
}

// ToEntity mapper Presenter to Entity
func (p *Presenter) ToEntity() *Group {
	return &Group{
		ID: p.ID,
		Name: p.Name,
		Description: p.Description,
		Owner: p.Owner,
	}
}

// FromEntity mapper Entity to Presenter
func (p *Presenter) FromEntity(entity *Group) {
	p.ID = entity.ID
	p.Name = entity.Name
	p.Description = entity.Description
	p.Owner = entity.Owner
}