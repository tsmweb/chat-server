package login

// Repository interface for login data source.
type Repository interface {
	Login(login Login) (bool, error)
	Update(login Login) error
}

// repository implementation for Repository interface.
// repository holds the dependencies for Service layer.
type repository struct {
	dao DAO
}

// NewRepository creates a new instance of Repository.
func NewRepository(dao DAO) Repository {
	return &repository{dao}
}

// Login returns if ID and password are valid.
func (r *repository) Login(login Login) (bool, error) {
	return r.dao.Login(login)
}

// Update login data in the data base.
func (r *repository) Update(login Login) error {
	return r.dao.Update(login)
}
