package adminUser

import (
	"context"
	"strconv"
	"strings"
	"time"

	loggermodel "github.com/caseapia/goproject-flush/internal/models/logger"
	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
)

func (s *AdminUserService) SetStaffRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, UserError.UserNotFound()
	}

	oldRank := u.StaffRank

	u.StaffRank = rank
	u.UpdatedAt = time.Now()

	addInfo := "Rank before: " + strconv.Itoa(oldRank) + ", Rank after: " + strconv.Itoa(rank)

	_ = s.logger.Log(
		ctx,
		loggermodel.CommonLogger,
		0,
		&userID,
		loggermodel.SetStaffRank,
		addInfo,
	)
	return u, s.repo.Update(ctx, u)
}

func (s *AdminUserService) SetDeveloperRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, UserError.UserNotFound()
	}

	oldRank := u.DeveloperRank

	u.DeveloperRank = rank
	u.UpdatedAt = time.Now()

	addInfo := "Rank before: " + strconv.Itoa(oldRank) + ", Rank after: " + strconv.Itoa(rank)

	_ = s.logger.Log(
		ctx,
		loggermodel.CommonLogger,
		0,
		&userID,
		loggermodel.SetDeveloperRank,
		addInfo,
	)
	return u, s.repo.Update(ctx, u)
}

func (s *AdminUserService) EditFlags(ctx context.Context, userID uint64, flags []string) (*usermodel.User, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, UserError.UserNotFound()
	}

	oldFlags := u.Flags

	u.Flags = flags
	u.UpdatedAt = time.Now()

	addInfo := "Flags before: " + strings.Join(oldFlags, ",") + ", Flags after: " + strings.Join(flags, ",")

	_ = s.logger.Log(ctx, loggermodel.CommonLogger, 0, &userID, loggermodel.ChangeFlags, addInfo)

	return u, s.repo.Update(ctx, u)
}
