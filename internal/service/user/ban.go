package service

import (
	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
)

func (s *UserService) BanUser(adminID uint, userID uint, reason string) (*models.User, error) {
	user, err := s.repo.GetByID(userID)

	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.IsBanned {
		return nil, ErrUserBanned
	}

	user.IsBanned = true
	user.BanReason = &reason

	reasonText := "Reason: " + reason
	s.logger.Log(adminID, userID, loggermodule.Ban, reasonText)

	return user, s.repo.Update(user)
}

func (s *UserService) UnbanUser(adminID uint, userID uint) (*models.User, error) {
	user, err := s.repo.GetByID(userID)

	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsBanned {
		return nil, ErrUserNotBanned
	}

	user.IsBanned = false
	user.BanReason = nil

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	s.logger.Log(adminID, userID, loggermodule.Unban)

	return user, s.repo.Update(user)
}
