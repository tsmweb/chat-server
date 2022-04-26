package dto

import (
	"github.com/tsmweb/user-service/app/group"
	"time"
)

// Group data
type Group struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Owner       string    `json:"owner"`
	Members     []*Member `json:"members,omitempty"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// ToEntity mapper dto.Group to group.Group
func (g *Group) ToEntity() *group.Group {
	return &group.Group{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Owner:       g.Owner,
	}
}

// FromEntity mapper group.Group to dto.Group
func (g *Group) FromEntity(entity *group.Group) {
	g.ID = entity.ID
	g.Name = entity.Name
	g.Description = entity.Description
	g.Owner = entity.Owner
	g.UpdatedBy = entity.UpdatedBy
	g.CreatedAt = entity.CreatedAt
	g.UpdatedAt = entity.UpdatedAt

	var members []*Member

	for _, member := range entity.Members {
		m := &Member{}
		m.FromEntity(member)
		members = append(members, m)
	}

	g.Members = members
}

// Member data
type Member struct {
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Admin     bool      `json:"admin"`
	UpdatedBy string    `json:"updated_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// ToEntity mapper dto.Member to group.Member
func (m *Member) ToEntity() *group.Member {
	return &group.Member{
		GroupID: m.GroupID,
		UserID:  m.UserID,
		Admin:   m.Admin,
	}
}

// FromEntity mapper group.Member to dto.Member
func (m *Member) FromEntity(entity *group.Member) {
	m.GroupID = entity.GroupID
	m.UserID = entity.UserID
	m.Admin = entity.Admin
	m.UpdatedBy = entity.UpdatedBy
	m.CreatedAt = entity.CreatedAt
	m.UpdatedAt = entity.UpdatedAt
}

// EntityToGroupDTO mapper []group.Group to []dto.Group
func EntityToGroupDTO(entities ...*group.Group) []*Group {
	var groups []*Group

	for _, group := range entities {
		g := &Group{}
		g.FromEntity(group)
		groups = append(groups, g)
	}

	return groups
}
