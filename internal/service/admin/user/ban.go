package adminUser

import (
	"context"
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
)

func (s AdminUserService) BanUser(
	ctx context.Context,
	adminID uint64,
	userID uint64,
	reason string,
) (*models.User, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, UserError.UserNotFound()
	}

	if u.IsBanned {
		return nil, UserError.UserBanned()
	}

	if u.UserHasFlag("NONBANNABLE") {
		return nil, UserError.UserInvalidStatus()
	}

	u.IsBanned = true
	u.BanReason = &reason
	u.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx, loggermodule.PunishmentLogger, adminID, &userID, loggermodule.Ban, "Reason: "+reason)

	return u, nil
}

func (s *AdminUserService) UnbanUser(
	ctx context.Context,
	adminID uint64,
	userID uint64,
) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, userID)

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, UserError.UserNotFound()
	}

	if !user.IsBanned {
		return nil, UserError.UserNotBanned()
	}

	user.IsBanned = false
	user.BanReason = nil
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx, loggermodule.PunishmentLogger, adminID, &userID, loggermodule.Unban)

	return user, nil
}
