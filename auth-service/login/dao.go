package login

// DAO interface for login data source.
type DAO interface {
	Login(login Login) (bool, error)
	Update(login Login) error
}
