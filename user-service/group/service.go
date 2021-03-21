package group

import "context"

// Service service is a façade for use cases.
type Service interface {
	Get(ctx context.Context, groupID string) (*Group, error)
	GetAll(ctx context.Context) ([]*Group, error)
	Create(ctx context.Context, name, description, owner string) (string, error)
	Update(ctx context.Context, group *Group) error
	Delete(ctx context.Context, groupID string) error
	AddMember(ctx context.Context, groupID string, userID string, admin bool) error
	RemoveMember(ctx context.Context, groupID, userID string) error
	SetAdmin(ctx context.Context, member *Member) error
}

type service struct {
	getUC GetUseCase
	getAllUC GetAllUseCase
	createUC CreateUseCase
	updateUC UpdateUseCase
	deleteUC DeleteUseCase
	addMemberUC AddMemberUseCase
	removeMemberUC RemoveMemberUseCase
	setAdminUC SetAdminUseCase
}

// NewService returns a Service instance.
func NewService(
	getUC GetUseCase,
	getAllUC GetAllUseCase,
	createUC CreateUseCase,
	updateUC UpdateUseCase,
	deleteUC DeleteUseCase,
	addMemberUC AddMemberUseCase,
	removeMemberUC RemoveMemberUseCase,
	setAdminUC SetAdminUseCase) Service {
	return &service{
		getUC,
		getAllUC,
		createUC,
		updateUC,
		deleteUC,
		addMemberUC,
		removeMemberUC,
		setAdminUC,
	}
}

// Get performs the GetUseCase.
func (s *service) Get(ctx context.Context, groupID string) (*Group, error) {
	return s.getUC.Execute(ctx, groupID)
}

// GetAll performs the GetAllUseCase.
func (s *service) GetAll(ctx context.Context) ([]*Group, error) {
	return s.getAllUC.Execute(ctx)
}

// Create performs the CreateUseCase.
func (s *service) Create(ctx context.Context, name, description, owner string) (string, error) {
	return s.createUC.Execute(ctx, name, description, owner)
}

// Update performs the UpdateUseCase.
func (s *service) Update(ctx context.Context, group *Group) error {
	return s.updateUC.Execute(ctx, group)
}

// Delete performs the DeleteUseCase.
func (s *service) Delete(ctx context.Context, groupID string) error {
	return s.deleteUC.Execute(ctx, groupID)
}

// AddMember performs the AddMemberUseCase.
func (s *service) AddMember(ctx context.Context, groupID string, userID string, admin bool) error {
	return s.addMemberUC.Execute(ctx, groupID, userID, admin)
}

// RemoveMember performs the RemoveMemberUseCase.
func (s *service) RemoveMember(ctx context.Context, groupID, userID string) error {
	return s.removeMemberUC.Execute(ctx, groupID, userID)
}

// SetAdmin performs the SetAdminUseCase.
func (s *service) SetAdmin(ctx context.Context, member *Member) error {
	return s.setAdminUC.Execute(ctx, member)
}
