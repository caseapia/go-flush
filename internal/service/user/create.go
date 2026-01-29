package service

import models "github.com/caseapia/goproject-flush/internal/models/user"

func (s *UserService) CreateUser(name string) (*models.User, error) {
	existing, _ := s.repo.GetByName(name)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	user := &models.User{Name: name}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
