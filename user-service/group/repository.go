package group

// Reader Group interface
type Reader interface {
	Get(ID string) (*Group, error)
	GetAll(userID string) ([]*Group, error)
}

// Writer Group writer
type Writer interface {
	Create(group *Group) error
	Update(group *Group) (int, error)
	Delete(groupID string) (int, error)
	AddMember(member *Member) error
	SetMemberAdmin(groupID, userID string, admin bool) (int, error)
	RemoveMember(groupID, userID string) (bool, error)

}

// Repository interface for Group data source.
type Repository interface {
	Reader
	Writer
}
