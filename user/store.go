package user

type Manager struct {
	// TODO: db
}

func (m *Manager) Get(req string) (User, error) {
	return User{}, nil
}

func (m *Manager) Create(req User) (User, error) {
	if req.Role == Regular {
		req.Home = req.Name
	} else {
		req.Home = "."
	}
	return User{}, nil
}

func (m *Manager) Update(u User) (User, error) {
	return User{}, nil
}

func (m *Manager) Delete(req string) (User, error) {
	return User{}, nil
}
