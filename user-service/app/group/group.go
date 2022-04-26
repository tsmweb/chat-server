package group

import (
	"context"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"strconv"
	"time"
)

var (
	ErrIDValidateModel      = &cerror.ErrValidateModel{Msg: "required id"}
	ErrNameValidateModel    = &cerror.ErrValidateModel{Msg: "required name"}
	ErrOwnerValidateModel   = &cerror.ErrValidateModel{Msg: "required owner"}
	ErrGroupIDValidateModel = &cerror.ErrValidateModel{Msg: "required group_id"}
	ErrUserIDValidateModel  = &cerror.ErrValidateModel{Msg: "required user_id"}
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
	}

	if err = g.Validate(CREATE); err != nil {
		return nil, err
	}

	return g, nil
}

// Validate model Group
func (g *Group) Validate(op Operation) error {
	if g.ID == "" {
		return ErrIDValidateModel
	}
	if g.Name == "" {
		return ErrNameValidateModel
	}
	if op == CREATE && g.Owner == "" {
		return ErrOwnerValidateModel
	}
	return nil
}

// Repository interface for Group data source.
type Repository interface {
	Get(ctx context.Context, groupID, userID string) (*Group, error)
	GetAll(ctx context.Context, userID string) ([]*Group, error)
	ExistsUser(ctx context.Context, userID string) (bool, error)
	ExistsGroup(ctx context.Context, groupID string) (bool, error)
	IsGroupAdmin(ctx context.Context, groupID, userID string) (bool, error)
	IsGroupOwner(ctx context.Context, groupID, userID string) (bool, error)
	Create(ctx context.Context, group *Group) error
	Update(ctx context.Context, group *Group) (bool, error)
	Delete(ctx context.Context, groupID string) (bool, error)
	AddMember(ctx context.Context, member *Member) error
	SetAdmin(ctx context.Context, member *Member) (bool, error)
	RemoveMember(ctx context.Context, groupID, userID string) (bool, error)
}
