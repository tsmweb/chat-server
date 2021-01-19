package login

// Repository interface for login data source.
type Repository interface {
	Login(login *Login) (bool, error)
	Update(login *Login) (int, error)
}
