package admin

import (
	adminRanks "github.com/caseapia/goproject-flush/internal/handler/admin/ranks"
	adminUser "github.com/caseapia/goproject-flush/internal/handler/admin/user"
	AdminRanksRepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	AdminRanksService "github.com/caseapia/goproject-flush/internal/service/admin/ranks"
	Contracts "github.com/caseapia/goproject-flush/internal/service/contracts"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type AdminModule struct {
	RanksHandler *adminRanks.Handler
	RanksService *AdminRanksService.RanksService
	UserHandler  *adminUser.Handler
}

func NewAdminModule(
	mainDB *bun.DB,
	userRankSetter Contracts.UserRankSetter,
	userHandler *adminUser.Handler,
	logger *logger.LoggerService,
) *AdminModule {
	ranksRepo := AdminRanksRepository.NewRanksRepository(mainDB)
	ranksSrv := AdminRanksService.NewRanksService(ranksRepo, userRankSetter, logger)
	ranksHandler := adminRanks.NewHandler(ranksSrv)

	return &AdminModule{
		RanksHandler: ranksHandler,
		RanksService: ranksSrv,
		UserHandler:  userHandler,
	}
}

func (m *AdminModule) RegisterRoutes(app fiber.Router) {
	m.RanksHandler.RegisterRoutes(app)

	if m.UserHandler != nil {
		m.UserHandler.RegisterRoutes(app)
	}
}
