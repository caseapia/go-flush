package adminUser

import (
	"context"
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
	AdminError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/admin"
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
)

func (s AdminUserService) BanUser(
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
		return nil, UserError.UserNotFound()
	}

	staffRankID := user.StaffRank
	devRankID := user.DeveloperRank

	adminRank, err := s.rankService.GetByID(ctx, staffRankID)
	if err != nil {
		return nil, err
	}

	if adminRank != nil && adminRank.HasFlag("MANAGER") {
		return nil, AdminError.ManagerRankCannotBeChanged()
	}

	devRank, err := s.rankService.GetByID(ctx, devRankID)
	if err != nil {
		return nil, err
	}

	if devRank != nil && devRank.HasFlag("MANAGER") {
		return nil, AdminError.ManagerRankCannotBeChanged()
	}

	if user.IsBanned {
		return nil, UserError.UserBanned()
	}

	user.IsBanned = true
	user.BanReason = &reason
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx, "punish", adminID, &userID, loggermodule.Ban, "Reason: "+reason)

	return user, nil
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

	_ = s.logger.Log(ctx, "punish", adminID, &userID, loggermodule.Unban)

	return user, nil
}
