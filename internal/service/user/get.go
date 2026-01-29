package service

import models "github.com/caseapia/goproject-flush/internal/models/user"

func (s *UserService) GetUser(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)

	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.IsDeleted {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetUsersList() ([]models.User, error) {
	var users []models.User

	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}
