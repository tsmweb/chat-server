package profile

// Reader interface
type Reader interface {
	Get(ID string) (*Profile, error)
}

// Writer profile writer
type Writer interface {
	Create(profile *Profile) error
	Update(profile *Profile) (int, error)
}

// Repository interface for profile data source.
type Repository interface {
	Reader
	Writer
}