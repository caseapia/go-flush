package adminRanks

import (
	"context"

	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
	AdminErrorConstructor "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/admin"
	adminRanksRepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	UserRepository "github.com/caseapia/goproject-flush/internal/repository/user"
	"github.com/caseapia/goproject-flush/internal/service/contracts"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/gookit/slog"
)

type UserRankSetter struct {
	userRepo  UserRepository.UserRepository
	ranksRepo *adminRanksRepository.RanksRepository
	logger    *logger.LoggerService
}

func NewUserRankSetter(
	userRepo *UserRepository.UserRepository,
	ranksRepo *adminRanksRepository.RanksRepository,
	logger *logger.LoggerService,
) contracts.UserRankSetter {
	return &UserRankSetter{
		userRepo:  *userRepo,
		ranksRepo: ranksRepo,
		logger:    logger,
	}
}

func (u *UserRankSetter) SetStaffRank(
	ctx context.Context,
	userID uint64,
	rankID int,
) (*usermodel.User, error) {

	rank, err := u.ranksRepo.GetByID(ctx, rankID)
	if err != nil {
		return nil, err
	}

	if rank.HasFlag("DEV") {
		slog.WithData(slog.M{
			"rankID": rankID,
			"userID": userID,
		}).Error("Rank has DEV flag")

		return nil, AdminErrorConstructor.CantIssueStaffRank()
	}

	userModel, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userModel, err = userModel.SetStaffRank(rankID)
	if err != nil {
		return nil, err
	}

	if err := u.userRepo.Update(ctx, userModel); err != nil {
		return nil, err
	}

	return userModel, nil
}

func (u *UserRankSetter) SetDeveloperRank(ctx context.Context, userID uint64, rankID int) (*usermodel.User, error) {
	rank, err := u.ranksRepo.GetByID(ctx, rankID)
	if err != nil {
		return nil, err
	}

	if !rank.HasFlag("DEV") && rank.Name != "None" && rank.Name != "Player" {
		slog.WithData(slog.M{
			"rankID": rankID,
			"userID": userID,
		}).Error("Rank hasn't DEV flag")

		return nil, AdminErrorConstructor.CantIssueDeveloperRank()
	}

	userModel, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userModel, err = userModel.SetDeveloperRank(rankID)
	if err != nil {
		return nil, err
	}

	if err := u.userRepo.Update(ctx, userModel); err != nil {
		return nil, err
	}

	return userModel, nil
}
