package user

type Manager struct {
	// TODO: db
}

func (m *Manager) Get(req string) (User, error) {
	// TODO sanitize
	return User{}, nil
}

func (m *Manager) Create(req User) (User, error) {
	// TODO sanitize
	u := req
	u.Home = req.Name
	// TODO push it
	return u, nil
}

func (m *Manager) Update(u User) (User, error) {
	// TODO sanitize
	return User{}, nil
}

func (m *Manager) Delete(req string) (User, error) {
	// TODO sanitize
	return User{}, nil
}
