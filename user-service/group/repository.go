package group

import (
	"context"
)

// Reader Group interface
type Reader interface {
	Get(ctx context.Context, groupID, userID string) (*Group, error)
	GetAll(ctx context.Context, userID string) ([]*Group, error)
	ExistsUser(ctx context.Context, userID string) (bool, error)
	ExistsGroup(ctx context.Context, groupID string) (bool, error)
	IsGroupAdmin(ctx context.Context, groupID, userID string) (bool, error)
	IsGroupOwner(ctx context.Context, groupID, userID string) (bool, error)
}

// Writer Group writer
type Writer interface {
	Create(ctx context.Context, group *Group) error
	Update(ctx context.Context, group *Group) (bool, error)
	Delete(ctx context.Context, groupID string) (bool, error)
	AddMember(ctx context.Context, member *Member) error
	SetAdmin(ctx context.Context, member *Member) (bool, error)
	RemoveMember(ctx context.Context, groupID, userID string) (bool, error)

}

// Repository interface for Group data source.
type Repository interface {
	Reader
	Writer
}
