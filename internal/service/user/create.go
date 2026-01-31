package service

import (
	"context"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (s *UserService) CreateUser(ctx context.Context, adminID int, name string) (*models.User, error) {
	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	user := &models.User{Name: name}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	newUser, err := s.repo.GetByName(ctx, name)

	if err != nil {
		return nil, err
	}

	if newUser != nil {
		return nil, ErrUserAlreadyExists
	}

	_ = s.logger.Log(ctx, uint64(adminID), 0, loggermodule.Create, "as user "+name)

	return user, nil
}
