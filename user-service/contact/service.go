package contact

import "context"

// Service service is a façade for use cases.
type Service interface {
	Get(ctx context.Context, userID, contactID string) (*Contact, error)
	GetAll(ctx context.Context, userID string) ([]*Contact, error)
	GetPresence(ctx context.Context, userID, contactID string) (PresenceType, error)
	Create(ctx context.Context, ID, name, lastname, profileID string) error
	Update(ctx context.Context, contact *Contact) error
	Delete(ctx context.Context, userID, contactID string) error
	Block(ctx context.Context, userID, contactID string) error
	Unblock(ctx context.Context, userID, contactID string) error
}

type service struct {
	getUC GetUseCase
	getAllUC GetAllUseCase
	getPresenceUC GetPresenceUseCase
	createUC CreateUseCase
	updateUC UpdateUseCase
	deleteUC DeleteUseCase
	blockUC BlockUseCase
	unblockUC UnblockUseCase
}

// NewService returns a Service instance.
func NewService(
	getUC GetUseCase,
	getAllUC GetAllUseCase,
	getPresenceUC GetPresenceUseCase,
	createUC CreateUseCase,
	updateUC UpdateUseCase,
	deleteUC DeleteUseCase,
	blockUC BlockUseCase,
	unblockUC UnblockUseCase) Service {
	return &service{
		getUC,
		getAllUC,
		getPresenceUC,
		createUC,
		updateUC,
		deleteUC,
		blockUC,
		unblockUC,
	}
}

// Get performs the get use case.
func (s *service) Get(ctx context.Context, userID, contactID string) (*Contact, error) {
	return s.getUC.Execute(ctx, userID, contactID)
}

// GetAll performs the use case to get all.
func (s *service) GetAll(ctx context.Context, userID string) ([]*Contact, error) {
	return s.getAllUC.Execute(ctx, userID)
}

// GetPresence performs the use case to get presence.
func (s *service) GetPresence(ctx context.Context, userID, contactID string) (PresenceType, error) {
	return s.getPresenceUC.Execute(ctx, userID, contactID)
}

// Create performs the creation use case.
func (s *service) Create(ctx context.Context, ID, name, lastname, profileID string) error {
	return s.createUC.Execute(ctx, ID, name, lastname, profileID)
}

// Update performs the update use case.
func (s *service) Update(ctx context.Context, contact *Contact) error {
	return s.updateUC.Execute(ctx, contact)
}

// Delete performs the delete use case.
func (s *service) Delete(ctx context.Context, userID, contactID string) error {
	return s.deleteUC.Execute(ctx, userID, contactID)
}

// Block perform the block use case.
func (s *service) Block(ctx context.Context, userID, contactID string) error {
	return s.blockUC.Execute(ctx, userID, contactID)
}

// Unblock perform the unblock use case.
func (s *service) Unblock(ctx context.Context, userID, contactID string) error {
	return s.unblockUC.Execute(ctx, userID, contactID)
}
