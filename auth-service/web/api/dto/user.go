package dto

import (
	"github.com/tsmweb/auth-service/app/user"
	"time"
)

// User data
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastname"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToEntity mapper dto.User to user.User
func (u *User) ToEntity() *user.User {
	return &user.User{
		ID:       u.ID,
		Name:     u.Name,
		LastName: u.LastName,
		Password: u.Password,
	}
}

// FromEntity mapper user.User to dto.User
func (u *User) FromEntity(entity *user.User) {
	u.ID = entity.ID
	u.Name = entity.Name
	u.LastName = entity.LastName
	u.CreatedAt = entity.CreatedAt
	u.UpdatedAt = entity.UpdatedAt
}
