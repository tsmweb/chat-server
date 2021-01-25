package user

// Service service is a façade for use cases.
type Service interface {
	Get(ID string) (*User, error)
	Create(ID, name, lastname, password string) error
	Update(user *User) error
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
func (s *service) Get(ID string) (*User, error) {
	return s.getUC.Execute(ID)
}

// Create performs the creation use case.
func (s *service) Create(ID, name, lastname, password string) error {
	return s.createUC.Execute(ID, name, lastname, password)
}

// Update performs the update use case.
func (s *service) Update(user *User) error {
	return s.updateUC.Execute(user)
}
