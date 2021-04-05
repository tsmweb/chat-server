package group

import (
	"time"
)

// Presenter data
type Presenter struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Members     []*MemberPresenter `json:"members"`
	UpdatedBy   string    `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToEntity mapper Presenter to Entity
func (p *Presenter) ToEntity() *Group {
	return &Group{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Owner:       p.Owner,
	}
}

// FromEntity mapper Entity to Presenter
func (p *Presenter) FromEntity(entity *Group) {
	p.ID = entity.ID
	p.Name = entity.Name
	p.Description = entity.Description
	p.Owner = entity.Owner
	p.UpdatedBy = entity.UpdatedBy
	p.CreatedAt = entity.CreatedAt
	p.UpdatedAt = entity.UpdatedAt

	var members []*MemberPresenter

	for _, member := range entity.Members {
		vm := &MemberPresenter{}
		vm.FromEntity(member)
		members = append(members, vm)
	}

	p.Members = members
}

// MemberPresenter data
type MemberPresenter struct {
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Admin     bool      `json:"admin"`
	UpdatedBy string    `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToEntity mapper MemberPresenter to Entity
func (p *MemberPresenter) ToEntity() *Member {
	return &Member{
		GroupID: p.GroupID,
		UserID:  p.UserID,
		Admin:   p.Admin,
	}
}

// FromEntity mapper Entity to MemberPresenter
func (p *MemberPresenter) FromEntity(entity *Member) {
	p.GroupID = entity.GroupID
	p.UserID = entity.UserID
	p.Admin = entity.Admin
	p.UpdatedBy = entity.UpdatedBy
	p.CreatedAt = entity.CreatedAt
	p.UpdatedAt = entity.UpdatedAt
}

// EntityToPresenters mapper Entities to Presenters
func EntityToPresenters(entities ...*Group) []*Presenter {
	var vms []*Presenter

	for _, group := range entities {
		vm := &Presenter{}
		vm.FromEntity(group)
		vms = append(vms, vm)
	}

	return vms
}
