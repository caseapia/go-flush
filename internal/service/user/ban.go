package service

import (
	"context"
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (s *UserService) BanUser(
	ctx context.Context,
	adminID uint64,
	userID uint64,
	reason string,
) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.IsBanned {
		return nil, ErrUserBanned
	}

	user.IsBanned = true
	user.BanReason = &reason
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	_ = s.logger.Log(
		ctx,
		adminID,
		userID,
		loggermodule.Ban,
		"Reason: "+reason,
	)

	return user, nil
}

func (s *UserService) UnbanUser(
	ctx context.Context,
	adminID uint64,
	userID uint64,
) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if !user.IsBanned {
		return nil, ErrUserNotBanned
	}

	user.IsBanned = false
	user.BanReason = nil
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	_ = s.logger.Log(
		ctx,
		adminID,
		userID,
		loggermodule.Unban,
	)

	return user, nil
}
