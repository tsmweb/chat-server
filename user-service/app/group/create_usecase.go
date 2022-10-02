package group

import (
	"context"

	"github.com/tsmweb/user-service/common/service"
)

// CreateUseCase creates a new Group, otherwise an error is returned.
type CreateUseCase interface {
	Execute(ctx context.Context, name, description, owner string) (string, error)
}

type createUseCase struct {
	tag        string
	repository Repository
}

// NewCreateUseCase create a new instance of CreateUseCase.
func NewCreateUseCase(r Repository) CreateUseCase {
	return &createUseCase{
		tag:        "group::CreateUseCase",
		repository: r,
	}
}

// Execute performs the creation use case.
func (u *createUseCase) Execute(ctx context.Context, name, description, owner string) (string, error) {
	g, err := NewGroup(name, description, owner)
	if err != nil {
		return "", err
	}

	err = u.repository.Create(ctx, g)
	if err != nil {
		service.Error(owner, u.tag, err)
		return "", err
	}

	return g.ID, nil
}
