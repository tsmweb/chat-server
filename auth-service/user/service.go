package user

import "context"

// Service service is a façade for use cases.
type Service interface {
	Get(ctx context.Context, ID string) (*User, error)
	Create(ctx context.Context, ID, name, lastname, password string) error
	Update(ctx context.Context, user *User) error
}

type service struct {
	getUC GetUseCase
	createUC CreateUseCase
	updateUC UpdateUseCase
}

// NewService returns a Service instance.
func NewService(
	getUC GetUseCase,
	createUC CreateUseCase,
	updateUC UpdateUseCase) Service {
	return &service{
		getUC,
		createUC,
		updateUC,
	}
}

// Get performs the get use case.
func (s *service) Get(ctx context.Context, ID string) (*User, error) {
	return s.getUC.Execute(ctx, ID)
}

// Create performs the creation use case.
func (s *service) Create(ctx context.Context, ID, name, lastname, password string) error {
	return s.createUC.Execute(ctx, ID, name, lastname, password)
}

// Update performs the update use case.
func (s *service) Update(ctx context.Context, user *User) error {
	return s.updateUC.Execute(ctx, user)
}
