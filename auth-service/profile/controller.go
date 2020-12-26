package profile

// Controller provides a gateway between routes and use cases
type Controller interface {
	Get(ID string) (Presenter, error)
	Create(profile Presenter) error
	Update(profile Presenter) error
}

type controller struct {
	getUseCase GetUseCase
	createUseCase CreateUseCase
	updateUseCase UpdateUseCase
}

// NewController creates a new instance of Controller.
func NewController(
	getUseCase GetUseCase,
	createUseCase CreateUseCase,
	updateUseCase UpdateUseCase) Controller {

	return &controller{
		getUseCase: getUseCase,
		createUseCase: createUseCase,
		updateUseCase: updateUseCase,
	}
}

// Get a profile by ID.
func (c *controller) Get(ID string) (Presenter, error) {
	presenter := Presenter{}

	p, err := c.getUseCase.Execute(ID)
	if err != nil {
		return presenter, err
	}

	presenter.FromEntity(p)

	return presenter, nil
}

// Create a new profile.
func (c *controller) Create(presenter Presenter) error {
	err := c.createUseCase.Execute(presenter.ID, presenter.Name, presenter.LastName, presenter.Password)
	if err != nil {
		return err
	}

	return nil
}

// Update updates profile data.
func (c *controller) Update(presenter Presenter) error {
	p := presenter.ToEntity()
	err := c.updateUseCase.Execute(p)
	if err != nil {
		return err
	}

	return nil
}