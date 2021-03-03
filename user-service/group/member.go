package group

import "time"

// Member data model
type Member struct {
	GroupID   string
	UserID    string
	Admin     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewMember create a new Member
func NewMember(groupID, userID string, admin bool) (*Member, error) {
	m := &Member{
		GroupID: groupID,
		UserID:  userID,
		Admin:   admin,
		CreatedAt: time.Now(),
	}

	err := m.Validate()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Validate model Group
func (m *Member) Validate() error {
	if m.GroupID == "" {
		return ErrGroupIDValidateModel
	}
	if m.UserID == "" {
		return ErrUserIDValidateModel
	}
	return nil
}
