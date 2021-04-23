package login

import "context"

// Service service is a façade for use cases.
type Service interface {
	Login(ctx context.Context, ID, password string) (string, error)
	Update(ctx context.Context, login *Login) error
}

type service struct {
	loginUC LoginUseCase
	updateUC UpdateUseCase
}

// NewService returns a Service instance.
func NewService(
	loginUC LoginUseCase,
	updateUC UpdateUseCase) Service {
	return &service{
		loginUC,
		updateUC,
	}
}

// Login performs the login use case.
func (s *service) Login(ctx context.Context, ID, password string) (string, error) {
	return s.loginUC.Execute(ctx, ID, password)
}

// Update performs the update use case.
func (s *service) Update(ctx context.Context, login *Login) error {
	return s.updateUC.Execute(ctx, login)
}