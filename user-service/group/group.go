package group

import (
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"strconv"
	"time"
)

// Group data model
type Group struct {
	ID          string
	Name        string
	Description string
	Owner       string
	Members     []*Member
	UpdatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewGroup create a new Group
func NewGroup(name, description, owner string) (*Group, error) {
	ID, err := hashutil.HashSHA1(owner + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return nil, err
	}

	g := &Group{
		ID:          ID,
		Name:        name,
		Description: description,
		Owner:       owner,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	err = g.Validate()
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
