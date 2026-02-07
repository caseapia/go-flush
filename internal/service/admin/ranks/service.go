package adminRanks

import (
	AdminRanksRepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	"github.com/caseapia/goproject-flush/internal/service/contracts"
	"github.com/caseapia/goproject-flush/internal/service/logger"
)

type RanksService struct {
	repo           *AdminRanksRepository.RanksRepository
	userRankSetter contracts.UserRankSetter
	logger         *logger.LoggerService
}

func NewRanksService(repo *AdminRanksRepository.RanksRepository, userRankSetter contracts.UserRankSetter, logger *logger.LoggerService) *RanksService {
	return &RanksService{
		repo:           repo,
		logger:         logger,
		userRankSetter: userRankSetter,
	}
}

func (s *RanksService) SetUserRankSetter(u contracts.UserRankSetter) {
	s.userRankSetter = u
}
