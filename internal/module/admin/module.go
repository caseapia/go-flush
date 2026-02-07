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

func NewAdminModule(db *bun.DB, userRankSetter Contracts.UserRankSetter, userHandler *adminUser.Handler, logger *logger.LoggerService) *AdminModule {
	ranksRepo := AdminRanksRepository.NewRanksRepository(db)
	ranksSrv := AdminRanksService.NewRanksService(ranksRepo, userRankSetter, logger)
	ranksHandler := adminRanks.NewHandler(ranksSrv)

	return &AdminModule{
		RanksHandler: ranksHandler,
		RanksService: ranksSrv,
		UserHandler:  userHandler,
	}
}

func (m *AdminModule) RegisterRoutes(app fiber.Router) {
	admin := app.Group("/admin")

	// Ranks
	admin.Get("/ranks", m.RanksHandler.GetRanksList)
	admin.Post("/rank/create", m.RanksHandler.CreateRank)
	admin.Post("/setstaff/:id", m.RanksHandler.SetStaffRank)
	admin.Post("/setdeveloper/:id", m.RanksHandler.SetDeveloperRank)

	// User actions
	if m.UserHandler != nil {
		admin.Delete("/delete/:id", m.UserHandler.DeleteUser)
		admin.Put("/restore/:id", m.UserHandler.RestoreUser)
		admin.Put("/create", m.UserHandler.CreateUser)
		admin.Patch("/ban/:id", m.UserHandler.BanUser)
		admin.Delete("/unban/:id", m.UserHandler.UnbanUser)
	}
}
