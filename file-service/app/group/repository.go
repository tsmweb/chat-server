package group

import "context"

// Repository interface for Group data source.
type Repository interface {
	ExistsGroup(ctx context.Context, groupID string) (bool, error)
	IsGroupMember(ctx context.Context, groupID, userID string) (bool, error)
	IsGroupAdmin(ctx context.Context, groupID, userID string) (bool, error)
}
