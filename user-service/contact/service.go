package contact

// Service service is a façade for use cases.
type Service interface {
	Get(userID, contactID string) (*Contact, error)
	GetAll(userID string) ([]*Contact, error)
	GetPresence(userID, contactID string) (PresenceType, error)
	Create(ID, name, lastname, profileID string) error
	Update(contact *Contact) error
	Delete(userID, contactID string) error
	Block(userID, contactID string) error
	Unblock(userID, contactID string) error
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
func (s *service) Get(userID, contactID string) (*Contact, error) {
	return s.getUC.Execute(userID, contactID)
}

// GetAll performs the use case to get all.
func (s *service) GetAll(userID string) ([]*Contact, error) {
	return s.getAllUC.Execute(userID)
}

// GetPresence performs the use case to get presence.
func (s *service) GetPresence(userID, contactID string) (PresenceType, error) {
	return s.getPresenceUC.Execute(userID, contactID)
}

// Create performs the creation use case.
func (s *service) Create(ID, name, lastname, profileID string) error {
	return s.createUC.Execute(ID, name, lastname, profileID)
}

// Update performs the update use case.
func (s *service) Update(contact *Contact) error {
	return s.updateUC.Execute(contact)
}

// Delete performs the delete use case.
func (s *service) Delete(userID, contactID string) error {
	return s.deleteUC.Execute(userID, contactID)
}

// Block perform the block use case.
func (s *service) Block(userID, contactID string) error {
	return s.blockUC.Execute(userID, contactID)
}

// Unblock perform the unblock use case.
func (s *service) Unblock(userID, contactID string) error {
	return s.unblockUC.Execute(userID, contactID)
}
