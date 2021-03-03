package group

import "time"

// Group data model
type Group struct {
	ID          string
	Name        string
	Description string
	Owner       string
	Members     []*Member
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewGroup create a new Group
func NewGroup(ID, name, description, owner string) (*Group, error) {
	g := &Group{
		ID:          ID,
		Name:        name,
		Description: description,
		Owner:       owner,
		CreatedAt:   time.Now(),
	}

	err := g.Validate()
	if err != nil {
		return nil, err
	}

	return g, nil
}

// Validate model Group
func (g *Group) Validate() error {
	if g.ID == "" {
		return ErrIDValidateModel
	}
	if g.Name == "" {
		return ErrNameValidateModel
	}
	if g.Owner == "" {
		return ErrOwnerValidateModel
	}
	return nil
}
