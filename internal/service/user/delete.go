package service

import (
	"context"
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (s *UserService) DeleteUser(ctx context.Context, id uint64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return user, err
	}
	if user == nil || user.IsDeleted {
		return user, ErrUserNotFound
	}

	user.IsDeleted = true
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return user, err
	}

	_ = s.logger.Log(ctx, 0, id, loggermodule.SoftDelete)

	return user, nil
}

func (s *UserService) RestoreUser(ctx context.Context, id uint64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return user, err
	}
	if user == nil || !user.IsDeleted {
		return user, ErrUserAlreadyExists
	}

	user.IsDeleted = false
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return user, err
	}

	_ = s.logger.Log(ctx, 0, id, loggermodule.RestoreUser)

	return user, nil
}
