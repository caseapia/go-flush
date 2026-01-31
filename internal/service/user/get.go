package service

import (
	"context"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (s *UserService) GetUser(ctx context.Context, id uint64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil || user.IsDeleted {
		return nil, ErrUserNotFound
	}

	_ = s.logger.Log(ctx, 0, id, loggermodule.SearchByUserID)

	return user, nil
}

func (s *UserService) GetUsersList(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx, 0, 0, loggermodule.SearchByAllUsers)

	return users, nil
}
