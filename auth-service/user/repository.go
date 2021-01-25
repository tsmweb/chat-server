package user

// Reader interface
type Reader interface {
	Get(ID string) (*User, error)
}

// Writer user writer
type Writer interface {
	Create(user *User) error
	Update(user *User) (int, error)
}

// Repository interface for user data source.
type Repository interface {
	Reader
	Writer
}